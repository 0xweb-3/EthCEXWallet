package initizlize

import (
	"fmt"
	"github.com/0xweb-3/EthCEXWallet/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() {
	v := viper.New()
	fileName := "../config/conf/config.yaml"

	v.SetConfigFile(fileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	fmt.Println(global.ServerConfig.Name)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		err := v.ReadInConfig()
		if err != nil {
			panic(err)
		}
		if err := v.Unmarshal(&global.ServerConfig); err != nil {
			panic(err)
		}
	})
}
