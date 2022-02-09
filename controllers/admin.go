package controllers

import (
	beego "github.com/beego/beego/v2/server/web"

	"class_notice/models"
)

type AdminController struct {
	beego.Controller
}

func (c *AdminController) Get() {
	_, _, is := isLogin(c)

	flash := beego.NewFlash()
	if _, ok := flash.Data["notice"]; ok {
		c.TplName = "admin.tpl"
		return
	}

	if !is && c.Ctx.Request.RequestURI != "/login" {
		// token 获取失败或失效，或用户被禁用将跳转登陆
		c.Ctx.Redirect(301, "/login")
		c.Finish()
		return
	}

	c.TplName = "admin.tpl"
}

func (c *AdminController) LoginPage() {
	token, _, is := isLogin(c)

	if is && c.Ctx.Request.RequestURI != "/admin" {

		flash := beego.NewFlash()
		flash.Notice(token)
		flash.Store(&c.Controller)

		c.Redirect("/admin", 301)
		c.Finish()
		return
	}

	c.TplName = "login.tpl"

}

func isLogin(c *AdminController) (token string, user *models.Users, isLogin bool) {
	token, tokenErr := c.GetSecureCookie("bin", "u_token")
	user, userErr := models.TokenGetUser(token)
	// token 获取失败或失效，或用户被禁用将视为未登陆
	if !tokenErr || userErr != nil || user == nil || user.Status != 1 {
		return token, user, false
	} else {
		return token, user, true
	}
}
