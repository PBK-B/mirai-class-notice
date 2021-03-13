package controllers

import (
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	"class_notice/models"
)

type UsersController struct {
	beego.Controller
}

// func (c *UsersController) Get() {
// 	c.Ctx.WriteString("Hello World!!!")
// }

func (c *UsersController) ApiLogin() {

	u_name := c.GetString("name")
	u_password := c.GetString("password")

	if u_name == "" || u_password == "" {
		callBackResult(c, 403, "参数异常", nil)
		return
	}
	user, err := models.LoginUser(u_name, u_password)

	if err == orm.ErrNoRows {
		callBackResult(c, 200, "用户名或密码错误", nil)
		return
	} else if err == orm.ErrMissPK {
		callBackResult(c, 403, "服务器错误", nil)
		return
	} else {
		if user.Status == 0 {
			callBackResult(c, 200, "账号被禁用", nil)
			return
		}

		c.Data["json"] = map[string]interface{}{
			"id":    user.Id,
			"name":  user.Name,
			"token": user.Token,
		}
		c.SetSecureCookie("bin", "u_token", user.Token)
		callBackResult(c, 200, "", c.Data["json"])

		return
	}
}

func getErrorJson(code int, msg string) string {
	return "{" + "\"code\": " + strconv.Itoa(code) + ",\"msg\": " + "\"" + msg + "\"}"
}

func callBackResult(c *UsersController, code int, msg string, data interface{}) {
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
