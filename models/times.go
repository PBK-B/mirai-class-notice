package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		time, err = GetTimesById(int(id))
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
	for _, time := range ls {
		startStr := time.Start
		s := strings.Split(startStr, ":")

		if s[0] != "" && s[1] != "" {
			// 获取到时间

			h, _ := strconv.Atoi(s[0])
			m, _ := strconv.Atoi(s[1])

			if h >= 0 && h < 24 && m >= 0 && m <= 60 {
				// 时间合法
				ts := fmt.Sprintf("0 %s %s * * *", fmt.Sprint(m), fmt.Sprint(h))
				tk := task.NewTask(time.Remarks, ts, func(ctx context.Context) error {
					fmt.Println("执行了 ", time.Remarks)
					PushTimeAllCourses(time.Id) // 通知筛选课程
					return nil
				})
				task.AddTask(time.Remarks, tk) // 将任务添加到全局任务
			}

			// fmt.Println(s[0], s[1], len(s))
		}
	}

	task.StartTask() // 启动全局任务
}
