package main

import (
	. "ACG/parser"
	"fmt"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	path := "./nika_vm.gram"
	var parser = NewBNFParser()
	rules, _ := parser.ParseFile(path)

	//generate parser
	var generator Generator
	generator.SetOutputPath("./parser/nika_parser.go")
	generator.GenerateLParser("NikaParser", rules)

	//gen compiler
	generator.SetOutputPath("./parser/nika_compiler.go")
	generator.PrintCompiler("NikaCompiler", rules)

}

func TestHelloEmpty(t *testing.T) {
	var file FileReader
	content := file.Read("./testcode.c")
	var prog string = string(content)

	var nika = NewNikaParser()
	err := nika.Tokenize(prog)
	if err != nil && t != nil {
		t.Fatalf("nika.Tokenize")
	}
	ret := nika.PROG()
	if ret.Err != nil && t != nil {
		t.Fatal("niake.PROG()")
	}
	count := NodeCount(ret.Nd)
	fmt.Printf("total tree nodes count:%d\n", count)

	var travel Travel
	var printer NodePrinter
	printer.OpenFile("./test_tree.txt")
	travel.DepthFirstTravel(ret.Nd, &printer)
	printer.CloseFile()

	compiler := NewNikaCompiler()
	compiler.Compile(ret.Nd)
	vm := NewVM(compiler)
	vm.Print()
	vm.Run()
	i := 0

	//last
	if vm.IsEmpty() == false {
		obj := vm.Back()
		if obj != nil {
			fmt.Printf("result:%s\n", obj.ToString())
		}
	}
	for vm.IsEmpty() == false {
		obj := vm.Pop()
		if obj != nil {
			fmt.Printf("result:%s\n", obj.ToString())
		}
		// if obj.(*NkInteger).Value != answers[i] {
		// 	t.Fatalf("%d not equal %d", obj.(*NkInteger).Value, answers[i])
		// }
		i++
	}

}
