package models

import (
	"class_notice/helper"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strconv"

	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Plugin struct {
	Id      int
	Name    string `orm:"size(128)"`
	Package string `orm:"size(128)"`
	Author  string `orm:"size(128)"`
	Website string `orm:"size(128)"`
	Path    string `orm:"size(128)"`
	Config  string `orm:"type(text)"`
	Version float64
	States  int
	Logs    string `orm:"type(text)"`
}

type PluginManifest struct {
	Id      string
	Name    string
	Version float64
	Author  string
	Website string
	Config  interface{}
}

func init() {
	// orm.RegisterModel(new(Plugin))
}

// 多字段索引
func (u *Plugin) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "Package"},
	}
}

// Plugin 更新
func (p *Plugin) UpdateInfo(np *Plugin) (err error) {
	np.Id = p.Id
	if np.Path != "" {
		p.Path = np.Path
	}
	if np.Package != "" {
		p.Package = np.Package
	}
	if np.Name != "" {
		p.Name = np.Name
	}
	if np.Author != "" {
		p.Author = np.Author
	}
	if np.Website != "" {
		p.Website = np.Website
	}
	if np.Version != p.Version {
		p.Version = np.Version
	}
	return UpdatePluginById(p)
}

// Plugin 更新配置
func (p *Plugin) UpdateConfig(config interface{}) (err error) {
	json, err := json.Marshal(config)
	if err != nil {
		return err
	}
	p.Config = string(json)
	return UpdatePluginById(p)
}

// Plugin 写入日志
func (p *Plugin) AddLog(log string) (err error) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	p.Logs = p.Logs + "\n" + "[" + timeStr + "] " + log
	return UpdatePluginById(p)
}

// Plugin 转 map 数据
func (p *Plugin) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      p.Id,
		"name":    p.Name,
		"version": p.Version,
		"package": p.Package,
		"logs":    p.Logs,
		"author":  p.Author,
		"website": p.Website,
		"states":  p.States,
	}
}

// Plugin 注册到 MiraiGo 中
func (p *Plugin) LoadRegisterModule() (err error) {
	// 获取插件 SO 文件绝对路径
	var pluginPathSO = filepath.Join(helper.GetCurrentAbPath(), p.Path+"/build.so")

	// 加载 so 文件，调用 RegisterModule 方法注册插件
	plu, err := plugin.Open(pluginPathSO)
	if plu == nil || err != nil {
		p.AddLog("加载插件 so 文件失败，" + err.Error())
		return
	}
	registerModule, err := plu.Lookup("RegisterModule")
	if err != nil {
		p.AddLog("加载插件 RegisterModule 方法未实现，" + err.Error())
		return
	}
	registerModule.(func())()
	p.AddLog("插件 (" + p.Package + ") 注册成功。")
	return
}

// AddPlugin insert a new Plugin into database and returns
// last inserted Id on success.
func AddPlugin(m *Plugin) (p *Plugin, err error) {
	o := orm.NewOrm()
	id, err := o.Insert(m)
	if err == nil {
		p, err = GetPluginById(int(id))
	}
	return
}

// GetPluginById retrieves Plugin by Id. Returns error if
// Id doesn't exist
func GetPluginById(id int) (v *Plugin, err error) {
	o := orm.NewOrm()
	v = &Plugin{Id: id}
	if err = o.QueryTable(new(Plugin)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPlugin retrieves all Plugin matches certain condition. Returns empty list if
// no records exist
func GetAllPlugin(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Plugin))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Plugin
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePlugin updates Plugin by Id and returns error if
// the record to be updated doesn't exist
func UpdatePluginById(m *Plugin) (err error) {
	o := orm.NewOrm()
	v := Plugin{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		// var num int64
		// if num, err = o.Update(m); err == nil {
		// 	fmt.Println("Number of records updated in database:", num)
		// }
		_, err = o.Update(m)
	}
	return
}

// DeletePlugin deletes Plugin by Id and returns error if
// the record to be deleted doesn't exist
func DeletePlugin(id int) (err error) {
	o := orm.NewOrm()
	v := Plugin{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Plugin{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// 获取全部插件
func AllPlugin(limit int, page int) (plugin []Plugin, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(&Plugin{})
	_, err = qs.Filter("id__isnull", false).Limit(limit, page).All(&plugin)
	return
}

// 安装插件文件
func InstallPlugin(file string) (p *Plugin, err error) {

	// 读取插件包 manifest 信息
	manifestStr, err := helper.TarReaderFile("./manifest.json", file)
	if err != nil || manifestStr == nil {
		return p, errors.New("插件包损坏。")
	}
	var manifestObj PluginManifest
	if err := json.Unmarshal(manifestStr, &manifestObj); err != nil {
		return p, errors.New("插件包 manifest 解析失败。")
	}

	// 判断插件是否安装过
	p, _ = GetPluginByPackage(manifestObj.Id)
	if p != nil {
		// 插件存在，判断是否需要覆盖安装
		if p.Version > manifestObj.Version {
			return p, errors.New("插件存在已安装的更高版本。")
		} else {
			// 记录旧版本插件安装路径
			var oldPluginPath = p.Path + ""

			// 解压插件，获得插件升级安装目录
			pluginPath, err := helper.UnzipPlugin(file)
			if err != nil || pluginPath == "" {
				return p, errors.New("插件包升级安装异常。")
			}

			err = p.UpdateInfo(&Plugin{
				Id:      p.Id,
				Package: manifestObj.Id,
				Name:    manifestObj.Name,
				Author:  manifestObj.Author,
				Website: manifestObj.Website,
				Version: manifestObj.Version,
				Path:    pluginPath,
			})

			if err != nil {
				// 失败，删除插件
				os.RemoveAll(pluginPath)
				return p, errors.New("插件包升级安装失败。")
			}

			// fmt.Println("oldPluginPath: " + oldPluginPath)
			// fmt.Println("p.Path: " + p.Path)
			// fmt.Println("pluginPath: " + pluginPath)

			// 成功，删除旧版本插件
			os.RemoveAll(oldPluginPath)

			p.AddLog("插件升级安装 (" + strconv.FormatFloat(manifestObj.Version, 'f', -1, 64) + ") 成功！")

			// 注册插件
			p.LoadRegisterModule()

			return p, nil
		}
	}

	// 解压插件，获得插件安装目录
	pluginPath, err := helper.UnzipPlugin(file)
	if err != nil || pluginPath == "" {
		return p, errors.New("插件包安装异常。")
	}

	p = &Plugin{
		Package: manifestObj.Id,
		Name:    manifestObj.Name,
		Author:  manifestObj.Author,
		Website: manifestObj.Website,
		Version: manifestObj.Version,
		States:  0,
		Path:    pluginPath,
	}
	p, err = AddPlugin(p)
	if err != nil || p == nil {
		// 失败，删除插件
		os.RemoveAll(pluginPath)
		return p, errors.New("插件包安装失败。")
	}
	p.AddLog("插件安装成功！")

	// 配置默认数据
	if manifestObj.Config != nil {
		cErr := p.UpdateConfig(manifestObj.Config)
		if cErr == nil {
			p.AddLog("插件默认数据配置成功！")
		}
	}

	// 注册插件
	p.LoadRegisterModule()

	return
}

// 通过包名获取 Plugin
func GetPluginByPackage(pack string) (c *Plugin, err error) {
	o := orm.NewOrm()
	c = &Plugin{Name: pack}
	if err = o.QueryTable(new(Plugin)).Filter("Package", pack).RelatedSel().One(c); err == nil {
		return c, nil
	}
	return nil, err
}
