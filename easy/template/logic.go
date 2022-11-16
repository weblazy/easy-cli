// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/api.got
// DO NOT EDIT!
package template

import (
	"bytes"
	"fmt"
)

func FromLogic(homePath,homedir,handlerName string , comments []string, functions []string,  buffer *bytes.Buffer) {

	for k1, v1 := range functions {
		buffer.WriteString(fmt.Sprintf(`
package %s

import (
	"%s/def"
	"github.com/weblazy/easy/utils/http/http_server/service"
)
    // %s
	func %s(svcCtx *service.ServiceContext, req *def.%sRequest) *service.Response {

		return nil
	}
`,handlerName,homePath, comments[k1],v1,v1))
	}

}