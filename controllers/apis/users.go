package controllers

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	"class_notice/helper"
	"class_notice/models"
)

type UsersController struct {
	beego.Controller
}

// 用户登陆后台接口
func (c *UsersController) ApiLogin() {

	u_name := c.GetString("name")
	u_password := c.GetString("password")

	if u_name == "" || u_password == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		return
	}

	// TODO: Bin BY 这里应该还可以做判断用户名和密码是否合法
	user, err := models.LoginUser(u_name, u_password)

	if user.Id == 0 || err == orm.ErrNoRows {
		callBackResult(&c.Controller, 200, "用户名或密码错误", nil)
		return
	} else if err == orm.ErrMissPK {
		callBackResult(&c.Controller, 403, "服务器错误", nil)
		return
	} else {
		if user.Status == 0 {
			callBackResult(&c.Controller, 200, "账号被禁用", nil)
			return
		}

		c.Data["json"] = map[string]interface{}{
			"id":    user.Id,
			"name":  user.Name,
			"token": user.Token,
		}
		c.SetSecureCookie("bin", "u_token", user.Token)
		callBackResult(&c.Controller, 200, "", c.Data["json"])

		return
	}
}

// 获取后台用户信息接口
func (c *UsersController) ApiGetMe() {

	_, user, _ := userAssistant(&c.Controller)

	c.Data["json"] = map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"token": user.Token,
	}
	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
}

// 创建后台管理用户接口
func (c *UsersController) ApiCreateUser() {

	// 要求登陆助理函数
	_, me, _ := userAssistant(&c.Controller)

	u_name := c.GetString("name")
	u_password := c.GetString("password")

	if u_name == "" || u_password == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	if !isAdminUser(me) {
		callBackResult(&c.Controller, 403, "权限不足", nil)
		c.Finish()
		return
	}

	// TODO: Bin BY 这里应该还可以做判断用户名和密码是否合法
	user := models.Users{Name: u_name, Password: helper.StringToMd5(u_password), Status: 1}
	user, err := models.AddUsers(&user)

	if err != nil {
		callBackResult(&c.Controller, 403, "用户创建失败"+err.Error(), nil)
		c.Finish()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"token": user.Token,
	}
	callBackResult(&c.Controller, 200, "", c.Data["json"])
}

func (c *UsersController) ApiUpStatusUser() {
	// 要求登陆助理函数
	_, me, _ := userAssistant(&c.Controller)

	u_id, _ := c.GetInt("id", 0)

	if u_id == 0 {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	// 防止全部用户被禁用状态导致系统进入假锁死
	if u_id == 1 {
		callBackResult(&c.Controller, 200, "默认 admin 用户不能禁用！", nil)
		c.Finish()
		return
	}

	if !isAdminUser(me) {
		callBackResult(&c.Controller, 403, "权限不足", nil)
		c.Finish()
		return
	}

	user, err := models.GetUserById(u_id)

	if user == nil || err != nil {
		callBackResult(&c.Controller, 200, "用户不存在", nil)
		c.Finish()
		return
	}

	if user.Status == 1 {
		// 禁用用户
		user.Status = 0
		err = models.UpdateUsersById(user)
	} else {
		// 启用用户
		user.Status = 1
		err = models.UpdateUsersById(user)
	}

	if err != nil {
		callBackResult(&c.Controller, 200, "出错啦，"+err.Error(), nil)
		c.Finish()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"token": user.Token,
	}

	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
}

// 修改用户密码接口
func (c *UsersController) ApiUpdateUser() {
	// 要求登陆助理函数
	_, me, _ := userAssistant(&c.Controller)

	u_id, _ := c.GetInt("id", 0)
	u_password := c.GetString("password")

	if u_id == 0 || u_password == "" {
		callBackResult(&c.Controller, 403, "参数异常", nil)
		c.Finish()
		return
	}

	if !isAdminUser(me) && me.Id != u_id {
		callBackResult(&c.Controller, 403, "权限不足", nil)
		c.Finish()
		return
	}

	user, err := models.GetUserById(u_id)

	if user == nil || err != nil {
		callBackResult(&c.Controller, 200, "用户不存在", nil)
		c.Finish()
		return
	}

	// TODO: Bin BY 这里应该还可以做判断用户名和密码是否合法
	user.Password = helper.StringToMd5(u_password)

	err = models.UpdateUsersById(user)

	if err != nil {
		callBackResult(&c.Controller, 200, "出错啦，"+err.Error(), nil)
		c.Finish()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"token": user.Token,
	}

	callBackResult(&c.Controller, 200, "", c.Data["json"])
	c.Finish()
}

// 获取用户列表接口
func (c *UsersController) ApiUserList() {
	// 要求登陆助理函数
	userAssistant(&c.Controller)
	u_count, _ := c.GetInt("count", 10)
	u_page, _ := c.GetInt("page", 0)

	users, err := models.AllUser(u_count, u_page)

	if err != nil {
		callBackResult(&c.Controller, 403, "服务器错误", nil)
		c.Finish()
		return
	}

	var new_users []interface{}

	for item := range users {
		i_u := users[item]
		new_u := map[string]interface{}{
			"id":     i_u.Id,
			"name":   i_u.Name,
			"status": i_u.Status,
		}
		new_users = append(new_users, new_u)
	}

	callBackResult(&c.Controller, 200, "", new_users)
}

// 做一个简单的判断用户是否属于超级管理员
func isAdminUser(u models.Users) bool {
	// TODO: 目前简单判断 ID 为 1 的用户为超级管理员
	if u.Id == 1 {
		return true
	} else {
		return false
	}
}
