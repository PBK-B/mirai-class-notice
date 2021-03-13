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

	// web.Router("/api/login", &apis.UsersController{}, "get:ApiLogin")

	api := web.NewNamespace("/api",
		// 控制台登陆
		web.NSRouter("/login", &apis.UsersController{}, "*:ApiLogin"),

		// 用户相关 API
		web.NSNamespace("/user",
			// 获取用户信息
			web.NSRouter("/me", &apis.UsersController{}, "*:GetMe"),
		),
	)

	web.AddNamespace(api)

}
