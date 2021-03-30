package controllers

import (
	"class_notice/helper"
	"class_notice/models"
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

	callBackResult(&c.Controller, 200, "ok", nil)
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
			go helper.InitBot(int64(b_account), b_password)
		}
	} else {
		callBackResult(&c.Controller, 200, "error: relogin error code -1", nil)
		c.Finish()
		return
	}

	callBackResult(&c.Controller, 200, "ok", nil)
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
			callBackResult(&c.Controller, 403, "bot 未登陆", nil)
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
