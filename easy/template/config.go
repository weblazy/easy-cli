// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/cmd_init.got
// DO NOT EDIT!
package template

import (
	"bytes"
	"fmt"
)

func FromConfigInit(name, pkgs, configStr, configVar, InitMysql, InitRedis string, buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf(`
package config

import (
	%s
    "os"
	"%s/common"
	"github.com/weblazy/easy/econfig/nacos"
	"github.com/weblazy/easy/elog"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/weblazy/easy/econfig/eviper"
	"github.com/weblazy/easy/econfig"
	"github.com/go-redis/redis/v8"
	"github.com/weblazy/easy/db/eredis"
	"github.com/weblazy/easy/db/emysql"
	"github.com/weblazy/easy/http/http_server/http_server_config"
	"github.com/weblazy/easy/db/emysql/emysql_config"
	"github.com/weblazy/easy/db/eredis/eredis_config"
)

type Config struct {
	BaseConfig       struct{}
	%s
	AppMysql         *emysql_config.Config
	AppRedis         *eredis_config.Config
}

var Conf = Config{
	BaseConfig:       struct{}{},
	%s
	AppMysql:         emysql_config.DefaultConfig(),
	AppRedis:         eredis_config.DefaultConfig(),
}

var Redis *eredis.RedisClient

var LocalConfig = ""
			
func InitConf() {
	switch os.Getenv(econfig.EasyConfigType) {
	case econfig.LocalType:
		common.Viper = eviper.NewViperFromString(LocalConfig)
	case econfig.FielType:
		common.Viper = eviper.NewViperFromFile("", os.Getenv(econfig.EasyConfigFile))
	case econfig.NacosType:
		nacos.NewNacosEnv()
		vt := nacos.GetViper()
		vt.SetDataIds(os.Getenv("ServiceName"), os.Getenv("DataId"))
		// 注册配置更新回调
		vt.NacosToViper()
		common.Viper = vt.Viper
	default:
		common.Viper = eviper.NewViperFromString(LocalConfig)
	}
	common.Viper.Unmarshal(&Conf)

	initMysql()
	initRedis()
}

// initMysql 初始化mysql服务
func initMysql() {
	%s
}

// initRedis 初始化redis服务
func initRedis() {
	%s
}`, pkgs,name, configStr, configVar, InitMysql, InitRedis))

}
