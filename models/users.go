package models

import (
	helper "class_notice/helper"

	"github.com/beego/beego/v2/client/orm"
)

type Users struct {
	Id       int
	Name     string `orm:"size(128)"`
	Password string `orm:"size(128)"`
	Token    string `orm:"size(128)"`
	Status   int
}

func init() {
	// orm.RegisterModel(new(Users))
}

// AddUsers insert a new Users into database and returns
// last inserted Id on success.
func AddUsers(m *Users) (user Users, err error) {
	o := orm.NewOrm()
	m.Token = helper.RandStringBytesMaskImprSrc(30)
	if m.Status == 0 {
		// 如果创建时未指定账号状态将默认给启用状态
		m.Status = 1
	}
	id, err := o.Insert(m)
	m.Id = int(id)
	user = Users{
		Id:     m.Id,
		Name:   m.Name,
		Token:  m.Token,
		Status: m.Status,
	}
	return
}

// 账号密码登陆
func LoginUser(name string, password string) (user Users, err error) {

	o := orm.NewOrm()
	user = Users{Name: name, Password: helper.StringToMd5(password)}
	err = o.Read(&user, "name", "password")

	if err == nil {
		user.Token = helper.RandStringBytesMaskImprSrc(30)
		o.Update(&user, "token") // 更新 Token
	}

	return
}

// Token 获取用户
func TokenGetUser(token string) (user Users, err error) {
	o := orm.NewOrm()
	user = Users{Token: token}
	err = o.Read(&user, "token")
	return
}
