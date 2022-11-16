// Code generated by hero.
// source: /Users/liuguoqiang/Desktop/go/mod/gocore/tools/gocore/template/cmd_job.got
// DO NOT EDIT!
package template

import "bytes"

func FromCmdJob(name, jobCmd, jobFunctions string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package jobs

import (
	"`)
	buffer.WriteString(name)
	buffer.WriteString(`/app/jobs"
	"github.com/urfave/cli/v2"
	"github.com/weblazy/easy/utils/closes"
)

// Job cmd 任务相关
var Cmd = &cli.Command{
	Name:    "job",
	Aliases: []string{"j"},
	Usage:   "job",
	Subcommands: []*cli.Command{
		`)
	buffer.WriteString(jobCmd)
	buffer.WriteString(`
	},
}
`)
	buffer.WriteString(jobFunctions)

}