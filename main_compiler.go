package main

import (
	. "nika/parser"
	"fmt"
	"os"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("请将 .nika 文件拖到此程序上运行")
		return
	}
	srcPath := os.Args[1]


	var file FileReader
	content := file.Read(srcPath)
	var prog string = string(content)

	var nika = NewNikaParser()
	err := nika.Tokenize(prog)
	if err != nil {
		fmt.Println("nika.Tokenize fail!")
		fmt.Scanln()
		return
	}
	ret := nika.PROG()
	if ret.Err != nil {
		fmt.Println("niake.PROG() fail!")
		fmt.Scanln()
		return
	}

	var travel Travel
	var printer NodePrinter
	printer.OpenFile("./test_tree.txt")
	travel.DepthFirstTravel(ret.Nd, &printer)
	printer.CloseFile()

	compiler := NewNikaCompiler()
	compiler.C(ret.Nd)
	vm := NewStackMachine(compiler)
	vm.Print()
	vm.Run()
}