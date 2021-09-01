package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/task"
)

type Times struct {
	Id      int
	Group   string
	Start   string
	End     string
	Remarks string `orm:"size(128)"`
}

func init() {
	// orm.RegisterModel(new(Times))
}

// AddTimes insert a new Times into database and returns
// last inserted Id on success.
func AddTimes(m *Times) (time *Times, err error) {

	s_s := strings.Split(m.Start, ":")
	s_e := strings.Split(m.End, ":")

	if len(s_s) != 2 || len(s_e) != 2 {
		return nil, errors.New("time: time Start or time End Illegal")
	}

	s_s0_i, err_0 := strconv.Atoi(s_s[0])
	s_s1_i, err_1 := strconv.Atoi(s_s[1])
	s_e0_i, err_2 := strconv.Atoi(s_e[0])
	s_e1_i, err_3 := strconv.Atoi(s_e[1])

	if err_0 != nil ||
		err_1 != nil ||
		err_2 != nil ||
		err_3 != nil ||
		s_s0_i > 24 ||
		s_e0_i > 24 ||
		s_s1_i > 60 ||
		s_e1_i > 60 {
		// 设置的时间过大或者不是数值
		return nil, errors.New("time: The time set is too large or is not a numerical value")
	}

	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		time, err = GetTimesById(int(id))
		TimesRunAllTask()
	}
	return
}

// GetTimesById retrieves Times by Id. Returns error if
// Id doesn't exist
func GetTimesById(id int) (v *Times, err error) {
	o := orm.NewOrm()
	v = &Times{Id: id}
	if err = o.QueryTable(new(Times)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTimes retrieves all Times matches certain condition. Returns empty list if
// no records exist
func GetAllTimes(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Times, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Times))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Times
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			// for _, v := range l {
			// 	m := make(map[string]interface{})
			// 	val := reflect.ValueOf(v)
			// 	for _, fname := range fields {
			// 		m[fname] = val.FieldByName(fname).Interface()
			// 	}
			// 	ml = append(ml, m)
			// }
		}
		return ml, nil
	}
	return nil, err
}

// UpdateTimes updates Times by Id and returns error if
// the record to be updated doesn't exist
func UpdateTimesById(m *Times) (err error) {
	o := orm.NewOrm()
	v := Times{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
		TimesRunAllTask()
	}
	return
}

// DeleteTimes deletes Times by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTimes(id int) (err error) {
	o := orm.NewOrm()
	v := Times{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Times{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllTimesGroups() (groups []*string) {
	o := orm.NewOrm()
	var l []Times
	o.QueryTable(new(Times)).Filter("group__isnull", false).Distinct().All(&l, "group")
	for item := range l {
		groups = append(groups, &l[item].Group)
	}
	return
}

func TimesToMap(time Times) map[string]interface{} {
	return map[string]interface{}{
		"id":     time.Id,
		"group":  time.Group,
		"start":  time.Start,
		"end":    time.End,
		"remark": time.Remarks,
	}
}

func TimesRunAllTask() {
	task.ClearTask() // 清理全部全局任务

	ls, _ := GetAllTimes(nil, nil, nil, nil, 0, -1)
	for index, time := range ls {
		startStr := time.Start
		s := strings.Split(startStr, ":")

		if len(s) < 2 {
			// 时间不合法
			continue
		}

		if s[0] != "" && s[1] != "" {
			// 获取到时间

			h, _ := strconv.Atoi(s[0])
			m, _ := strconv.Atoi(s[1])

			if h >= 0 && h < 24 && m >= 0 && m <= 60 {
				// 时间合法
				// TODO: 这里需要加上一个提前一个设定的时间
				var notifyTime int64 = 60
				if timeConfig, timeErr := GetConfigsDataByName("system"); timeConfig != nil && timeErr == nil {
					notifyTime = int64(timeConfig["notice_minute"].(float64))
				}

				nTime := transformEstimateTime(startStr, notifyTime)
				h = nTime.Hour()
				m = nTime.Minute()

				// h = h - 1
				// if h < 0 {
				// 	// FIX 修复凌晨 0 点减去一个小时后导致变成负数
				// 	h = 23
				// }
				// 之前是默认提前一个小时通知

				ts := fmt.Sprintf("0 %s %s * * *", fmt.Sprint(m), fmt.Sprint(h))
				new_time := ls[index]
				taskName := new_time.Remarks + fmt.Sprint(new_time.Id)
				tk := task.NewTask(taskName, ts, func(ctx context.Context) error {
					log.Println("执行了 ", new_time.Remarks)
					PushTimeAllCourses(new_time.Id) // 通知筛选课程
					return nil
				})
				task.AddTask(taskName, tk) // 将任务添加到全局任务
			}

			task.StartTask()

			// fmt.Println(s[0], s[1], len(s))
		}
	}

	task.StartTask() // 启动全局任务
	log.Println("【models.Time】注册全部上课时间定时任务成功！")
}

// 转换时间工具方法
func transformEstimateTime(date string, minutes int64) time.Time {
	s := minutes * 60
	ts := fmt.Sprintf("2010-01-01 %s:00", date)
	t, _ := time.ParseInLocation("2006-01-02 15:04:00", ts, time.Local)
	newT := t.Unix() - s
	tm := time.Unix(newT, 0)
	return tm
}
