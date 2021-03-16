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
		web.NSRouter("/login", &apis.UsersController{}, "get:ApiLogin"),

		// 用户相关 API
		web.NSNamespace("/user",
			// 获取用户信息
			web.NSRouter("/me", &apis.UsersController{}, "get:ApiGetMe"),
			web.NSRouter("/create", &apis.UsersController{}, "post:ApiCreateUser"),
			web.NSRouter("/list", &apis.UsersController{}, "get:ApiUserList"),
		),

		// 用户相关 API
		web.NSNamespace("/time",
			web.NSRouter("/create", &apis.TimesController{}, "post:ApiCreateTime"),
			web.NSRouter("/list", &apis.TimesController{}, "get:ApiTimeList"),
		),
	)

	web.AddNamespace(api)

}
