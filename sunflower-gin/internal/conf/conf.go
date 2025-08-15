package conf

import (
	"github.com/spf13/viper"
)

// Load 加载配置文件，参数是配置文件的路径
func Load(confPath string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(confPath)

	err := conf.ReadInConfig() // 读取配置信息
	if err != nil {
		panic(err) // 读取配置信息失败时，返回并退出程序
	}
	// conf.GetString("server.mode")
	// conf.GetInt("server.port")

	return conf
}
