package models

import (
	"encoding/json"
	"errors"

	"github.com/beego/beego/v2/client/orm"
)

type Configs struct {
	Id   int
	Name string `orm:"size(128)"`
	Data string `orm:"type(longtext)"`
}

func init() {
	// orm.RegisterModel(new(Configs))
}

// 添加一个 config 文件
func AddConfigs(m *Configs) (c *Configs, err error) {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		c, err = GetConfigsById(int(id))
	}
	return
}

// 更新一个 config 文件
func UpdateConfigs(m *Configs) (c *Configs, err error) {
	o := orm.NewOrm()
	var v Configs
	if m.Name != "" {
		v = Configs{Name: m.Name}
	} else {
		v = Configs{Id: m.Id}
	}
	// ascertain id exists in the database
	if err = o.Read(&v, "Name"); err == nil {
		var id int64
		v.Data = m.Data
		if id, err = o.Update(&v); err == nil {
			return GetConfigsById(int(id))
		}
	}
	return nil, err
}

// 通过 id 获取 config
func GetConfigsById(id int) (c *Configs, err error) {
	o := orm.NewOrm()
	c = &Configs{Id: id}
	if err = o.QueryTable(new(Configs)).Filter("Id", id).RelatedSel().One(c); err == nil {
		return c, nil
	}
	return nil, err
}

// 通过 name 获取 config
func GetConfigsByName(name string) (c *Configs, err error) {
	o := orm.NewOrm()
	c = &Configs{Name: name}
	if err = o.QueryTable(new(Configs)).Filter("Name", name).RelatedSel().One(c); err == nil {
		return c, nil
	}
	return nil, err
}

// 添加一个 config 文件：interface 类型
func AddConfigByData(name string, data map[string]interface{}) (c *Configs, err error) {

	if name == "" || data == nil {
		return nil, errors.New("name or data is nil")
	}

	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	m := &Configs{
		Name: name,
		Data: string(json),
	}

	return AddConfigs(m)
}

// 添加一个 config 文件：Array 类型
func AddConfigByArray(name string, array []interface{}) (c *Configs, err error) {

	if name == "" || array == nil {
		return nil, errors.New("name or data is nil")
	}

	json, err := json.Marshal(array)
	if err != nil {
		return nil, err
	}

	m := &Configs{
		Name: name,
		Data: string(json),
	}

	return AddConfigs(m)
}

// 通过 id 获取 config Data：map[string]interface 类型
func GetConfigsDataById(id int) (data map[string]interface{}, err error) {
	c, err := GetConfigsById(id)
	if err != nil || c == nil {
		return
	}
	dataStr := c.Data
	var data_obj interface{}
	json.Unmarshal([]byte(dataStr), &data_obj)
	if data_obj != nil {
		data = data_obj.(map[string]interface{})
	}
	return
}

// 通过 name 获取 config Data：map[string]interface 类型
func GetConfigsDataByName(name string) (data map[string]interface{}, err error) {
	c, err := GetConfigsByName(name)
	if err != nil || c == nil {
		return
	}
	dataStr := c.Data
	var data_obj interface{}
	json.Unmarshal([]byte(dataStr), &data_obj)
	if data_obj != nil {
		data = data_obj.(map[string]interface{})
	}
	return
}

// 更新一个 config 文件：interface 类型
func UpdateConfigByData(name string, data map[string]interface{}) (c *Configs, err error) {

	if name == "" || data == nil {
		return nil, errors.New("name or data is nil")
	}

	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	m := &Configs{
		Name: name,
		Data: string(json),
	}

	return UpdateConfigs(m)
}

// 添加一个 config 文件：Array 类型
func UpdateConfigByArray(name string, array []interface{}) (c *Configs, err error) {

	if name == "" || array == nil {
		return nil, errors.New("name or data is nil")
	}

	json, err := json.Marshal(array)
	if err != nil {
		return nil, err
	}

	m := &Configs{
		Name: name,
		Data: string(json),
	}

	return UpdateConfigs(m)
}
