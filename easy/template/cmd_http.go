// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/cmd_api.got
// DO NOT EDIT!
package template

import (
	"bytes"
	"fmt"
)

func FromCmdApi(serviceName, homePath string, buffer *bytes.Buffer) {
	buffer.WriteString(fmt.Sprintf(`
package %s

import (
	"%s/routes"
	"%s/config"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/weblazy/easy/econfig/eviper"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/weblazy/easy/closes"
	"github.com/urfave/cli/v2"
	"github.com/weblazy/easy/http/http_server"
	"github.com/weblazy/easy/http/http_server/http_server_config"
)

var Cmd = &cli.Command{
	Name:    "api",
	Aliases: []string{"a"},
	Usage:   "api start",
	Subcommands: []*cli.Command{
		{
			Name:   "start",
			Usage:  "开启运行api服务",
			Action: Run,
		},
	},
}

func Run(c *cli.Context) error {
	defer closes.Close()
	econfig.InitGlobalViper(&config.Conf, config.LocalConfig)

	s, err := http_server.NewHttpServer(config.Conf.HttpServerConfig)
	if err != nil {
		return err
	}
	// 注册路由
	routes.Routes(s.Engine)

	err = s.Start()
	if err != nil {
		return err
	}
	return nil
}`, serviceName, homePath, homePath))

}
