package models

import (
	"class_notice/helper"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Courses struct {
	Id           int    `orm:"description(ID)"`
	Status       int    `orm:"description(状态: 通知[1] 暂停[0])"`
	Title        string `orm:"size(128);description(课程名称)"`
	Classroom    string `orm:"size(128);description(上课教学楼)"`
	ClassroomId  string `orm:"size(128);description(上课教室)"`
	Teacher      string `orm:"size(128);description(老师名字)"`
	Remarks      string `orm:"size(128);description(备注)"`
	WeekTime     int    `orm:"description(星期几上课)"`
	LessonSerial int    `orm:"description(第几节课)"`
	Cycle        string `orm:"description(哪些周需要上课)"`
	Times        *Times `orm:"rel(fk)"`
}

func init() {
	// orm.RegisterModel(new(Courses))
}

// AddCourses insert a new Courses into database and returns
// last inserted Id on success.
func AddCourses(m *Courses) (courses *Courses, err error) {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		courses, err = GetCoursesById(int(id))
	}
	return
}

// GetCoursesById retrieves Courses by Id. Returns error if
// Id doesn't exist
func GetCoursesById(id int) (v *Courses, err error) {
	o := orm.NewOrm()
	v = &Courses{Id: id}
	if err = o.QueryTable(new(Courses)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCourses retrieves all Courses matches certain condition. Returns empty list if
// no records exist
func GetAllCourses(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []Courses, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Courses))
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

	var l []Courses
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

// UpdateCourses updates Courses by Id and returns error if
// the record to be updated doesn't exist
func UpdateCoursesById(m *Courses) (err error) {
	o := orm.NewOrm()
	v := Courses{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCourses deletes Courses by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCourses(id int) (err error) {
	o := orm.NewOrm()
	v := Courses{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Courses{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// 判断课程是否已经添加过
func CoursesRepeat(m *Courses) bool {
	o := orm.NewOrm()
	err := o.QueryTable(new(Courses)).Filter("Title", m.Title).Filter("Classroom", m.Classroom).Filter("ClassroomId", m.ClassroomId).Filter("Teacher", m.Teacher).Filter("WeekTime", m.WeekTime).Filter("LessonSerial", m.LessonSerial).Filter("Cycle", m.Cycle).RelatedSel().One(m)
	if err == nil {
		return true
	} else {
		return false
	}
}

func CoursesToMap(course Courses) map[string]interface{} {
	var _cycle_obj interface{}
	json.Unmarshal([]byte(course.Cycle), &_cycle_obj)
	cycle_obj := _cycle_obj.([]interface{})

	return map[string]interface{}{
		"id":           course.Id,
		"title":        course.Title,
		"status":       course.Status,
		"classroom":    course.Classroom,
		"classroom_id": course.ClassroomId,
		"teacher":      course.Teacher,
		"remark":       course.Remarks,
		"week_time":    course.WeekTime,
		"cycle":        cycle_obj,
		"time":         TimesToMap(*course.Times),
	}
}

// 生成提醒消息
func GenerateNoticeString(course Courses) string {
	return fmt.Sprintf("提醒上课小助手，今天 %s 将在%s %s 上【%s】课，%s", course.Times.Start, course.Classroom, course.ClassroomId, course.Title, course.Remarks)
}

// 通知指定时间 id 的全部课程
func PushTimeAllCourses(timeId int) {
	// t_week := int(time.Now().Weekday())

	bot_group_code := int64(0)

	botConfig, botConfigErr := GetConfigsDataByName("bot")
	systemConfig, systemConfigErr := GetConfigsDataByName("system")
	if botConfigErr != nil || botConfig == nil {
		log.Panicln("机器人未配置！")
		return
	}

	// 获取用户设置的通知 QQ 群号
	bot_group_code = int64(botConfig["group_code"].(float64))

	// 获取当前是星期几
	timer := time.Now()
	t_week := timer.Weekday()

	// 获取当前是距离开学之后的第几周
	t_cycle := 0
	timerOld := timer
	if systemConfig != nil && systemConfigErr == nil {
		// 获取设置的时间成功
		timerOld = time.Unix(int64(systemConfig["school_time"].(float64)), 0)
	}
	// 计算两个时间的差值 返回的是纳秒 按需求自行计算其他单位
	timeDuration := timer.Sub(timerOld)
	t_cycle = int(timeDuration/1000/1000/1000/24/60/60/7) + 1 // 这里 +1 是因为当去除小数点加一才等于当前周

	// 获取指定时间状态为启用通知的课程
	ls, _ := GetAllCourses(map[string]string{"times_id": fmt.Sprint(timeId), "status": "1"}, nil, nil, nil, 0, -1)
	for _, course := range ls {

		// 将用户设置的周日(7) 转换为国际时间规定的周日(0)
		week_time := course.WeekTime
		if week_time >= 7 {
			week_time = 0
		}

		if week_time == int(t_week) {
			// 是今天的课，开始判断是否是这周的课。

			var _cycle_obj interface{}
			json.Unmarshal([]byte(course.Cycle), &_cycle_obj)
			cycle_list := _cycle_obj.([]interface{})
			for _, cycle := range cycle_list {
				c_cycle := int(cycle.(float64))

				// fmt.Printf("%d  课程是第几周：%d  当前是第几周：%d \n", course.Id, c_cycle, t_cycle)

				if c_cycle == t_cycle {
					// 是这周且是今天的课
					msg := GenerateNoticeString(course)
					helper.SendGroupMessage(bot_group_code, msg)
					break
				}
			}
		}
	}
}
