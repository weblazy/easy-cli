package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/weblazy/easy-cli/orm/conf"
	"github.com/weblazy/easy-cli/orm/utils"
	"github.com/weblazy/easy-cli/orm/writer"
	"gopkg.in/yaml.v2"
)

// createMysqlYaml 创建 MysqlYaml
func createMysqlYaml(c *cli.Context) error {
	root := c.String("dir")
	mysqlPath := "mysql.yaml"
	if root != "" {
		mysqlPath = root + "/mysql.yaml"
	}
	var writer = writer.NewWriter()
	yamlByte, err := yaml.Marshal([]*conf.MysqlDb{new(conf.MysqlDb)})
	if err != nil {
		return err
	}
	writer.Add(yamlByte)
	writer.WriteToFile(mysqlPath)
	utils.PrintHint("mysql.yaml file has been generated.")
	return nil
}
