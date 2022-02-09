package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"

	helper "class_notice/helper"
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

func GetUserById(id int) (v *Users, err error) {
	o := orm.NewOrm()
	v = &Users{Id: id}
	if err = o.QueryTable(new(Users)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
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

// UpdateUsers updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *Users) (err error) {
	o := orm.NewOrm()
	v := Users{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
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
func TokenGetUser(token string) (user *Users, err error) {
	o := orm.NewOrm()
	user = &Users{Token: token}
	err = o.Read(user, "token")
	if err != nil {
		return nil, err
	}

	return
}

func AllUser(limit int, page int) (users []Users, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Users))
	_, err = qs.Filter("id__isnull", false).Limit(limit, page).All(&users)
	return
}
