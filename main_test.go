package main

import (
	. "ACG/parser"
	"encoding/binary"
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

func TestEmpty(t *testing.T) {
	data := make([]byte, 10)
	binary.LittleEndian.PutUint16(data[0:], 64+256)
	fmt.Printf("data0:%d", data[0])
	fmt.Printf("data1:%d", data[1])
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
	vm := NewStackMachine(compiler)
	vm.Print()
	vm.Run()

	answers := []int32{0, 428, 0, -64, 138, 35860, 3, -17, 36348, 33, 14, 32, 0, 1, 1, 0, 3, 1, 1, 0, 0, 0, 1, 3, 6, 10, 15, 21, 28, 36, 45, 45}
	fmt.Printf("answers_length:%d\n", len(answers))
	fmt.Printf("debugstack_length:%d\n", vm.DebugStack.Length())
	counter := 0
	for vm.DebugStack.Length() > 0 {
		obj := vm.DebugStack.PopInteger()
		if obj != answers[len(answers)-1-counter] {
			t.Fatalf("not equal")
		}
		fmt.Printf("result:%d\n", obj)
		counter++
	}

}
