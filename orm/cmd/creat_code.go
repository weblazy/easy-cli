package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/weblazy/easy-cli/easy/conf"
	"github.com/weblazy/easy-cli/orm/def"
	"github.com/weblazy/easy-cli/orm/file"
	"github.com/weblazy/easy-cli/orm/utils"
)

// creatCode 创建服务并创建初始化配置
func creatCode(c *cli.Context) error {
	config := conf.GetConfig()
	yamlPath := c.String("config")
	root := "."

	if yamlPath == "" {
		yamlPath = root + "/orm.yaml"
	}

	if !file.CheckFileIsExist(yamlPath) {
		return fmt.Errorf("%s is not found", yamlPath)
	}

	// 创建配置&读取配置
	config, err := InitYaml(yamlPath, config)
	if err != nil {
		panic(err)
	}

	modPath := root + "/go.mod"
	if file.CheckFileIsExist(modPath) {
		resp, err := utils.Cmd("go", []string{"fmt", "./..."})
		if err != nil {
			fmt.Println(resp)
			panic(err)
		}
	} else {
		utils.PrintHint("Run go mod init.")
		resp, err := utils.Cmd("go", []string{"mod", "init", config.Service.ProjectName})
		if err != nil {
			fmt.Println(resp)
			panic(err)
		}
	}

	CreateModel(root, config.Service.ProjectName, config.MysqlList)

	utils.PrintHint("Run go mod tidy.")

	resp, err := utils.Cmd("go", []string{"mod", "tidy"})
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	utils.PrintHint("Run go fmt.")
	resp, err = utils.Cmd("go", []string{"fmt", "./..."})
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	utils.PrintHint("goimports -l -w .")
	resp, err = utils.Cmd("goimports", []string{"-l", "-w", "."})
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}
	utils.PrintHint("Welcome to orm, the project has been initialized.")

	return nil
}

func CreateModel(root, projectName string, mysqlList []conf.Mysql) {

	initDb := ""
	for _, v1 := range mysqlList {
		dir := root + "/model/" + v1.Name
		initDb += `orm.NewDB(conf.DB` + strings.Title(v1.Name) + `)` + "\n" + v1.Name + `.SchemaMigrate()` + "\n"
		err := file.MkdirIfNotExist(dir)
		if err != nil {
			panic(err)
		}
		tables := v1.Models
		tableStr := ""

		for _, v2 := range tables {
			tableName := v2.Name
			tableStruct := file.UnderlineToCamel(v2.Name)
			tableStr += "_ =  GetDB(ctx).Set(\"gorm:table_options\", \"CHARSET=utf8mb4 comment='" + v2.Comment + "' AUTO_INCREMENT=1;\").AutoMigrate(&" + tableStruct + "{})\n"
			tabelPath := dir + "/" + tableName + ".go"
			fieldStr := ""
			fields := v2.Fields
			for _, v3 := range fields {
				fieldStr += createField(v3)
			}
			createTable(v1.Name, tableStruct, tableName, fieldStr, fileBuffer)
			fileWriter(fileBuffer, tabelPath)

		}

		createSchema(v1.Name, tableStr, fileBuffer, projectName, mysqlList)
		fileForceWriter(fileBuffer, dir+"/mysql_client.go")
		utils.PrintHint(dir + " Has been created.")

	}

}

// createField 创建gorm对应的字段
func createField(field string) string {
	tags := strings.Split(field, ";")
	if len(tags) == 0 {
		return ""
	}

	fieldMap := make(map[string]string)
	for _, v1 := range tags {
		attributes := strings.Split(v1, ":")
		if len(attributes) < 2 {
			continue
		}
		fieldMap[attributes[0]] = attributes[1]
	}
	fieldName := fieldMap["column"]
	upFieldName := file.UnderlineToCamel(fieldName)
	fieldType := def.GetTypeName(fieldMap["type"])
	return upFieldName + "  " + fieldType + " `json:\"" + fieldName + "\" gorm:\"" + field + "\"`\n"
}

func createSchema(dbName, tabels string, buffer *bytes.Buffer, projectName string, mysqlList []conf.Mysql) {
	buffer.WriteString(`
package `)
	buffer.WriteString(dbName)
	buffer.WriteString(`

import (
	"fmt"
	"context"

	`)
	// buffer.WriteString(config.ProjectName)
	buffer.WriteString(fmt.Sprintf(`
	"github.com/weblazy/easy/db/emysql"
	"gorm.io/gorm"
)
const %sMysql = "%sMysql"

func GetDB(ctx context.Context) *gorm.DB {
	return emysql.GetMysql(ctx, %sMysql)
}

func SchemaMigrate() {
	fmt.Println("开始初始化`, file.UnderlineToCamel(dbName), file.UnderlineToCamel(dbName), file.UnderlineToCamel(dbName)))
	buffer.WriteString(dbName)
	buffer.WriteString(`数据库")
	//自动建表，数据迁移
	ctx := context.Background()
    `)
	buffer.WriteString(tabels)
	buffer.WriteString(`
	fmt.Println("数据库`)
	buffer.WriteString(dbName)
	buffer.WriteString(`初始化完成")
}`)

}

func createTable(dbName, tableStruct, tableName, fields string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package `)
	buffer.WriteString(dbName)
	buffer.WriteString(`
import(
	"time"
	 "gorm.io/gorm"
	 "github.com/shopspring/decimal"
	 "github.com/weblazy/easy/db/emysql"
	 )
var `)
	buffer.WriteString(tableStruct)
	buffer.WriteString(`Handler = &`)
	buffer.WriteString(tableStruct)
	buffer.WriteString(`{}

type `)
	buffer.WriteString(tableStruct)
	buffer.WriteString(` struct {
	`)
	buffer.WriteString(fields)
	buffer.WriteString(`
}

func (t * `)
	buffer.WriteString(tableStruct)
	buffer.WriteString(`) TableName() string {
	return "`)
	buffer.WriteString(tableName)
	buffer.WriteString(`"
}

`)
	buffer.WriteString(`
	func (t * ` + tableStruct + `) Insert(ctx context.Context, data *` + tableStruct + `) error {
		db := GetDB(ctx)
	return db.Create(data).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) BulkInsert(ctx context.Context, fields []string, params []map[string]interface{}) error {
		db := GetDB(ctx)
	return emysql.BulkInsert(db, t.TableName(), fields, params)
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) BulkSave(ctx context.Context, fields []string, params []map[string]interface{}) error {
		db := GetDB(ctx)
	return emysql.BulkSave(db, t.TableName(), fields, params)
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) Delete(ctx context.Context, where string, args ...interface{}) error {
		db := GetDB(ctx)
	return db.Where(where, args...).Delete(&` + tableStruct + `{}).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) Updates(ctx context.Context, data map[string]interface{}, where string, args ...interface{}) (int64, error) {
		db := GetDB(ctx)
	db = db.Model(&` + tableStruct + `{}).Where(where, args...).Updates(data)
	return db.RowsAffected, db.Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetOne(ctx context.Context, where string, args ...interface{})(*` + tableStruct + `, error) {
	var obj ` + tableStruct + `
	return &obj, GetDB(ctx).Where(where, args...).Take(&obj).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetFirst(ctx context.Context, order string, where string, args ...interface{})(*` + tableStruct + `, error) {
	var obj ` + tableStruct + `
	return &obj, GetDB(ctx).Where(where, args...).Order(order).Take(&obj).Error
}`)

	buffer.WriteString(`
	func (* ` + tableStruct + `) GetList(ctx context.Context, where string, args ...interface{}) ([]*` + tableStruct + `, error) {
	var list []*` + tableStruct + `
	return list, GetDB(ctx).Where(where, args...).Find(&list).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetListWithLimit(ctx context.Context, limit int, where string, args ...interface{}) ([]*` + tableStruct + `, error) {
	var list []*` + tableStruct + `
	return list, GetDB(ctx).Where(where, args...).Limit(limit).Find(&list).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetListOrder(ctx context.Context, order string, where string, args ...interface{}) ([]*` + tableStruct + `, error) {
	var list []*` + tableStruct + `
	return list,GetDB(ctx).Where(where, args...).Order(order).Find(&list).Error}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetListOrderLimit(ctx context.Context, order string, limit int, where string, args ...interface{}) ([]*` + tableStruct + `, error) {
	var list []*` + tableStruct + `
	if limit == 0 || limit > 10000 {
		limit = 10
	}
	return list,GetDB(ctx).Where(where, args...).Order(order).Limit(limit).Find(&list).Error}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetListPage(ctx context.Context, pageNum, limit int, where string, args ...interface{}) ([]*` + tableStruct + `, error) {
	var list []*` + tableStruct + `
	offset := (pageNum - 1) * limit
	return list, GetDB(ctx).Where(where, args...).Offset(offset).Limit(limit).Find(&list).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetCount(ctx context.Context, where string, args ...interface{}) (int64, error) {
	var count int64
	return count, GetDB(ctx).Model(&` + tableStruct + `{}).Where(where, args...).Count(&count).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetSumInt64(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	type sum struct {
		Num int64 ` + "`" + `json:"num" gorm:"column:num"` + "`" + `
	}
	var obj sum
	return obj.Num, GetDB(ctx).Raw(sql, args...).Scan(&obj).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetSumFloat64(ctx context.Context, sql string, args ...interface{}) (float64, error) {
	type sum struct {
		Num float64 ` + "`" + `json:"num" gorm:"column:num"` + "`" + `
	}
	var obj sum
	return obj.Num, GetDB(ctx).Raw(sql, args...).Scan(&obj).Error
}`)

	buffer.WriteString(`
	func (t * ` + tableStruct + `) GetSumDecimal(ctx context.Context, sql string, args ...interface{}) (decimal.Decimal, error) {
	type sum struct {
		Num decimal.Decimal ` + "`" + `json:"num" gorm:"column:num"` + "`" + `
	}
	var obj sum
	return obj.Num, GetDB(ctx).Raw(sql, args...).Scan(&obj).Error
}`)
}

func fileForceWriter(buffer *bytes.Buffer, path string) {
	w.Add(buffer.Bytes())
	w.ForceWriteToFile(path)
	buffer.Reset()
}

func fileWriter(buffer *bytes.Buffer, path string) {
	w.Add(buffer.Bytes())
	w.WriteToFile(path)
	buffer.Reset()
}
