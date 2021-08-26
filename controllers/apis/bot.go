package controllers

import (
	"class_notice/helper"
	"class_notice/models"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
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

// 登陆机器人账号（重构）
func (c *BotController) ApiLoginBot() {

	userAssistant(&c.Controller) // 认证

	b_account, _ := c.GetInt64("account")
	b_password := c.GetString("password")

	if b_account == 0 || b_password == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		return
	}

	// 初始化 Bot
	bot.InitBot(b_account, b_password)

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	bot.UseProtocol(bot.IPad)

	// 使用用户自定义的 device.json 文件
	deviceByte := utils.ReadFile(helper.GetCurrentAbPath() + "/device.json")
	if deviceByte == nil {
		deviceByte = utils.ReadFile(helper.GetCurrentAbPath() + "/conf/device.json")
	}
	if deviceByte != nil {
		fmt.Printf("[botLog] use user device.json flie…\n")
		if useDeviceErr := client.SystemDeviceInfo.ReadJson(deviceByte); useDeviceErr != nil {
			fmt.Println("[botLog] device.json error")
		}
	}

	// 登录
	resp, err := bot.Instance.Login()

	for {
		if err != nil {
			// logger.WithError(err).Fatal("unable to login")
			callBackResult(&c.Controller, 200, "QQ 账号登陆异常，"+err.Error(), nil)
			c.Finish()
			return
		}

		if !resp.Success {
			// 登陆失败
			c.Data["json"] = botCallBackToMap(resp)
			callBackResult(&c.Controller, 200, "", c.Data["json"])
			c.Finish()
			return

		} else {
			// 刷新好友列表，群列表
			bot.RefreshList()

			// 返回数据
			c.Data["json"] = map[string]interface{}{}
			callBackResult(&c.Controller, 200, "QQ 账号登陆成功", c.Data["json"])
			c.Finish()

			// 刷新全部机器人账号信息
			RefreshBotInfo()
			return
		}
	}
}

func RefreshBotInfo() {

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

			// 使用用户自定义的 device.json 文件
			deviceByte := utils.ReadFile(helper.GetCurrentAbPath() + "/device.json")
			if deviceByte == nil {
				deviceByte = utils.ReadFile(helper.GetCurrentAbPath() + "/conf/device.json")
			}
			if deviceByte != nil {
				fmt.Printf("[botLog] use user device.json flie…\n")
				if useDeviceErr := client.SystemDeviceInfo.ReadJson(deviceByte); useDeviceErr != nil {
					fmt.Println("[botLog] device.json error")
				}
			}

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
		callBackResult(&c.Controller, 200, "bot 未登陆", nil)
		c.Finish()
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
		// models.UpdateConfigByData("bot_info", botInfo)
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

// 认证机器人账号登陆滑块
func (c *BotController) ApiBotVerifyTicket() {
	userAssistant(&c.Controller) // 认证

	u_ticket := c.GetString("ticket")

	if u_ticket == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	if bot.Instance == nil {
		callBackResult(&c.Controller, 200, "没有待认证的账号", nil)
		c.Finish()
		return
	}

	resp, err := bot.Instance.SubmitTicket(u_ticket)

	if !resp.Success || err != nil {

		if err != nil {
			callBackResult(&c.Controller, 200, "登陆出错，"+err.Error(), nil)
		} else {
			c.Data["json"] = botCallBackToMap(resp)
			callBackResult(&c.Controller, 200, "登陆出错，"+resp.ErrorMessage, c.Data["json"])
		}

		c.Finish()
		return
	}

	c.Data["json"] = ""
	callBackResult(&c.Controller, 200, "成功。", c.Data["json"])
	c.Finish()
}

// 转换登陆结果为 map 数据
func botCallBackToMap(resp *client.LoginResponse) map[string]interface{} {
	switch resp.Error {

	case client.NeedCaptcha:

		err := ioutil.WriteFile("static/img/acptcha/log.jpg", resp.CaptchaImage, os.FileMode(0755))
		if err != nil {
			return map[string]interface{}{
				"error": 10010,
				"text":  "(验证码获取失败) login failed",
			}
		}

		return map[string]interface{}{
			"error": 10011,
			"text":  "(登陆需要验证码) login failed",
			"url":   "/static/img/acptcha/log.jpg",
			"sign":  resp.CaptchaSign,
		}

	case client.UnsafeDeviceError:
		// 不安全设备错误
		return map[string]interface{}{
			"error": 10020,
			"text":  "(不安全设备错误) login failed",
			"url":   resp.VerifyUrl,
		}
	case client.SMSNeededError:

		// 需要SMS错误
		return map[string]interface{}{
			"error": 10030,
			"text":  "(需要短信验证码) login failed",
		}

	case client.TooManySMSRequestError:
		// 短信请求错误太多

		return map[string]interface{}{
			"error": 10040,
			"text":  "(短信请求错误太多) login failed",
		}

	case client.SMSOrVerifyNeededError:
		// SMS或验证所需的错误

		return map[string]interface{}{
			"error": 10050,
			"text":  "(需要短信验证码或扫描二维码) login failed",
			"url":   resp.VerifyUrl,
		}

	case client.SliderNeededError:
		// 需要滑动认证

		return map[string]interface{}{
			"error": 10060,
			"text":  "(需要滑动认证) please look at the doc https://github.com/Mrs4s/go-cqhttp/blob/master/docs/slider.md to get ticket",
			"url":   resp.VerifyUrl,
		}

	case client.OtherLoginError, client.UnknownLoginError:
		// 其他登录错误

		return map[string]interface{}{
			"error": 10070,
			"text":  "(其他登陆错误) login failed: " + resp.ErrorMessage,
		}

	}

	return nil
}
