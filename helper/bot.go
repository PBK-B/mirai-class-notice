package helper

import (
	"os"
	"os/signal"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/message"

	_ "github.com/Logiase/MiraiGo-Template/modules/logging"
)

func init() {
	utils.WriteLogToFS()
	bot.Stop()
	// config.Init()
	// bot.GenRandomDevice()
}

func InitBot(account int64, password string) {
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
	bot.Login()

	// 刷新好友列表，群列表
	bot.RefreshList()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}

func SendGroupMessage(groupCode int64, msg string) {
	m := message.NewSendingMessage().Append(message.NewText(msg))
	bot.Instance.QQClient.SendGroupMessage(groupCode, m)
}
