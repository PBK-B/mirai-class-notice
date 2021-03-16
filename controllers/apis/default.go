package controllers

import (
	"class_notice/models"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// 在需要登陆的接口加上这个助理函数，未登陆就会返回统一的错误
func userAssistant(c *beego.Controller) (token string, user models.Users, is bool) {
	token, user, is = isLogin(c)
	if !is {
		callBackResult(c, 200, "用户未登陆或 token 失效！", nil)
		c.Finish()
	}
	return
}

// 错误返回结果 json 数据获取函数
func getErrorJson(code int, msg string) string {
	return "{" + "\"code\": " + strconv.Itoa(code) + ",\"msg\": " + "\"" + msg + "\"}"
}

// 接口数据返回接口整理函数
func callBackResult(c *beego.Controller, code int, msg string, data interface{}) {
	c.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	if data == nil {
		_code_str := -1
		c.CustomAbort(code, getErrorJson(_code_str, msg))
	} else {
		_code_str := 1
		c.Data["json"] = map[string]interface{}{
			"code": _code_str,
			"msg":  msg,
			"data": data,
		}
		c.ServeJSON()
	}
}

// 判断是否登陆函数，该函数还可以获取登陆用户
func isLogin(c *beego.Controller) (token string, user models.Users, isLogin bool) {
	token, tokenErr := c.GetSecureCookie("bin", "u_token")
	user, userErr := models.TokenGetUser(token)
	// token 获取失败或失效，或用户被禁用将视为未登陆
	if !tokenErr || userErr != nil || user.Status != 1 {
		return token, user, false
	} else {
		return token, user, true
	}
}
