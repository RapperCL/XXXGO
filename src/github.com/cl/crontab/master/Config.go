package master

import (
	"encoding/json"
	"io/ioutil"
)

// json: 序列化之后的别名
type Config struct {
	ApiPort         int `json:"apiPort"`
	ApiReadTimeout  int `json:"apiReadTimeout"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	// 读取
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	// json反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	// 返回单例配置
	G_config = &conf
	return
}
