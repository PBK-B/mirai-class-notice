package controllers

import (
	"class_notice/models"

	beego "github.com/beego/beego/v2/server/web"
)

type BotController struct {
	beego.Controller
}

type BotConfigDataJsonType struct {
	account      int64
	password     string
	password_md5 string
	group_code   int64
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
