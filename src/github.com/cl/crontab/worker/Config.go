package worker

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

// json: 序列化之后的别名
type Config struct {
	EtcdEndpoints     []string      `json:"etcdEndpoints"`
	EtcdDialTimeout   int           `json:"etcdDialTimeout"`
	ScheduleSleepTime time.Duration `json:"scheduleSleepTime"`
	WorkerSleepTime   time.Duration `json:"workerSleepTime"`
	BashDir           string        `json:"bashDir"`
}

var (
	// 单例
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	// 1 读取配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	// 2 做json序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	G_config = &conf
	return
}
