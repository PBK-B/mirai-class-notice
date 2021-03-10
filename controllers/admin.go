package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type AdminController struct {
	beego.Controller
}

func (c *AdminController) Get() {
	c.Ctx.WriteString("Hello World!!!")
}
