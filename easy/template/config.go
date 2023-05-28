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
}

var Conf = Config{
	BaseConfig:       struct{}{},
	%s
}

var LocalConfig = ""
			
`, pkgs,name, configStr, configVar))

}
