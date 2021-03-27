package controllers

import (
	"class_notice/models"

	beego "github.com/beego/beego/v2/server/web"
)

type BotController struct {
	beego.Controller
}

func (c *BotController) ApiBotLogin() {

	b_account, _ := c.GetInt64("account")
	b_password := c.GetString("password")

	if b_account == 0 || b_password == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)
	// helper.InitBot(b_account, b_password)

	config, err := models.GetConfigsDataByName("bot")

	if err == nil && config != nil {
		b_p_md5 := config["password_md5"].(string)
		b_g_c := config["group_code"].(float64)
		models.UpdateConfigByData("bot", map[string]interface{}{
			"account":      b_account,
			"password":     b_password,
			"password_md5": b_p_md5,
			"group_code":   b_g_c,
		})
	} else {

		models.AddConfigByData("bot", map[string]interface{}{
			"account":      b_account,
			"password":     b_password,
			"password_md5": "",
			"group_code":   0,
		})
	}

	callBackResult(&c.Controller, 200, "ok", nil)
	c.Finish()
}

func (c *BotController) ApiUpdateBotGroupcode() {
	group_code, _ := c.GetInt64("group_code")

	if group_code == 0 {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)
	// helper.InitBot(b_account, b_password)

	config, err := models.GetConfigsDataByName("bot")

	if err == nil && config != nil {
		b_a := config["account"].(float64)
		b_p := config["password"].(string)
		b_p_md5 := config["password_md5"].(string)
		models.UpdateConfigByData("bot", map[string]interface{}{
			"account":      b_a,
			"password":     b_p,
			"password_md5": b_p_md5,
			"group_code":   group_code,
		})
	} else {
		models.AddConfigByData("bot", map[string]interface{}{
			"account":      0,
			"password":     "",
			"password_md5": "",
			"group_code":   group_code,
		})
	}

	callBackResult(&c.Controller, 200, "ok", nil)
	c.Finish()
}
