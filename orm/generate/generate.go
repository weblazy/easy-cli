package generate

import (
	"strings"

	"github.com/weblazy/easy-cli/orm/conf"
	"github.com/weblazy/easy-cli/orm/def"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Genertate(mysqlDb *conf.MysqlDb, config *conf.Config) (*conf.Mysql, error) {
	gormDb, err := OpenORM(mysqlDb)
	if err != nil {
		return nil, err
	}

	tableNames := []string{}
	for k1 := range mysqlDb.Tables {
		tableNames = append(tableNames, mysqlDb.Tables[k1])
	}

	//生成所有表信息
	tables := getTables(gormDb, mysqlDb.Database, tableNames)
	models := make([]conf.Model, 0)
	for _, table := range tables {
		fileds := make([]string, 0)
		fieldList := getFields(gormDb, table.Name)
		for _, v1 := range fieldList {
			field := "column:" + v1.Field + ";type:" + v1.Type
			if v1.Null == "NO" {
				field += "  NOT NULL"
			}
			field += ";default:" + v1.Default + " " + v1.Extra + ";" + "comment:" + v1.Comment + ";"
			fileds = append(fileds, field)
		}
		model := conf.Model{
			Name:    table.Name,
			Comment: table.Comment,
			Fields:  fileds,
		}
		models = append(models, model)
	}
	databaseName := strings.ReplaceAll(mysqlDb.Database, "-", "_")
	mysql := conf.Mysql{
		Name:   databaseName,
		Models: models,
	}
	return &mysql, nil
}

func OpenORM(mysqlDb *conf.MysqlDb) (*gorm.DB, error) {
	dsn := mysqlDb.User + ":" + mysqlDb.Password + "@tcp(" + mysqlDb.Host + ":" + mysqlDb.Port + ")/" + mysqlDb.Database + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, err
}

// 获取表信息
func getTables(gormDb *gorm.DB, databaseName string, tableNames []string) []def.Table {
	var tables []def.Table
	if len(tableNames) == 0 {
		gormDb.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema=?;", databaseName).Find(&tables)
	} else {
		gormDb.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (?) AND table_schema=?;", tableNames, databaseName).Find(&tables)
	}
	return tables
}

// 获取所有字段信息
func getFields(gormDb *gorm.DB, tableName string) []def.Field {
	var fields []def.Field
	gormDb.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}
