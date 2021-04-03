package helper

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"

	_ "github.com/Logiase/MiraiGo-Template/modules/logging"
	asc2art "github.com/yinghau76/go-ascii-art"
)

func init() {
	utils.WriteLogToFS()
	bot.Stop()
	// config.Init()
	// bot.GenRandomDevice()
}

func InitBot(account int64, password string) error {
	// 快速初始化
	// bot.Init()

	bot.InitBot(account, password)

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.IPad)

	// 登录
	resp, err := bot.Instance.Login()
	console := bufio.NewReader(os.Stdin)

	for {
		if err != nil {
			// logger.WithError(err).Fatal("unable to login")
			return err
		}

		var text string
		if !resp.Success {
			switch resp.Error {

			case client.NeedCaptcha:

				// TODO: 登陆需要验证码
				img, _, _ := image.Decode(bytes.NewReader(resp.CaptchaImage))
				fmt.Println(asc2art.New("image", img).Art)

				fmt.Printf("please input captcha (%x) : \n", resp.CaptchaSign)
				fmt.Printf("please input captcha (%s) : \n", convert(resp.CaptchaSign))

				// for text == "" {
				// 	// text, _ = console.ReadString('\n')
				// 	log.Panicln("待输入验证码")
				// 	if text != "" {
				// 		log.Panicln("输入的验证码：" + text)
				// 		break
				// 	}
				// }
				// log.Panicln("输入的验证码：" + text)
				// resp, err = bot.Instance.SubmitCaptcha(strings.ReplaceAll(text, "\n", ""), resp.CaptchaSign)

				return errors.New("bot login: 登陆需要验证码")
				continue

			case client.UnsafeDeviceError:

				// 不安全设备错误
				// fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
				// os.Exit(4)
				return errors.New("bot login: 不安全设备错误")

			case client.SMSNeededError:

				// 需要SMS错误
				fmt.Println("device lock enabled, Need SMS Code")
				fmt.Printf("Send SMS to %s ? (yes)", resp.SMSPhone)
				t, _ := console.ReadString('\n')
				t = strings.TrimSpace(t)
				if t != "yes" {
					// os.Exit(2)
				}
				if !bot.Instance.RequestSMS() {
					// logger.Warnf("unable to request SMS Code")
					log.Panicln("unable to request SMS Code")
					// os.Exit(2)
				}
				// logger.Warn("please input SMS Code: ")
				fmt.Println("please input SMS Code: ")
				text, _ = console.ReadString('\n')
				resp, err = bot.Instance.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))

				return errors.New("bot login: 需要SMS错误")
				continue

			case client.TooManySMSRequestError:

				// 短信请求错误太多
				// fmt.Printf("too many SMS request, please try later.\n")
				// os.Exit(6)
				return errors.New("bot login: 短信请求错误太多")

			case client.SMSOrVerifyNeededError:

				// SMS或验证所需的错误
				// fmt.Println("device lock enabled, choose way to verify:")
				// fmt.Println("1. Send SMS Code to ", resp.SMSPhone)
				// fmt.Println("2. Scan QR Code")
				// fmt.Print("input (1,2):")
				// text, _ = console.ReadString('\n')
				// text = strings.TrimSpace(text)
				// switch text {
				// case "1":
				// 	if !Instance.RequestSMS() {
				// 		fmt.Println("unable to request SMS Code")
				// 		os.Exit(2)
				// 	}
				// 	fmt.Print("please input SMS Code: ")
				// 	text, _ = console.ReadString('\n')
				// 	resp, err = Instance.SubmitSMS(strings.ReplaceAll(strings.ReplaceAll(text, "\n", ""), "\r", ""))
				// 	continue
				// case "2":
				// 	fmt.Printf("device lock -> %v\n", resp.VerifyUrl)
				// 	os.Exit(2)
				// default:
				// 	fmt.Println("invalid input")
				// 	os.Exit(2)
				// }
				return errors.New("bot login: SMS或验证所需的错误")

			case client.SliderNeededError:

				// 需要滑块错误
				if client.SystemDeviceInfo.Protocol == client.AndroidPhone {
					return errors.New("bot login: Android Phone Protocol DO NOT SUPPORT Slide verify")
					// fmt.Println("Android Phone Protocol DO NOT SUPPORT Slide verify")
					// fmt.Println("please use other protocol")
					// os.Exit(2)
				}
				bot.Instance.AllowSlider = false
				bot.Instance.Disconnect()
				resp, err = bot.Instance.Login()

				continue

			case client.OtherLoginError, client.UnknownLoginError:
				// 其他登录错误
				// logger.Fatalf("login failed: %v", resp.ErrorMessage)
				return errors.New("bot login: " + resp.ErrorMessage)
			}
		}

		break
	}

	// 刷新好友列表，群列表
	bot.RefreshList()

	return nil
	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, os.Interrupt, os.Kill)
	// <-ch

	// bot.Stop()

}

func SendGroupMessage(groupCode int64, msg string) {
	m := message.NewSendingMessage().Append(message.NewText(msg))

	if bot.Instance != nil {
		bot.Instance.QQClient.SendGroupMessage(groupCode, m)
	} else {
		log.Println("通知消息发送失败，Bot 账号未登陆！")
	}
}

func BotSubmitCaptcha(code string, sign []uint8) {

	fmt.Printf("String to byte： %s  \n", convert(sign))

	if bot.Instance != nil {
		_, _ = bot.Instance.SubmitCaptcha(code, sign)
	} else {
		log.Println("验证码输入失败，Bot 账号未登陆！")
	}
}

func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}
