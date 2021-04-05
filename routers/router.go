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

		// 时间相关 API
		web.NSNamespace("/time",
			web.NSRouter("/create", &apis.TimesController{}, "post:ApiCreateTime"),
			web.NSRouter("/update", &apis.TimesController{}, "post:ApiUpdateTime"),
			web.NSRouter("/list", &apis.TimesController{}, "get:ApiTimeList"),
			web.NSRouter("/groups", &apis.TimesController{}, "get:ApiTimeGroupList"),
		),

		// 课表相关 API
		web.NSNamespace("/course",
			web.NSRouter("/create", &apis.CoursesController{}, "post:ApiCreateCourses"),
			web.NSRouter("/list", &apis.CoursesController{}, "get:ApiCoursesList"),
		),

		// 系统相关 API
		web.NSNamespace("/system",
			web.NSRouter("/info", &apis.SystemController{}, "get:ApiSystemGetInfo"),
			web.NSRouter("/schooltime", &apis.SystemController{}, "post:ApiSystemSetSchoolTime"),
			web.NSRouter("/fewweeks", &apis.SystemController{}, "post:ApiSystemSetFewWeeks"),
			web.NSRouter("/noticeminute", &apis.SystemController{}, "post:ApiSystemSetNoticeMinute"),
		),

		// 机器人相关 API
		web.NSNamespace("/bot",
			web.NSRouter("/login", &apis.BotController{}, "post:ApiBotLogin"),
			web.NSRouter("/relogin", &apis.BotController{}, "post:ApiBotReLogin"),
			web.NSRouter("/captcha", &apis.BotController{}, "post:ApiBotSubmitCaptcha"),
			web.NSRouter("/info", &apis.BotController{}, "get:ApiBotGetInfo"),
			// web.NSRouter("/test", &apis.BotController{}, "get:ApiBotTest"),
			web.NSNamespace("/update",
				web.NSRouter("/groupcode", &apis.BotController{}, "post:ApiUpdateBotGroupcode"),
			),
		),
	)

	web.AddNamespace(api)

}
