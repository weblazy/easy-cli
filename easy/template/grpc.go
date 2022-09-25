package template

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/weblazy/easy-cli/easy/conf"
	"github.com/weblazy/easy-cli/easy/file"
)

func createGrpcs(root, name string) {
	for _, v := range goCoreConfig.Grpcs {
		homedir := root + "/grpcs/" + v.Name
		createGrpcProtoHandler(root, name, homedir, v)
		// createGrpcHandler(root, name, homedir, v)
		// createDef(root, homedir, v)
		// createErrCode(root, homedir, v)
	}
}

func createGrpcProtoHandler(root, name, homedir string, grpc conf.Grpc) {
	handlersList := grpc.Apis
	if len(handlersList) == 0 {
		return
	}

	apiDir := homedir + "/handler/"
	err := file.MkdirIfNotExist(apiDir)
	if err != nil {
		panic(err)
	}

	logicDir := homedir + "/logic/"
	err = file.MkdirIfNotExist(logicDir)
	if err != nil {
		panic(err)
	}

	protoDir := homedir + "/proto/"
	err = file.MkdirIfNotExist(protoDir)
	if err != nil {
		panic(err)
	}

	for _, v1 := range handlersList {
		handlerName := v1.ModuleName

		serviceProtoDir := protoDir + file.CamelToUnderline(handlerName) + "/"
		err = file.MkdirIfNotExist(serviceProtoDir)
		if err != nil {
			panic(err)
		}
		// fileForceWriter(fileBuffer, serviceProtoDir+file.CamelToUnderline(handlerName)+".go")

		servicelogicDir := logicDir + file.CamelToUnderline(handlerName) + "/"
		err = file.MkdirIfNotExist(servicelogicDir)
		if err != nil {
			panic(err)
		}
		// fileForceWriter(fileBuffer, servicelogicDir+file.CamelToUnderline(handlerName)+".go")

		// apiPath := apiDir + file.CamelToUnderline(handlerName) + ".go"
		routes := v1.Handle
		// FromDomain(fileBuffer)
		// fileForceWriter(fileBuffer, domainDir+file.CamelToUnderline(handlerName)+".go")
		// if len(routes) == 0 {
		// 	continue
		// }

		rpcFunc := ""
		params := ""
		for _, v2 := range routes {
			params += createRpcParam(v2)
			route := v2.Name
			function := strings.Title(route)
			rpcFunc += fmt.Sprintf("  rpc %s(%sRequest) returns (%sResponse);", function, function, function)

			logic := fmt.Sprintf(`
			func (l *%sLogic) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {
					resp := &%s.%sResponse{}
					return resp,nil
}`, function, function, handlerName, function, handlerName, function, handlerName, function)
			fileBuffer.WriteString(logic)
			fileForceWriter(fileBuffer, servicelogicDir+file.CamelToUnderline(route)+".go")
		}

		CreateProto(fileBuffer, handlerName, params, rpcFunc)
		fileForceWriter(fileBuffer, serviceProtoDir+file.CamelToUnderline(handlerName)+".proto")
	}

}

func CreateProto(buffer *bytes.Buffer, service, param, rpcFunc string) {
	buffer.WriteString(fmt.Sprintf(`
syntax = "proto3";

package %s;

%s

service %s{
	%s
}
`, service, param, service, rpcFunc))

}

func createRpcParam(handle conf.Handle) string {
	params := ""
	fields := handle.RequestParams
	for k1, v3 := range fields {
		params += fmt.Sprintf("  %s %s = %d;\n", v3.Type, v3.Name, k1)
	}

	req := fmt.Sprintf(`message %sRequest{
  %s
}`, strings.Title(handle.Name), params)

	params = ""
	fields = handle.ResponseParams
	for k1, v3 := range fields {
		params += fmt.Sprintf("  %s %s = %d;\n", v3.Type, v3.Name, k1)
	}
	resp := fmt.Sprintf(`message %sResponse{
  %s
}`, strings.Title(handle.Name), params)

	return req + resp
}
