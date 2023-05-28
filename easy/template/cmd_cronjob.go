// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/cmd_cronjob.got
// DO NOT EDIT!
package template

import (
	"bytes"
	"fmt"
)

func FromCmdCronJob(name, cronjobs string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package cronjobs

import (
	`)
	buffer.WriteString(fmt.Sprintf("\"%s/cronjobs/config\"\n",name))
	buffer.WriteString(fmt.Sprintf("\"%s/cronjobs/handler\"\n",name)+
	`"github.com/robfig/cron/v3"
	"github.com/weblazy/easy/closes"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:    "cron",
	Aliases: []string{"c"},
	Usage:   "cron start",
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
	// 初始化必要内容
	econfig.InitGlobalViper(&config.Conf, config.LocalConfig)
	cronJob := cron.New()

    `)
	buffer.WriteString(cronjobs)
	buffer.WriteString(`

	cronJob.Start()

	closes.AddShutdown(closes.ModuleClose{
		Name:     "CronTable",
		Priority: 0,
		Func: func() {
			cronJob.Stop()
		},
	})
	closes.SignalClose()
	return nil
}`)

}
