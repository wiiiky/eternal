package config

import (
	"github.com/spf13/viper"
)

/*
 * 读取配置
 * https://github.com/spf13/viper
 */
func Init(APPNAME string) {
	viper.SetConfigName(APPNAME)             // name of config file (without extension)
	viper.AddConfigPath("/etc/" + APPNAME)   // path to look for the config file in
	viper.AddConfigPath("$HOME/." + APPNAME) // call multiple times to add many search paths
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(err)
	}
	viper.SetDefault("debug", true)
}
