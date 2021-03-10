package routers

import (
	"class_notice/controllers"
	apis "class_notice/controllers/apis"

	web "github.com/beego/beego/v2/server/web"
)

func init() {
	web.Router("/", &controllers.MainController{})
	web.Router("/admin", &controllers.AdminController{})
	web.Router("/login", &controllers.AdminController{}, "get:LoginPage")

	web.NewNamespace("/api",
		web.NSNamespace("/user",
			// 控制台登陆
			web.NSRouter("/login", &apis.UsersController{}, "get:ApiLogin"),
		),
	)

}
