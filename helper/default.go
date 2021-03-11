package helper

import (
	beego "github.com/beego/beego/v2/server/web"
)

/**
*******************************
*	环境变量获取
*******************************
 */

func Env(key string) (string, error) {
	return beego.AppConfig.String(key)
}

func EnvInt(key string) (int, error) {
	return beego.AppConfig.Int(key)
}
func EnvInt64(key string) (int64, error) {
	return beego.AppConfig.Int64(key)
}
func EnvFloat(key string) (float64, error) {
	return beego.AppConfig.Float(key)
}
func EnvBool(key string) (bool, error) {
	return beego.AppConfig.Bool(key)
}

/**
*******************************
*	MySQL 连接
*******************************
 */

func GetConnectionURL(username string, password string, host string, database string) string {
	// root:xxx@tcp(127.0.0.1:3306)/xxx?charset=utf8
	return username + ":" + password + "@tcp(" + host + ")/" + database + "?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
}
