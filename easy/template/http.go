package template

import (
	"fmt"
	"strings"

	"github.com/weblazy/easy-cli/easy/conf"
	"github.com/weblazy/easy-cli/easy/file"
)

func createHttps(root, name string) {
	for _, v := range easyConfig.HttpApis {
		homedir := root + "/https/" + v.Name
		homePath := name + "/https/" + v.Name
		createApi(root, name, homedir, homePath, v)
		createDef(root, homedir, v)
		createErrCode(root, homedir, v)
	}
}

func createApi(root, name, homedir, homePath string, httpApi conf.HttpApi) {

	handlersList := httpApi.Apis
	if len(handlersList) == 0 {
		return
	}

	err := file.MkdirIfNotExist(homedir)
	if err != nil {
		panic(err)
	}

	FromCmdApi(httpApi.Name, homePath, fileBuffer)
	fileForceWriter(fileBuffer, homedir+"/"+httpApi.Name+".go")

	handlerDir := homedir + "/handler/"
	err = file.MkdirIfNotExist(handlerDir)
	if err != nil {
		panic(err)
	}
	configDir := homedir + "/config/"
	err = file.MkdirIfNotExist(configDir)
	if err != nil {
		panic(err)
	}
	configStr := "	HttpServerConfig *http_server_config.Config"
	configVar := "	HttpServerConfig: http_server_config.DefaultConfig(),"
	FromConfigInit(name, "", configStr, configVar, "", "", fileBuffer)
	fileForceWriter(fileBuffer, configDir+"config.go")

	logicDir := homedir + "/logic/"
	err = file.MkdirIfNotExist(logicDir)
	if err != nil {
		panic(err)
	}

	routesStr := ""
	routesPkg := ""
	//handlers := make([]string, 0)
	err = file.MkdirIfNotExist(homedir + "/routes")
	if err != nil {
		panic(err)
	}
	for _, v1 := range handlersList {
		handlerName := v1.ModuleName
		routesStr += "\n" + handlerName + "Group:=router.Group(\"" + v1.Prefix + "\")\n"
		routesStr += fmt.Sprintf("\n%sInterceptor(%sGroup)\n", handlerName, handlerName)
		routesPkg += fmt.Sprintf("\"%s/handler/%s\"\n", homePath, handlerName)
		// handlerPath := handlerDir + file.CamelToUnderline(handlerName) + ".go"
		routes := v1.Handle
		// FromDomain(fileBuffer)
		// fileForceWriter(fileBuffer, logicDir+file.CamelToUnderline(handlerName)+".go")
		// if len(routes) == 0 {
		// 	continue
		// }
		// 首字母大写
		// handler := strings.Title(handlerName)
		// handlers = append(handlers, handler)
		functions := make([]string, 0)
		comments := make([]string, 0)

		// handlerContent := ""
		// handlerFile, err := os.Open(handlerPath)
		// if err == nil {
		// 	content, err := ioutil.ReadAll(handlerFile)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	handlerContent = string(content)
		// }
		// err = file.MkdirIfNotExist(homedir + "/handler/" + handlerName)
		// if err != nil {
		// 	panic(err)
		// }
		err = file.MkdirIfNotExist(homedir + "/logic/" + handlerName)
		if err != nil {
			panic(err)
		}
		for _, v2 := range routes {
			route := v2.Name
			function := strings.Title(route)
			// if strings.Contains(handlerContent, "func "+function+"(g *gin.Context)") {
			// 	continue
			// }
			functions = append(functions, function)
			comments = append(comments, v2.Comment)
			routesStr += handlerName + "Group." + v2.Method + "(\"/" + file.CamelToUnderline(handlerName) + "/" + file.CamelToUnderline(route) + "\"," + "handler." + function + ") //" + v2.Comment + "\n"
			// FromDomain(name, handler, function, req, fileBuffer)
			// fileForceWriter(fileBuffer, domainDir+file.CamelToUnderline(route)+".go")

			FromHttpLogic(homePath, handlerName, function, v2.Comment, fileBuffer)
			fileForceWriter(fileBuffer, logicDir+file.CamelToUnderline(handlerName)+"/"+file.CamelToUnderline(function)+".go")
		}
		FromHttpHandler(httpApi.Name, homePath, homedir, handlerName, comments, functions, fileBuffer)
		fileForceWriter(fileBuffer, handlerDir+file.CamelToUnderline(handlerName)+".go")

		fileBuffer.WriteString(CreateInterceptor(homePath, handlerName))
		fileForceWriter(fileBuffer, homedir+"/routes/"+handlerName+"_interceptor.go")
	}
	routesPkg += fmt.Sprintf("\"%s/config\"\n", homePath)
	FromRoutes(routesPkg, routesStr, fileBuffer)

	fileForceWriter(fileBuffer, homedir+"/routes/routes.go")

	createPostman(root, name, httpApi.CommonHeaders, httpApi.Apis)
}
