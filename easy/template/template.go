package template

import (
	"bytes"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/weblazy/easy-cli/easy/conf"
	"github.com/weblazy/easy-cli/easy/def"
	"github.com/weblazy/easy-cli/easy/file"
)

var writer = file.NewWriter()

var fileBuffer = new(bytes.Buffer)

var localConf = `
[base]
debug = true
`

var goCoreConfig *conf.GoCore

// CreateCode 更具配置文件生成项目
// 模板引擎生成语句 hero -source=./tools/gocore/template -extensions=.got,.md,.docker
func CreateCode(root, name string, config *conf.GoCore) {

	goCoreConfig = config
	newProgress(11, "start preparing...")
	time.Sleep(time.Second)
	progressNext("Initialize the directory structure...")
	mkdir(root)
	progressNext("Initialize the configuration file...")
	createConf(root, name)
	progressNext("Initialize the main program...")
	createMain(root, name)
	progressNext("Initialize the Dockerfile...")
	createDockerfile(root)
	progressNext("Initialize the Readme...")
	createReadme(root)
	progressNext("Initialize the DB Model...")
	createModel(root, name)
	progressNext("Initialize the Cronjob...")
	createCronjob(name, root)
	progressNext("Initialize the Job...")
	createJob(name, root)
	progressNext("Initialize the Http...")
	createHttps(root, name)
	progressNext("Initialize the Grpc...")
	createGrpcs(root, name)
	progressNext("Initialize the Request return parameters...")

}

// CreateField 创建gorm对应的字段
func CreateField(field string) string {
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

func createMain(root, name string) {
	var cmdList []string
	for _, v := range goCoreConfig.HttpApis {
		cmdList = append(cmdList, v.Name+".Cmd,")
	}
	for _, v := range goCoreConfig.Grpcs {
		cmdList = append(cmdList, v.Name+".Cmd,")
	}
	if len(goCoreConfig.CronJobs) > 0 {
		cmdList = append(cmdList, "cronjobs.Cmd,")
	}
	if len(goCoreConfig.Jobs) > 0 {
		cmdList = append(cmdList, "jobs.Cmd,")
	}
	FromMain(name, cmdList, fileBuffer)
	fileForceWriter(fileBuffer, root+"/main.go")
}

func createConf(root string, name string) {
	FromConfConst(name, fileBuffer)
	fileForceWriter(fileBuffer, root+"/conf/const.go")
}

func createDockerfile(root string) {
	FromDockerfile(fileBuffer)
	fileForceWriter(fileBuffer, root+"/Dockerfile")
}

func createReadme(root string) {
	FromREADME(fileBuffer)
	fileForceWriter(fileBuffer, root+"/README.md")
}

func createErrCode(root, homedir string, httpApi conf.HttpApi) {
	FromErrCode(fileBuffer)
	err := file.MkdirIfNotExist(homedir + "/errcode")
	if err != nil {
		panic(err)
	}
	fileWriter(fileBuffer, homedir+"/errcode/errcode.go")
}

func createModel(root, name string) {
	mysqlMap := goCoreConfig.Config.CMysql
	pkgs := ""
	dbUpdate := ""
	dbUpdateRedis := ""
	baseConf := ""
	if len(goCoreConfig.Config.CRedis) > 0 {
		dbUpdateRedis = "var err error"
	}
	if len(mysqlMap) > 0 {
		dbUpdate = "var err error"
	}
	InitDB := ""
	initRedis := ""
	for _, v1 := range mysqlMap {
		pkgs += `"` + name + `/model/` + v1.Name + `"` + "\n"
		dir := root + "/model/" + v1.Name
		dbUpdate += `
				err = orm.NewOrUpdateDB(conf.DB` + strings.Title(v1.Name) + `)
				if err != nil {
					elog.ErrorCtx(context.Background(), "InitMysqlErr", elog.FieldError(err))
				}
		`
		InitDB += `orm.NewDB(conf.DB` + strings.Title(v1.Name) + `)` + "\n" + v1.Name + `.SchemaMigrate()` + "\n"
		err := file.MkdirIfNotExist(dir)
		if err != nil {
			panic(err)
		}
		tables := v1.Models
		tableStr := ""

		for _, v2 := range tables {
			tableName := v2.Name
			tableStruct := file.UnderlineToCamel(v2.Name)
			tableStr += "_ = Orm().Set(\"gorm:table_options\", \"CHARSET=utf8mb4 comment='" + v2.Comment + "' AUTO_INCREMENT=1;\").AutoMigrate(&" + tableStruct + "{})\n"
			tabelPath := dir + "/" + tableName + ".go"
			fieldStr := ""
			fields := v2.Fields
			for _, v3 := range fields {
				fieldStr += CreateField(v3)
			}
			FromModelTable(v1.Name, tableStruct, tableName, fieldStr, fileBuffer)
			fileWriter(fileBuffer, tabelPath)

		}

		FromModel(v1.Name, tableStr, fileBuffer)
		fileForceWriter(fileBuffer, dir+"/mysql_client.go")

		buff := new(bytes.Buffer)
		FromConfMysql(v1.Name, buff)
		localConf += buff.String()

		for _, v1 := range goCoreConfig.Config.CRedis {
			for k2 := range v1.Index {
				localConf += `
[` + v1.Name + `]
host = "" 
port = ":6379"
auth = ""
prefix = ""
`

				baseConf += `[` + v1.Name + `.redisDB]
` + k2 + ` = ` + cast.ToString(v1.Index[k2])
				initRedis += "redis.NewRedis(conf." + strings.Title(v1.Name) + strings.Title(k2) + "Redis)\n"
				dbUpdateRedis += `		
				err = redis.NewOrUpdateRedis(conf.` + strings.Title(v1.Name) + strings.Title(k2) + `Redis)
				if err != nil {
					elog.ErrorCtx(context.Background(), "InitRedisErr", elog.FieldError(err))
				}
		`
			}
		}
		if goCoreConfig.Config.CRocketMQConfig {
			localConf += `
			
[aliyunmq]
NameServer = ""
AccessKey = ""
SecretKey = ""
Namespace = ""

			`
		}

	}
	if !goCoreConfig.Config.CNacos {
		FromConfLocal("DevConfig", localConf, fileBuffer)
		fileWriter(fileBuffer, root+"/conf/dev.go")
		FromConfLocal("TestConfig", localConf, fileBuffer)
		fileWriter(fileBuffer, root+"/conf/test.go")
		FromConfLocal("UatConfig", localConf, fileBuffer)
		fileWriter(fileBuffer, root+"/conf/uat.go")
		FromConfLocal("OnlConfig", localConf, fileBuffer)
		fileWriter(fileBuffer, root+"/conf/onl.go")
	}
	FromConfLocal("LocalConfig", localConf, fileBuffer)
	fileWriter(fileBuffer, root+"/conf/local.go")
	FromCmdInit(name, pkgs, dbUpdate, InitDB, initRedis, dbUpdateRedis, fileBuffer)
	fileForceWriter(fileBuffer, root+"/cmd/init.go")

	FromConfBase(baseConf, fileBuffer)
	fileForceWriter(fileBuffer, root+"/conf/base.go")
}

func createCronjob(name, root string) {
	jobs := goCoreConfig.CronJobs
	if len(jobs) == 0 {
		return
	}

	dir := root + "/cronjobs/"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}

	handlerDir := dir + "/handler/"
	err = file.MkdirIfNotExist(handlerDir)
	if err != nil {
		panic(err)
	}

	cronjobs := ""
	for _, v1 := range jobs {

		jobPath := handlerDir + file.CamelToUnderline(v1.Job.Name) + ".go"
		FromCronJob(v1.Job.Name, v1.Job.Comment, fileBuffer)
		fileForceWriter(fileBuffer, jobPath)
		cronjobs += "_,_ = cronJob.AddFunc(\"" + v1.Spec + "\", handler." + v1.Job.Name + ")\n"
	}

	FromCmdCronJob(name, cronjobs, fileBuffer)
	fileForceWriter(fileBuffer, dir+"cronjobs.go")
}

func createJob(name, root string) {

	jobs := goCoreConfig.Jobs
	if len(jobs) == 0 {
		return
	}

	dir := root + "/jobs/"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}

	handlerDir := dir + "/handler/"
	err = file.MkdirIfNotExist(handlerDir)
	if err != nil {
		panic(err)
	}

	jobCmd := ""
	jobFunctions := ""
	for _, v1 := range jobs {
		FromJob(v1.Name, v1.Comment, fileBuffer)
		fileForceWriter(fileBuffer, handlerDir+file.CamelToUnderline(v1.Name)+".go")
		jobCmd += `		{
			Name:   "` + v1.Name + `",
			Usage:  "` + v1.Comment + `",
			Action: ` + v1.Name + `,
		},`
		jobFunctions += `
func ` + v1.Name + `(c *cli.Context) error {
	defer closes.Close()
	// 初始化必要内容
	cmd.InitConf()
	cmd.InitDB()
	handler.` + v1.Name + `()
	return nil
}
`
	}

	FromCmdJob(name, jobCmd, jobFunctions, fileBuffer)
	fileForceWriter(fileBuffer, root+"/jobs/jobs.go")
}

// ------------------------------------------------------------------------------

func fileForceWriter(buffer *bytes.Buffer, path string) {
	writer.Add(buffer.Bytes())
	writer.ForceWriteToFile(path)
	buffer.Reset()
}

func fileWriter(buffer *bytes.Buffer, path string) {
	writer.Add(buffer.Bytes())
	writer.WriteToFile(path)
	buffer.Reset()
}

func unResetfileWriter(buffer *bytes.Buffer, path string) {
	writer.Add(buffer.Bytes())
	writer.WriteToFile(path)
}

func mkdir(root string) {
	var dirList = []string{
		"/common",
		"/cmd",
		// "/app/domain",
		"/model",
		// "/app/errcode",
		// "/app/routes",
		"/conf",
		"/pkg",
	}
	for _, dir := range dirList {
		dir = root + dir
		err := file.MkdirIfNotExist(dir)
		if err != nil {
			panic(err)
		}
	}
}
