// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/domain_handler.got
// DO NOT EDIT!
package template

import "bytes"

func FromDomainHandler(handlers []string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package domain
`)
	for _, v1 := range handlers {
		buffer.WriteString(`
    var `)
		buffer.WriteString(v1)
		buffer.WriteString(`Handler = &`)
		buffer.WriteString(v1)
		buffer.WriteString(`{}
    type `)
		buffer.WriteString(v1)
		buffer.WriteString(` struct{}
`)
	}

}
