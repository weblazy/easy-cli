package conf

type MysqlDb struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Database string   `yaml:"database"`
	Tables   []string `yaml:"tables"`
}

type Config struct {
	ProjectName string  `yaml:"project_name"`
	MysqlList   []Mysql `yaml:"mysql"`
}

type Mysql struct {
	Name   string  `yaml:"name"` // Mysql名称，默认default
	Models []Model `yaml:"models"`
}

type Model struct {
	Name    string   `yaml:"name"`   // 表名
	Auto    bool     `yaml:"auto"`   // 是否自动创建表结构
	Fields  []string `yaml:"fields"` // 字段列表
	Comment string   `yaml:"comment"`
}

func GetConfig() *Config {
	return &Config{
		ProjectName: "gorm",
		MysqlList: []Mysql{
			{
				Name: "app",
				Models: []Model{
					{
						Name: "user",
						Auto: false,
						Fields: []string{
							"column:id;primary_key;type:int AUTO_INCREMENT",
							"column:name;type:varchar(100) NOT NULL;default:'';comment:'用户名';unique_index",
						},
						Comment: "用户表",
					},
				},
			},
		},
	}
}
