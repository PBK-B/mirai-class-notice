/*
 * @Author: Bin
 * @Date: 2021-09-29
 * @FilePath: /class_notice/controllers/apis/plugin.go
 */
package controllers

import (
	"class_notice/helper"
	"class_notice/models"
	"io"
	"os"
	"path/filepath"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
)

type PluginController struct {
	beego.Controller
}

// 插件上传接口
func (c *PluginController) ApiUploadPlugin() {

	// 要求登陆助理函数
	// userAssistant(&c.Controller)

	file, head, err := c.GetFile("plugin")
	if err != nil {
		callBackResult(&c.Controller, 200, "插件上传失败。", nil)
		return
	}
	defer file.Close()
	if file == nil || head == nil {
		callBackResult(&c.Controller, 200, "插件包异常。", nil)
	}

	// err = helper.UnzipPlugin(file)
	// context, err := helper.TarReaderFile("./manifest.json", file)
	// if err != nil || context == nil {
	// 	callBackResult(&c.Controller, 200, "插件包损坏。"+err.Error(), nil)
	// }

	// var v models.PluginManifest
	// if err := json.Unmarshal(context, &v); err != nil {
	// 	callBackResult(&c.Controller, 200, "插件包配置读取失败。"+err.Error(), nil)
	// }

	// 缓存插件安装包文件
	var cache_dir_str = "files/cache/"
	err = os.MkdirAll(cache_dir_str, os.ModePerm) // 创建一个插件存放文件夹
	cache_dir_str = filepath.Join(helper.GetCurrentAbPath(), cache_dir_str)
	cache_path_id, err := uuid.NewUUID() // 生成一个路径的 uuid
	cache_path_dir := cache_dir_str + "/" + cache_path_id.String()
	f, err := os.OpenFile(cache_path_dir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		callBackResult(&c.Controller, 200, err.Error(), nil)
	}
	io.Copy(f, file)
	defer f.Close()

	// 安装插件
	plugin, err := models.InstallPlugin(cache_path_dir)
	// 清空缓存文件夹
	os.RemoveAll(cache_dir_str)

	if plugin == nil || err != nil {
		callBackResult(&c.Controller, 200, err.Error(), nil)
	}

	callBackResult(&c.Controller, 200, "", plugin.ToMap())
	c.Finish()
}
