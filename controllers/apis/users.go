package controllers

import (
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	"class_notice/helper"
	"class_notice/models"
)

type UsersController struct {
	beego.Controller
}

// 用户登陆后台接口
func (c *UsersController) ApiLogin() {

	u_name := c.GetString("name")
	u_password := c.GetString("password")

	if u_name == "" || u_password == "" {
		callBackResult(c, 403, "参数异常", nil)
		return
	}

	// TODO: Bin BY 这里应该还可以做判断用户名和密码是否合法
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

// 获取后台用户信息接口
func (c *UsersController) ApiGetMe() {

	_, user, _ := userAssistant(c)

	c.Data["json"] = map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"token": user.Token,
	}
	callBackResult(c, 200, "", c.Data["json"])
	c.Finish()
}

// 创建后台管理用户接口
func (c *UsersController) ApiCreateUser() {

	// 要求登陆助理函数
	userAssistant(c)

	u_name := c.GetString("name")
	u_password := c.GetString("password")

	if u_name == "" || u_password == "" {
		callBackResult(c, 403, "参数异常", nil)
		c.Finish()
		return
	}

	// TODO: Bin BY 这里应该还可以做判断用户名和密码是否合法
	user := models.Users{Name: u_name, Password: helper.StringToMd5(u_password), Status: 1}
	user, err := models.AddUsers(&user)

	if err != nil {
		callBackResult(c, 403, "用户创建失败"+err.Error(), nil)
		c.Finish()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"token": user.Token,
	}
	callBackResult(c, 200, "", c.Data["json"])
}

// 在需要登陆的接口加上这个助理函数，未登陆就会返回统一的错误
func userAssistant(c *UsersController) (token string, user models.Users, is bool) {
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

// 判断是否登陆函数，该函数还可以获取登陆用户
func isLogin(c *UsersController) (token string, user models.Users, isLogin bool) {
	token, tokenErr := c.GetSecureCookie("bin", "u_token")
	user, userErr := models.TokenGetUser(token)
	// token 获取失败或失效，或用户被禁用将视为未登陆
	if !tokenErr || userErr != nil || user.Status != 1 {
		return token, user, false
	} else {
		return token, user, true
	}
}
