package main

import (
	_ "class_notice/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

