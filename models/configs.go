package models

type Configs struct {
	Id   int
	Name string `orm:"size(128)"`
	Data string `orm:"type(longtext)"`
}

func init() {
	// orm.RegisterModel(new(Configs))
}
