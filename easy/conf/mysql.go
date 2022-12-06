package conf

type MysqlDb struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Database string   `yaml:"database"`
	Tables   []string `yaml:"tables"`
}

func GetMysqlConfig() *Config {
	return &Config{
		Service: Service{
			ProjectName: projectName,
			Version:     version,
		},
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
