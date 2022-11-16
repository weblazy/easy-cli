package template

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sunmi-OS/gocore/v2/utils"
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
	handlersList := grpc.GrpcServers
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

	handlerStr := ""
	handlerRegister := ""
	pkgStr := ""
	for _, v1 := range handlersList {
		handlerName := v1.ModuleName

		handlerStr += fmt.Sprintf("%sServer := %s.New%sServer(svcCtx)\n", handlerName, handlerName, strings.Title(handlerName))
		handlerRegister += fmt.Sprintf("%s.Register%sServer(s, %sServer)", handlerName, strings.Title(handlerName), handlerName)
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

		fileBuffer.WriteString(fmt.Sprintf(`package %s
		type %sServer struct{
			svcCtx *svc.ServiceContext		
		}

		func New%sServer(svcCtx *svc.ServiceContext) *%sServer {
			return &%sServer{
				svcCtx: svcCtx,
			}
		}
		`, handlerName, strings.Title(handlerName), strings.Title(handlerName), strings.Title(handlerName), strings.Title(handlerName)))
		fileForceWriter(fileBuffer, servicelogicDir+file.CamelToUnderline(handlerName)+"_server.go")

		rpcFunc := ""
		params := ""
		for _, v2 := range routes {
			params += createRpcParam(v2)
			route := v2.Name
			function := strings.Title(route)
			rpcFunc += fmt.Sprintf("  rpc %s(%sRequest) returns (%sResponse);", function, function, function)

			logic := fmt.Sprintf(`
			package %s
			func (l *%sServer) %s(ctx context.Context, req *%s.%sRequest) (*%s.%sResponse, error) {
					resp := &%s.%sResponse{}
					return resp,nil
}`, handlerName, strings.Title(handlerName), function, handlerName, function, handlerName, function, handlerName, function)
			fileBuffer.WriteString(logic)
			fileForceWriter(fileBuffer, servicelogicDir+file.CamelToUnderline(route)+".go")
		}

		params += createRpcParams(v1.Params)

		CreateProto(fileBuffer, handlerName, params, rpcFunc)
		fileForceWriter(fileBuffer, serviceProtoDir+file.CamelToUnderline(handlerName)+".proto")
		createPb(serviceProtoDir + file.CamelToUnderline(handlerName) + ".proto")
		pkgStr += "\"" + name + "/grpcs/" + grpc.Name + "/" + file.CamelToUnderline(handlerName) + "\"\n"
	}
	GreateCmd(grpc.Name, homedir, pkgStr, handlerStr, handlerRegister)
}

func CreateProto(buffer *bytes.Buffer, service, param, rpcFunc string) {
	buffer.WriteString(fmt.Sprintf(`
syntax = "proto3";

package %s;

option go_package = "./%s";

%s

service %sService{
	%s
}
`, service, service, param, strings.Title(service), rpcFunc))

}

func createRpcParam(handle conf.Handle) string {
	params := ""
	fields := handle.RequestParams
	for k1, v3 := range fields {
		params += fmt.Sprintf("  %s %s = %d;\n", v3.Type, v3.Name, k1+1)
	}

	req := fmt.Sprintf(`message %sRequest{
  %s
}
`, handle.Name, params)

	params = ""
	fields = handle.ResponseParams
	for k1, v3 := range fields {
		params += fmt.Sprintf("  %s %s = %d;\n", v3.Type, v3.Name, k1+1)
	}
	resp := fmt.Sprintf(`message %sResponse{
  %s
}
`, strings.Title(handle.Name), params)

	return req + resp
}

func createRpcParams(paramsMap map[string][]conf.Param) string {
	paramStr := ""
	for k1, v1 := range paramsMap {
		paramStrContent := ""
		for k2, v2 := range v1 {
			paramStrContent += fmt.Sprintf("  %s %s = %d;\n", v2.Type, v2.Name, k2+1)
		}

		paramStr += fmt.Sprintf(`message %s{
  %s
}`, strings.Title(k1), paramStrContent)
	}

	return paramStr
}

func GreateCmd(grpcName, homedir, pkgStr string, handlerStr, handlerRegister string) {
	cmd := fmt.Sprintf(`
	package %s
	import (
		"log"
		"net"
		"order/cmd"

		"github.com/urfave/cli/v2"

		"github.com/weblazy/easy/utils/closes"
		"github.com/weblazy/gocore/viper"
		"google.golang.org/grpc"
		%s
	)


	var Cmd = &cli.Command{
	Name:    "%s",
	Aliases: []string{},
	Usage:   "%s start",
	Subcommands: []*cli.Command{
		{
			Name:   "start",
			Usage:  "start service",
			Action: Run,
		},
	},
}

func Run(c *cli.Context) error {
	defer closes.Close()

	cmd.InitConf()
	cmd.InitDB()
	cmd.InitCache()
	svcCtx := svc.NewServiceContext(c)

	%s

	listen, err := net.Listen("tcp", viper.C.GetString("network.GrpcServicePort"))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
    %s
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}

	return nil
}
`, grpcName, pkgStr, grpcName, grpcName, handlerStr, handlerRegister)
	fileBuffer.WriteString(cmd)
	fileForceWriter(fileBuffer, homedir+"/"+file.CamelToUnderline(grpcName)+".go")
}

func createPb(path string) {
	resp, err := utils.Cmd("protoc", []string{"--go_out=.", "--go_opt=paths=source_relative", "--go-grpc_out=.", "--go-grpc_opt=paths=source_relative", path})
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}
}
