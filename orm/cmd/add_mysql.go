package cmd

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/weblazy/easy-cli/orm/conf"
	"github.com/weblazy/easy-cli/orm/file"
	"github.com/weblazy/easy-cli/orm/generate"
	"github.com/weblazy/easy-cli/orm/utils"
	"github.com/weblazy/easy-cli/orm/writer"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var w = writer.NewWriter()

var fileBuffer = new(bytes.Buffer)

// addMysql 添加已有mysql
func addMysql(c *cli.Context) error {
	root := c.String("dir")
	mysqlPath := "mysql.yaml"
	ormPath := "orm.yaml"
	if root != "" {
		mysqlPath = root + "/mysql.yaml"
		ormPath = root + "/orm.yaml"
	}
	mysqlByte, err := os.ReadFile(mysqlPath)
	if err != nil {
		return err
	}
	mysqlDbList := []*conf.MysqlDb{}
	err = yaml.Unmarshal(mysqlByte, &mysqlDbList)
	if err != nil {
		return err
	}

	config, err := InitYaml(root, conf.GetConfig())
	if err != nil {
		return err
	}
	m1 := make(map[string]conf.Mysql)
	for _, v := range config.MysqlList {
		m1[v.Name] = v
	}
	for _, mysqlDb := range mysqlDbList {
		mysql, err := generate.Genertate(mysqlDb, config)
		if err != nil {
			return err
		}
		m1[mysql.Name] = *mysql
	}
	list := make([]conf.Mysql, 0)
	for _, v := range m1 {
		list = append(list, v)
	}
	config.MysqlList = list
	_, err = CreateYaml(ormPath, config)
	if err != nil {
		return err
	}
	utils.PrintHint("orm.yaml file has been generated.")
	return nil
}

// CreateYaml 创建Yaml文件
func CreateYaml(yamlPath string, config *conf.Config) (*conf.Config, error) {
	var writer = writer.NewWriter()
	yamlByte, err := yaml.Marshal(config)
	if err != nil {
		return config, err
	}
	writer.Add(yamlByte)
	writer.ForceWriteToFile(yamlPath)
	return config, nil
}

// InitYaml 初始化Yaml配置文件
func InitYaml(yamlPath string, config *conf.Config) (*conf.Config, error) {
	if file.CheckFileIsExist(yamlPath) {
		apiFile, err := os.Open(yamlPath)
		if err == nil {
			content, err := ioutil.ReadAll(apiFile)
			if err != nil {
				panic(err)
			}
			cfg := conf.Config{}
			err = yaml.Unmarshal(content, &cfg)
			if err != nil {
				panic(err)
			}
			return &cfg, nil
		}
		panic(err)
	}

	return CreateYaml(yamlPath, config)
}
