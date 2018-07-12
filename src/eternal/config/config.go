package config

import (
	"github.com/spf13/viper"
)

var DEBUG = true

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
	DEBUG = viper.GetBool("debug")
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetStringDefault(key, def string) string {
	viper.SetDefault(key, def)
	return viper.GetString(key)
}

func GetIntDefault(key string, def int) int {
	viper.SetDefault(key, def)
	return viper.GetInt(key)
}

func GetStringSliceDefault(key string, def []string) []string {
	viper.SetDefault(key, def)
	return viper.GetStringSlice(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func GetBoolDefault(key string, def bool) bool {
	viper.SetDefault(key, def)
	return viper.GetBool(key)
}
