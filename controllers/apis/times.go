package controllers

import (
	"class_notice/models"

	beego "github.com/beego/beego/v2/server/web"
)

type TimesController struct {
	beego.Controller
}

// 创建上课时间
func (c *TimesController) ApiCreateTime() {

	t_group := c.GetString("group")
	t_start := c.GetString("start")
	t_end := c.GetString("end")
	t_remark := c.GetString("remark")

	if t_group == "" || t_start == "" || t_end == "" || t_remark == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	time := models.Times{
		Group:   t_group,
		Start:   t_start,
		End:     t_end,
		Remarks: t_remark,
	}
	_time, err := models.AddTimes(&time)

	if err != nil || _time == nil {
		callBackResult(&c.Controller, 200, "时间创建失败 "+err.Error(), nil)
		c.Finish()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"id":     _time.Id,
		"group":  _time.Group,
		"start":  _time.Start,
		"end":    _time.End,
		"remark": _time.Remarks,
	}
	callBackResult(&c.Controller, 200, "", c.Data["json"])
}

// 更新上课时间
func (c *TimesController) ApiUpdateTime() {

	t_id, _ := c.GetInt("id")
	t_group := c.GetString("group")
	t_start := c.GetString("start")
	t_end := c.GetString("end")
	t_remark := c.GetString("remark")

	if t_id == 0 || t_group == "" || t_start == "" || t_end == "" || t_remark == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	time := models.Times{
		Id:      t_id,
		Group:   t_group,
		Start:   t_start,
		End:     t_end,
		Remarks: t_remark,
	}
	err := models.UpdateTimesById(&time)

	if err != nil {
		callBackResult(&c.Controller, 200, "时间更新失败 "+err.Error(), nil)
		c.Finish()
		return
	}

	_time, err := models.GetTimesById(t_id)
	if err != nil {
		callBackResult(&c.Controller, 200, "时间更新失败 "+err.Error(), nil)
		c.Finish()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"id":     _time.Id,
		"group":  _time.Group,
		"start":  _time.Start,
		"end":    _time.End,
		"remark": _time.Remarks,
	}
	callBackResult(&c.Controller, 200, "", c.Data["json"])
}

// 获取全部上课时间
func (c *TimesController) ApiTimeList() {

	u_count, _ := c.GetInt("count", 10)
	u_page, _ := c.GetInt("page", 0)

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	times, err := models.GetAllTimes(nil, nil, nil, nil, int64(u_page), int64(u_count))

	if err != nil {
		callBackResult(&c.Controller, 403, "时间列表获取失败", nil)
		c.Finish()
		return
	}

	var new_time []interface{}

	for item := range times {
		i_t := times[item]
		new_t := map[string]interface{}{
			"id":     i_t.Id,
			"group":  i_t.Group,
			"start":  i_t.Start,
			"end":    i_t.End,
			"remark": i_t.Remarks,
		}
		new_time = append(new_time, new_t)
	}

	callBackResult(&c.Controller, 200, "", new_time)
	c.Finish()
}

func (c *TimesController) ApiTimeGroupList() {
	groups := models.GetAllTimesGroups()
	callBackResult(&c.Controller, 200, "", groups)
	c.Finish()
}

func (c *TimesController) ApiTimeTest() {
	// fmt.Println("-------------- Hello Test 001 --------------")
	// helper.InitBot()
	models.PushTimeAllCourses(1)
	callBackResult(&c.Controller, 200, "", nil)
	c.Finish()
}
