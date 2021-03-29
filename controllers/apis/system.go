package controllers

import (
	"class_notice/models"

	beego "github.com/beego/beego/v2/server/web"
)

type SystemController struct {
	beego.Controller
}

// 获取系统设置
func (c *SystemController) ApiSystemGetInfo() {

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	config, err := models.GetConfigsDataByName("system")

	if err != nil || config == nil {

		config = map[string]interface{}{
			"school_time":   nil, // 开学时间
			"few_weeks":     0,   // 这学期共几周
			"notice_minute": 0,   // 提前多少分钟通知
		}

	}

	// 时间戳转时间对象
	// timer := time.Unix(b_school_time, 0)

	callBackResult(&c.Controller, 200, "", config)
	c.Finish()
}

// 设置开学时间
func (c *SystemController) ApiSystemSetSchoolTime() {
	b_school_time, err := c.GetInt64("time")

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	if b_school_time <= 0 || err != nil {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	config, err := models.GetConfigsDataByName("system")

	if err == nil && config != nil {

		config["school_time"] = b_school_time
		models.UpdateConfigByData("system", config)

	} else {

		models.AddConfigByData("system", map[string]interface{}{
			"school_time": b_school_time,
		})

	}

	// 时间戳转时间对象
	// timer := time.Unix(b_school_time, 0)

	callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
	c.Finish()
}

// 设置这学期共几周
func (c *SystemController) ApiSystemSetFewWeeks() {
	b_few_weeks, err := c.GetInt("weeks")

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	if b_few_weeks <= 0 || err != nil {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	config, err := models.GetConfigsDataByName("system")

	if err == nil && config != nil {

		config["few_weeks"] = b_few_weeks
		models.UpdateConfigByData("system", config)

	} else {

		models.AddConfigByData("system", map[string]interface{}{
			"few_weeks": b_few_weeks,
		})

	}

	// 时间戳转时间对象
	// timer := time.Unix(b_school_time, 0)

	callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
	c.Finish()
}

// 设置提前多少分钟通知
func (c *SystemController) ApiSystemSetNoticeMinute() {
	b_notice_minute, err := c.GetInt("minute")

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	if b_notice_minute <= 0 || err != nil {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	config, err := models.GetConfigsDataByName("system")

	if err == nil && config != nil {

		config["notice_minute"] = b_notice_minute
		models.UpdateConfigByData("system", config)

	} else {

		models.AddConfigByData("system", map[string]interface{}{
			"notice_minute": b_notice_minute,
		})

	}

	// 时间戳转时间对象
	// timer := time.Unix(b_school_time, 0)

	callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
	c.Finish()
}
