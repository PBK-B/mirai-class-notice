package helper

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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

/**
*******************************
*	随机字符串生成
*******************************
 */
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

/**
*******************************
*	md5 获取
*******************************
 */
func StringToMd5(str string) string {
	_str := md5.Sum([]byte(str))
	_md5 := hex.EncodeToString(_str[:])
	return _md5
}

// 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return GetCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		// log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
