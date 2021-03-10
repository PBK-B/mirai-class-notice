package routers

import (
	"class_notice/controllers"

	web "github.com/beego/beego/v2/server/web"
	context "github.com/beego/beego/v2/server/web/context"
)

func init() {
	web.Router("/", &controllers.MainController{})
	web.Router("/admin", &controllers.AdminController{})
	web.Get("/a", func(ctx *context.Context) {
		ctx.Output.Body([]byte("hello world"))
	})
}
