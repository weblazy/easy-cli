// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/conf_base.got
// DO NOT EDIT!
package template

import "bytes"

func FromConfBase(baseConf string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package conf

var BaseConfig = ` + "`" + `
[network]
ApiServiceHost = "`)
	buffer.WriteString("127.0.0.1")
	buffer.WriteString(`"
ApiServicePort = "`)
	buffer.WriteString("80")
	buffer.WriteString(`"

`)
	buffer.WriteString(baseConf)
	buffer.WriteString(`
` + "`" + ``)

}
