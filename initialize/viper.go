package initialize

import (
	"flag"
	"github.com/biningo/boil-gin/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

/**
*@Author lyer
*@Date 2/20/21 15:21
*@Describe
**/

func InitViper() *viper.Viper {
	path := ""
	flag.StringVar(&path, "c", "config.yaml", "choose config file")
	flag.Parse()
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	//监听变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config changed:", e.Name)
		if err := v.Unmarshal(&global.G_CONFIG); err != nil {
			panic(err)
		}
	})

	if err := v.Unmarshal(&global.G_CONFIG); err != nil {
		panic(err)
	}
	return v
}
