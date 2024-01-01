package main

import (
	. "ACG/parser"
	"fmt"
	"testing"
)

type inter interface {
	a()
}

type inter2 interface {
	a()
	b()
}

type interA struct {
}

func (s *interA) a() {

}

type interB struct {
	interA
}

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

	generator.SetOutputPath("./parser/nika_nodes.go")
	generator.PrintNodes("NikaCompiler", rules)

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

	var travel Travel
	var printer NodePrinter
	printer.OpenFile("./test_tree.txt")
	travel.DepthFirstTravel(ret.Nd, &printer)
	printer.CloseFile()

	compiler := NewNikaCompiler()
	compiler.C(ret.Nd)
	vm := NewVM(compiler)
	vm.Print()
	vm.Run()

	answers := []int{0, 428, 0, -64, 138, 35860, 3, -17, 36348, 33, 14, 32, 0, 1, 1, 0, 3, 1, 1, 0, 0, 0, 1, 3, 6, 10, 15, 21, 28, 36, 45, 45}
	for i := 0; i < len(vm.DebugStack); i++ {
		obj := vm.DebugStack[i]
		if obj.(*NkInteger).Value != answers[i] {
			t.Fatalf("not equal")
		}
		fmt.Printf("result:%s\n", obj.ToString())
	}

}
