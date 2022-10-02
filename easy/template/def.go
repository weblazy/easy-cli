package template

import (
	"strings"

	"github.com/weblazy/easy-cli/easy/conf"
	"github.com/weblazy/easy-cli/easy/file"
)

func createDef(root, homedir string, httpApi conf.HttpApi) {
	modules := httpApi.Apis
	if len(modules) == 0 {
		return
	}
	dir := homedir + "/def"
	err := file.MkdirIfNotExist(dir)
	if err != nil {
		panic(err)
	}

	writer.Add([]byte(`package def` + "\n"))
	for _, v1 := range modules {
		for _, v2 := range v1.Handle {

			params := ""
			fields := v2.RequestParams
			for _, v3 := range fields {
				// field := strings.Split(v2.String(), ";")
				// if len(field) < 3 {
				// 	continue
				// }
				params += file.UnderlineToCamel(v3.Name) + " " + v3.Type + " `json:\"" + v3.Name + "\" binding:\"" + v3.Validate + "\"` // " + v3.Comment + "\n"
			}
			FromApiRequest(strings.Title(v2.Name)+"Request", params, fileBuffer)

			params = ""
			fields = v2.ResponseParams
			for _, v3 := range fields {
				params += file.UnderlineToCamel(v3.Name) + " " + v3.Type + " `json:\"" + v3.Name + "\" binding:\"" + v3.Validate + "\"` // " + v3.Comment + "\n"
			}
			FromApiRequest(strings.Title(v2.Name)+"Response", params, fileBuffer)
		}
	}
	for k1, v1 := range httpApi.Params {
		params := ""
		fields := v1
		for _, v2 := range fields {
			params += file.UnderlineToCamel(v2.Name) + " " + v2.Type + " `json:\"" + v2.Name + "\" binding:\"" + v2.Validate + "\"` // " + v2.Comment + "\n"
		}
		FromApiRequest(k1, params, fileBuffer)
	}
	fileForceWriter(fileBuffer, dir+"/def.go")
}
