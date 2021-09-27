/*
 * @Author: Bin
 * @Date: 2021-09-27
 * @FilePath: /class_notice/plugins/system.go
 */
package plugins

import (
	"fmt"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

func init() {
	bot.RegisterModule(SystemInstance)
}

var SystemInstance = &system{}

var logger = utils.GetModuleLogger("bin.system")

// var tem map[string]string

type system struct {
}

func (mov *system) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "bin.system",
		Instance: SystemInstance,
	}
}

func (mov *system) Init() {
	fmt.Printf("【Bot Plugins】模块初始化=%+v\n", mov.MiraiGoModule().ID)
	// path := config.GlobalConfig.GetString("logiase.autoreply.path")

	// if path == "" {
	// 	path = "./autoreply.yaml"
	// }

	// bytes := utils.ReadFile(path)
	// err := yaml.Unmarshal(bytes, &tem)
	// if err != nil {
	// 	logger.WithError(err).Errorf("unable to read config file in %s", path)
	// }
}

func (mov *system) PostInit() {
}

func (mov *system) Serve(b *bot.Bot) {

	// b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {

	// 	fmt.Printf("message=%+v\n", msg.Sender.Nickname)

	// 	if botObj := bot.Instance; botObj == nil {
	// 		// 机器人已下线，直接结束回复流程
	// 		fmt.Println("【收到消息】机器人已下线，直接结束回复流程")
	// 		return
	// 	}

	// 	groupInfo := c.FindGroup(msg.GroupCode)
	// 	if groupInfo == nil {
	// 		// QQ 群信息获取失败，结束流程
	// 		fmt.Println("【收到消息】QQ 群信息获取失败，直接结束回复流程")
	// 		return
	// 	}
	// 	groupMemberInfo := groupInfo.FindMember(c.Uin)
	// 	botName := ""
	// 	if groupMemberInfo == nil {
	// 		// QQ 群我的数据获取失败，直接赋值昵称
	// 		botName = b.Nickname
	// 	} else {
	// 		botName = groupMemberInfo.DisplayName()
	// 	}

	// 	fmt.Printf("群昵称=%+v\n", botName)

	// 	fmt.Println("【收到消息】" + msg.ToString())

	// })

	b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {

		// fmt.Printf("message=%+v\n", msg.ToString())

		if botObj := bot.Instance; botObj == nil {
			// 机器人已下线，直接结束回复流程
			fmt.Println("【收到消息】机器人已下线，直接结束回复流程")
			return
		}

		if msg.ToString() == "\\help" {
			m := message.NewSendingMessage().Append(message.NewText(`【指令菜单】
\help (帮助菜单)
\info (系统状态)
`))
			c.SendPrivateMessage(msg.Sender.Uin, m)
			return
		}

		if msg.ToString() == "\\info" {
			m := message.NewSendingMessage().Append(message.NewText("系统状态: 正常\n机器人状态: 在线\n"))
			c.SendPrivateMessage(msg.Sender.Uin, m)
			return
		}

	})

	// b.OnPrivateMessage(func(c *client.QQClient, msg *message.PrivateMessage) {
	// 	out := autoreply(msg.ToString())
	// 	if out == "" {
	// 		return
	// 	}
	// 	m := message.NewSendingMessage().Append(message.NewText(out))
	// 	c.SendPrivateMessage(msg.Sender.Uin, m)
	// })
}

func (mov *system) Start(bot *bot.Bot) {
}

func (mov *system) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}
