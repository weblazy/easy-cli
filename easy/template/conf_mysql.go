// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/conf_mysql.got
// DO NOT EDIT!
package template

import (
	"bytes"
	"strings"
)

func FromConfMysql(dbName string, buffer *bytes.Buffer) {
	buffer.WriteString(`
[db`)
	buffer.WriteString(strings.Title(dbName))
	buffer.WriteString(`]
Host = ""           #数据库连接地址
Name = "`)
	buffer.WriteString(dbName)
	buffer.WriteString(`"           #数据库名称
User = ""           #数据库用户名
Passwd = ""         #数据库密码
Port = "3306"       #数据库端口号
`)

}
