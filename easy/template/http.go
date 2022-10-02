package template

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/weblazy/easy-cli/easy/conf"
	"github.com/weblazy/easy-cli/easy/file"
)

func createHttps(root, name string) {
	for _, v := range goCoreConfig.HttpApis {
		homedir := root + "/https/" + v.Name
		createApi(root, name, homedir, v)
		createDef(root, homedir, v)
		createErrCode(root, homedir, v)
	}
}

func createApi(root, name, homedir string, httpApi conf.HttpApi) {

	handlersList := httpApi.Apis
	if len(handlersList) == 0 {
		return
	}

	err := file.MkdirIfNotExist(homedir)
	if err != nil {
		panic(err)
	}

	FromCmdApi(name, httpApi.Name, fileBuffer)
	fileForceWriter(fileBuffer, homedir+"/"+httpApi.Name+".go")

	handlerDir := homedir + "/handler/"
	err = file.MkdirIfNotExist(handlerDir)
	if err != nil {
		panic(err)
	}

	domainDir := homedir + "/logic/"
	err = file.MkdirIfNotExist(domainDir)
	if err != nil {
		panic(err)
	}

	routesStr := ""

	//handlers := make([]string, 0)

	for _, v1 := range handlersList {
		handlerName := v1.ModuleName
		routesStr += "\n" + handlerName + ":=router.Group(\"" + v1.Prefix + "\")\n"
		handlerPath := handlerDir + file.CamelToUnderline(handlerName) + ".go"
		routes := v1.Handle
		FromDomain(fileBuffer)
		fileForceWriter(fileBuffer, domainDir+file.CamelToUnderline(handlerName)+".go")
		if len(routes) == 0 {
			continue
		}
		// 首字母大写
		handler := strings.Title(handlerName)
		//handlers = append(handlers, handler)
		functions := make([]string, 0)
		comments := make([]string, 0)
		reqs := make([]string, 0)

		handlerContent := ""
		handlerFile, err := os.Open(handlerPath)
		if err == nil {
			content, err := ioutil.ReadAll(handlerFile)
			if err != nil {
				panic(err)
			}
			handlerContent = string(content)
		}
		for _, v2 := range routes {
			route := v2.Name
			function := strings.Title(route)
			if strings.Contains(handlerContent, "func "+function+"(g *gin.Context)") {
				continue
			}
			functions = append(functions, function)
			req := strings.Title(v2.Name)
			reqs = append(reqs, req)
			comments = append(comments, v2.Comment)
			routesStr += handlerName + "." + v2.Method + "(\"/" + file.CamelToUnderline(handlerName) + "/" + file.CamelToUnderline(route) + "\",handler." + function + ") //" + v2.Comment + "\n"
			// FromDomain(name, handler, function, req, fileBuffer)
			// fileForceWriter(fileBuffer, domainDir+file.CamelToUnderline(route)+".go")

		}

		FromHandler(name, httpApi.Name, handler, handlerContent, comments, functions, reqs, fileBuffer)
		// writer.Add(fileBuffer.Bytes())
		fileForceWriter(fileBuffer, handlerPath)
	}
	FromRoutes(name, httpApi.Name, routesStr, fileBuffer)
	err = file.MkdirIfNotExist(homedir + "/routes")
	if err != nil {
		panic(err)
	}
	fileForceWriter(fileBuffer, homedir+"/routes/routers.go")

	createPostman(root, name, httpApi.CommonHeaders, httpApi.Apis)
}
