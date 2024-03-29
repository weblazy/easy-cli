// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/api.got
// DO NOT EDIT!
package template

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/weblazy/easy-cli/easy/conf"
)

func FromHttpHandler(httpName,homePath, homedir, handlerName string, comments []string, functions []string, buffer *bytes.Buffer) {
	pkgStr := ""
	funcStr := ""
	for k1, v1 := range functions {
		pkgStr += fmt.Sprintf(`"%s/logic/%s"`, homePath, handlerName) + "\n"
		funcStr += fmt.Sprintf(`    
	// %s
	func %s(g *gin.Context) {
		 svcCtx := service.NewServiceContext(g)
		 req := new(def.%sRequest)
		 err := svcCtx.BindValidator(req)
		 if err != nil {
			svcCtx.Error(err)
			return
		 }
		svcCtx.Return(%s.%s(svcCtx, req))
	}
	`, comments[k1], v1, v1, handlerName, v1)
	}
	buffer.WriteString(fmt.Sprintf(`
package handler

import (
	"%s/def"
	%s

	"github.com/gin-gonic/gin"
	"github.com/weblazy/easy/http/http_server/service"
)

%s
`, homePath, pkgStr, funcStr))
}

func FromRpcHandler(homePath, handlerName string, functions []conf.Handle, buffer *bytes.Buffer) {
	pkgStr := ""
	funcStr := ""
	service := fmt.Sprintf(`
		type %sService struct{
			%s.Unimplemented%sServiceServer
		}

		func New%sService() *%sService {
			return &%sService{
			}
		}
		`,  strings.Title(handlerName),handlerName,strings.Title(handlerName), strings.Title(handlerName), strings.Title(handlerName), strings.Title(handlerName))
	for k1 := range functions {
		v1 := functions[k1]
		funcName := strings.Title(v1.Name)
		pkgStr += fmt.Sprintf(`
		"%s/logic/%s_logic"
		"%s/proto/%s"
		`, homePath, handlerName, homePath, handlerName)
		funcStr += fmt.Sprintf(`  
	// %s  
	func (h *%sService) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {
		svcCtx := &%s_logic.%sCtx{
			Log: code_err.NewLog(ctx),
			Req:      req,
			Res: new(%s.%sResponse),
		}
		err := %s_logic.%s(svcCtx)
		if err != nil {
			svcCtx.Res.Code = err.Code
			svcCtx.Res.Msg = err.Msg
		}
		return svcCtx.Res, nil
	}
	`, v1.Comment, strings.Title(handlerName), funcName, handlerName, funcName, handlerName, funcName, handlerName, funcName, handlerName, funcName, handlerName, funcName)
	}
	buffer.WriteString(fmt.Sprintf(`
package handler

import (
	"context"
     %s
	"github.com/weblazy/easy/code_err"
)
%s
%s
`, pkgStr, service, funcStr))
}
