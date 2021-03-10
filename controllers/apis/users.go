package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type UsersController struct {
	beego.Controller
}

// func (c *UsersController) Get() {
// 	c.Ctx.WriteString("Hello World!!!")
// }

func (c *UsersController) ApiLogin() {
	c.Ctx.WriteString("Hello World!!!")
}
