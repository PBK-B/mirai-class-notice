package controllers

import (
	"class_notice/helper"
	"class_notice/models"
	"fmt"
	"strconv"

	"github.com/Logiase/MiraiGo-Template/bot"
	beego "github.com/beego/beego/v2/server/web"
)

type BotController struct {
	beego.Controller
}

func (c *BotController) ApiBotLogin() {

	b_account, _ := c.GetInt64("account")
	b_password := c.GetString("password")

	if b_account == 0 || b_password == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)
	// helper.InitBot(b_account, b_password)

	config, err := models.GetConfigsDataByName("bot")

	if err == nil && config != nil {
		b_p_md5 := config["password_md5"].(string)
		b_g_c := config["group_code"].(float64)
		models.UpdateConfigByData("bot", map[string]interface{}{
			"account":      b_account,
			"password":     b_password,
			"password_md5": b_p_md5,
			"group_code":   b_g_c,
		})
	} else {

		models.AddConfigByData("bot", map[string]interface{}{
			"account":      b_account,
			"password":     b_password,
			"password_md5": "",
			"group_code":   0,
		})
	}

	go helper.InitBot(b_account, b_password)

	if err != nil {
		callBackResult(&c.Controller, 200, err.Error(), nil)
		c.Finish()
		return
	}

	callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
	c.Finish()
}

func (c *BotController) ApiBotReLogin() {

	// 要求登陆助理函数
	userAssistant(&c.Controller)
	// helper.InitBot(b_account, b_password)

	config, err := models.GetConfigsDataByName("bot")

	if err == nil && config != nil {
		b_account := config["account"].(float64)
		b_password := config["password"].(string)

		info := bot.Instance
		if info != nil && info.Online {
			callBackResult(&c.Controller, 200, "已登陆，无需重新登陆", nil)
			c.Finish()
			return
		} else {
			fmt.Printf("重新登陆： %s  \n", config)
			// 初始化 Bot
			bot.InitBot(int64(b_account), b_password)

			// 初始化 Modules
			bot.StartService()

			// 使用协议
			// 不同协议可能会有部分功能无法使用
			// 在登陆前切换协议
			bot.UseProtocol(bot.IPad)

			resp, err := bot.Instance.Login()
			// 登录
			if err == nil && resp.Success {
				// 刷新好友列表，群列表
				bot.RefreshList()

				// 将登陆成功的对象加入序列
				// bot.Instances[account.Account] = bot.Instance

				callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
			}
		}
	} else {
		callBackResult(&c.Controller, 200, "error: relogin error code -1", nil)
		c.Finish()
		return
	}

	callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
	c.Finish()
}

func (c *BotController) ApiBotGetInfo() {

	// 要求登陆助理函数
	userAssistant(&c.Controller)
	// helper.InitBot(b_account, b_password)

	// 读取缓存的数据
	config, err := models.GetConfigsDataByName("bot_info")

	// 读取 bot 配置数据
	bot_c, bot_c_err := models.GetConfigsDataByName("bot")

	// 判断账号是否登陆
	info := bot.Instance
	if info == nil {
		if err == nil && config != nil {
			config["on_line"] = false

			if bot_c_err == nil && bot_c != nil {
				config["group_code"] = bot_c["group_code"]
			}

			c.Data["json"] = config
			callBackResult(&c.Controller, 200, "", c.Data["json"])
			c.Finish()
		} else {
			callBackResult(&c.Controller, 200, "bot 未登陆", nil)
			c.Finish()
		}
		return
	}

	// 获取最新的信息数据
	var groupList []interface{}
	for _, item := range info.GroupList {
		new_group := map[string]interface{}{
			"code": item.Code,
			"name": item.Name,
		}
		groupList = append(groupList, new_group)
	}

	botInfo := map[string]interface{}{
		"account":    info.Uin,
		"nick_name":  info.Nickname,
		"avatar":     "https://q2.qlogo.cn/headimg_dl?spec=100&dst_uin=" + strconv.FormatInt(info.Uin, 10),
		"on_line":    info.Online,
		"group_list": groupList,
		"group_code": 0,
	}

	if bot_c_err == nil && bot_c != nil {
		botInfo["group_code"] = bot_c["group_code"]
	}

	// 更新缓存数据
	if err == nil && config != nil {
		models.UpdateConfigByData("bot_info", botInfo)
	} else {
		models.AddConfigByData("bot_info", botInfo)
	}

	c.Data["json"] = botInfo
	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
}

func (c *BotController) ApiBotTest() {

	// 要求登陆助理函数
	userAssistant(&c.Controller)
	// helper.InitBot(b_account, b_password)

	// defer task.StopTask()

	// models.PushTimeAllCourses(1)
	callBackResult(&c.Controller, 200, "ok!!!", nil)
	c.Finish()
}

func (c *BotController) ApiUpdateBotGroupcode() {
	group_code, _ := c.GetInt64("group_code")

	if group_code == 0 {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)
	// helper.InitBot(b_account, b_password)

	config, err := models.GetConfigsDataByName("bot")

	if err == nil && config != nil {
		b_a := config["account"].(float64)
		b_p := config["password"].(string)
		b_p_md5 := config["password_md5"].(string)
		models.UpdateConfigByData("bot", map[string]interface{}{
			"account":      b_a,
			"password":     b_p,
			"password_md5": b_p_md5,
			"group_code":   group_code,
		})
	} else {
		models.AddConfigByData("bot", map[string]interface{}{
			"account":      0,
			"password":     "",
			"password_md5": "",
			"group_code":   group_code,
		})
	}

	callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
	c.Finish()
}

func (c *BotController) ApiBotSubmitCaptcha() {
	code := c.GetString("code")
	// sign := c.GetString("sign")

	var c_sign []string
	c.Ctx.Input.Bind(&c_sign, "sign")

	if code == "" || len(c_sign) == 0 {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		return
	}

	// 要求登陆助理函数
	userAssistant(&c.Controller)

	fmt.Printf("String to byte： %s  \n", c_sign)
	fmt.Printf("String to byte： %b  \n", len(c_sign))

	// helper.BotSubmitCaptcha(code, c_sign)

	callBackResult(&c.Controller, 200, "ok", map[string]interface{}{})
	c.Finish()
}
