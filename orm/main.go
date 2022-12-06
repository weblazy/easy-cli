package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/weblazy/easy-cli/orm/cmd"
	"github.com/weblazy/easy-cli/orm/utils"
)

// 配置常量
const (
	PROJECT_NAME    = "gorm"
	PROJECT_VERSION = "v1.0.0"
)

func main() {
	// 打印banner
	utils.PrintBanner(PROJECT_NAME)
	// 配置cli参数
	app := cli.NewApp()
	app.Name = PROJECT_NAME
	app.Usage = PROJECT_NAME
	app.Version = PROJECT_VERSION
	// 指定命令运行的函数
	app.Commands = []*cli.Command{
		cmd.Mysql,
	}
	// 启动cli
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
