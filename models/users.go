package models

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
