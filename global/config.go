package global

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局配置
var Conf = &AppConf{}

type AppConf struct {
	// mapstructure：通用结构体tag
	Sensitive    []string `mapstructure:"sensitive"`
	TokenSecret  string   `mapstructure:"token_secret"`
	OfflineNum   int      `mapstructure:"offline_num"`
	MessageQueue int      `mapstructure:"message_queue"`
}

// Config init
func configInit() {
	// 设置配置文件名称
	viper.SetConfigName("config")
	// 查找配置文件 的 路径
	viper.AddConfigPath(RootDir + "/config/")
	//读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		//读取配置信息错误
		log.Printf("viper.ReadInConfig() failed: %v\n", err)
		return
	}
	// 把读取到的信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		log.Printf("viper.Unmarshal() failed: %v\n", err)
		return
	}
	viper.WatchConfig() // 实时监控配置文件
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("配置文件发生修改...")
		// 重新反序列化
		if err := viper.Unmarshal(Conf); err != nil {
			log.Printf("viper.Unmarshal() failed: %v\n", err)
			return
		}
	})
}
