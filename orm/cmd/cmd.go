package cmd

import (
	"github.com/urfave/cli/v2"
)

// Mysql 生成脚手架
var Mysql = &cli.Command{
	Name:  "mysql",
	Usage: "mysql",
	Subcommands: []*cli.Command{
		{
			Name: "create_yaml",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dir",
					Aliases:     []string{"d"},
					Usage:       "dir",
					DefaultText: ".",
				}},
			Usage:  "mysql create_yaml -d .",
			Action: createMysqlYaml,
		},
		{
			Name: "add_mysql",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dir",
					Aliases:     []string{"d"},
					Usage:       "dir",
					DefaultText: ".",
				}},
			Usage:  "mysql add_mysql -d .",
			Action: addMysql,
		},
		{
			Name: "create_code",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dir",
					Aliases:     []string{"d"},
					Usage:       "dir",
					DefaultText: ".",
				}},
			Usage:  "mysql create_code -d .",
			Action: creatCode,
		},
	},
}
