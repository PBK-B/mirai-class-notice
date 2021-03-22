package controllers

import (
	"class_notice/helper"

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
	helper.InitBot(b_account, b_password)

	callBackResult(&c.Controller, 200, "ok", nil)
	c.Finish()
}
