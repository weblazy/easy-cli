package utils

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

// PrintBanner 打印项目名称
func PrintBanner(name string) {
	myFigure := figure.NewFigure(name, "", true)
	myFigure.Print()
	fmt.Println()
}

// PrintHint 打印提示
func PrintHint(str string) {
	_, err := color.New(color.FgCyan, color.Bold).Print(str + "\n")
	if err != nil {
		return
	}
}
