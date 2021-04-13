package controllers

import (
	"class_notice/models"
	"encoding/json"

	beego "github.com/beego/beego/v2/server/web"
)

type CoursesController struct {
	beego.Controller
}

// 添加一个课程
func (c *CoursesController) ApiCreateCourses() {

	c_time_id, _ := c.GetInt("time_id")             // 绑定时间 id
	c_title := c.GetString("title")                 // 标题
	c_classroom := c.GetString("classroom")         // 教室
	c_classroom_id := c.GetString("classroom_id")   // 教室 id
	c_teacher := c.GetString("teacher")             // 教师名字
	c_remarks := c.GetString("remarks")             // 备注
	c_week_time, _ := c.GetInt("week_time")         // 星期几
	c_lesson_serial, _ := c.GetInt("lesson_serial") // 第几节课
	c_cycle := c.GetString("cycle")                 // 哪些周需要上课

	is_c_cycle := json.Valid([]byte(c_cycle)) // 判断获取到的数据是否为 json

	// 获取 body 数组数据
	// var c_cycle []string // 哪些周需要上课
	// c.Ctx.Input.Bind(&c_cycle, "cycle")

	if c_time_id == 0 ||
		c_title == "" ||
		c_classroom == "" ||
		c_classroom_id == "" ||
		c_teacher == "" ||
		c_remarks == "" ||
		c_week_time == 0 ||
		c_lesson_serial == 0 ||
		!is_c_cycle {
		callBackResult(&c.Controller, 501, "参数异常", nil)
		c.Finish()
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	// 解码 json
	// var test interface{}
	// json.Unmarshal([]byte("{ \"week_time\": 1, \"lesson_id\": \"Hello\", \"cycle\": { \"a\": 111 }, \"list\": [1,2,3] }"), &test)
	// test_obj := test.(map[string]interface{})
	// fmt.Println(test_obj["cycle"].(map[string]interface{})["a"])

	// 编码 json
	// b, err := json.Marshal(test_obj)
	// if err == nil {
	// 	log.Println(string(b))
	// }

	time, err_time := models.GetTimesById(c_time_id)

	if err_time != nil || time == nil {
		callBackResult(&c.Controller, 501, "上课时间 ID 异常", nil)
		c.Finish()
		return
	}

	course := &models.Courses{
		Title:        c_title,
		Status:       1, // 默认启用状态
		Classroom:    c_classroom,
		ClassroomId:  c_classroom_id,
		Teacher:      c_teacher,
		Remarks:      c_remarks,
		WeekTime:     c_week_time,
		LessonSerial: c_lesson_serial,
		Cycle:        c_cycle,
		Times:        time,
	}

	if is_courses := models.CoursesRepeat(course); is_courses {
		callBackResult(&c.Controller, 200, "该课程已经添加过了", nil)
		c.Finish()
		return
	}

	course, err := models.AddCourses(course)

	if err != nil {
		callBackResult(&c.Controller, 403, "添加课程失败"+err.Error(), nil)
		c.Finish()
		return
	}

	c.Data["json"] = models.CoursesToMap(*course)
	callBackResult(&c.Controller, 200, "", c.Data["json"])
}

// 修改一个课程数据
func (c *CoursesController) ApiUpdateCourses() {

	c_id, _ := c.GetInt("id")                       // 课程 id
	c_time_id, _ := c.GetInt("time_id")             // 绑定时间 id
	c_title := c.GetString("title")                 // 标题
	c_classroom := c.GetString("classroom")         // 教室
	c_classroom_id := c.GetString("classroom_id")   // 教室 id
	c_teacher := c.GetString("teacher")             // 教师名字
	c_remarks := c.GetString("remarks")             // 备注
	c_week_time, _ := c.GetInt("week_time")         // 星期几
	c_lesson_serial, _ := c.GetInt("lesson_serial") // 第几节课
	c_cycle := c.GetString("cycle")                 // 哪些周需要上课

	is_c_cycle := json.Valid([]byte(c_cycle)) // 判断获取到的数据是否为 json

	// 获取 body 数组数据
	// var c_cycle []string // 哪些周需要上课
	// c.Ctx.Input.Bind(&c_cycle, "cycle")

	if c_id == 0 ||
		c_time_id == 0 ||
		c_title == "" ||
		c_classroom == "" ||
		c_classroom_id == "" ||
		c_teacher == "" ||
		c_remarks == "" ||
		c_week_time == 0 ||
		c_lesson_serial == 0 ||
		!is_c_cycle {
		callBackResult(&c.Controller, 501, "参数异常", nil)
		c.Finish()
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	// 解码 json
	// var test interface{}
	// json.Unmarshal([]byte("{ \"week_time\": 1, \"lesson_id\": \"Hello\", \"cycle\": { \"a\": 111 }, \"list\": [1,2,3] }"), &test)
	// test_obj := test.(map[string]interface{})
	// fmt.Println(test_obj["cycle"].(map[string]interface{})["a"])

	// 编码 json
	// b, err := json.Marshal(test_obj)
	// if err == nil {
	// 	log.Println(string(b))
	// }

	time, err_time := models.GetTimesById(c_time_id)

	if err_time != nil || time == nil {
		callBackResult(&c.Controller, 501, "上课时间 ID 异常", nil)
		c.Finish()
		return
	}
	is_course, is_course_err := models.GetCoursesById(c_id)

	if is_course == nil || is_course_err != nil {
		callBackResult(&c.Controller, 200, "该课程不存在或课程数据异常！", nil)
		c.Finish()
		return
	}

	course := &models.Courses{
		Id:           c_id,
		Title:        c_title,
		Status:       is_course.Status,
		Classroom:    c_classroom,
		ClassroomId:  c_classroom_id,
		Teacher:      c_teacher,
		Remarks:      c_remarks,
		WeekTime:     c_week_time,
		LessonSerial: c_lesson_serial,
		Cycle:        c_cycle,
		Times:        time,
	}

	err := models.UpdateCoursesById(course)

	if err != nil {
		callBackResult(&c.Controller, 403, "修改课程失败"+err.Error(), nil)
		c.Finish()
		return
	}

	c.Data["json"] = models.CoursesToMap(*course)
	callBackResult(&c.Controller, 200, "", c.Data["json"])
}

// 获取全部上课时间
func (c *CoursesController) ApiCoursesList() {

	u_count, _ := c.GetInt("count", 10)
	u_page, _ := c.GetInt("page", 0)

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	courses, err := models.GetAllCourses(nil, nil, nil, nil, int64(u_page), int64(u_count))

	if err != nil {
		callBackResult(&c.Controller, 403, "时间列表获取失败", nil)
		c.Finish()
		return
	}

	var new_courses []interface{}

	for item := range courses {
		new_c := models.CoursesToMap(courses[item])
		new_courses = append(new_courses, new_c)
	}

	callBackResult(&c.Controller, 200, "", new_courses)
	c.Finish()
}

// 获取一个课程数据
func (c *CoursesController) ApiGetCourses() {

	c_id, _ := c.GetInt("id") // 课程 id

	if c_id == 0 {
		callBackResult(&c.Controller, 501, "参数异常", nil)
		c.Finish()
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	// 解码 json
	// var test interface{}
	// json.Unmarshal([]byte("{ \"week_time\": 1, \"lesson_id\": \"Hello\", \"cycle\": { \"a\": 111 }, \"list\": [1,2,3] }"), &test)
	// test_obj := test.(map[string]interface{})
	// fmt.Println(test_obj["cycle"].(map[string]interface{})["a"])

	// 编码 json
	// b, err := json.Marshal(test_obj)
	// if err == nil {
	// 	log.Println(string(b))
	// }

	course, err := models.GetCoursesById(c_id)

	if err != nil || course == nil {
		callBackResult(&c.Controller, 501, "课程不存在或数据异常！", nil)
		c.Finish()
		return
	}

	c.Data["json"] = models.CoursesToMap(*course)
	callBackResult(&c.Controller, 200, "", c.Data["json"])
}
