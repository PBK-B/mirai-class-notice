package controllers

import (
	"class_notice/models"

	beego "github.com/beego/beego/v2/server/web"
)

type AdminController struct {
	beego.Controller
}

func (c *AdminController) Get() {
	_, _, is := isLogin(c)

	if !is {
		// token 获取失败或失效，或用户被禁用将跳转登陆
		c.Redirect("/login", 301)
		c.Finish()
	}

	c.TplName = "admin.tpl"
}

func (c *AdminController) LoginPage() {
	_, _, is := isLogin(c)
	if is {
		c.Redirect("/admin", 301)
	} else {
		c.TplName = "login.tpl"
	}
}

func isLogin(c *AdminController) (token string, user models.Users, isLogin bool) {
	token, tokenErr := c.GetSecureCookie("bin", "u_token")
	user, userErr := models.TokenGetUser(token)
	// token 获取失败或失效，或用户被禁用将视为未登陆
	if !tokenErr || userErr != nil || user.Status != 1 {
		return token, user, false
	} else {
		return token, user, true
	}
}
