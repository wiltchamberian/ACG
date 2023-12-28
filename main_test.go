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
	var prog = `(3+ 4*(8-2)) < 0+13;
2*((3*(4-5+13)/2 *13-7)-13);
77*13==15;
(13-29)*(10-2)/2;
123-10/(3-1)+20;
(13*20 + (((127-123)*(5+18*4)-27)*13-(25+29*10)))*10-120;
(1||0&&1) + (1||1) + (0||1) + (0&&1);
-13-+5----4+---(---(---3));
((3+ 4*(8-2)) < 0+13)+
	2*((3*(4-5+13)/2 *13-7)-13) + 
	(77*13==15) +
	(13-29)*(10-2)/2 +
	(123-10/(3-1)+20) +
	((13*20 + (((127-123)*(5+18*4)-27)*13-(25+29*10)))*10-120)+
	((1||0&&1) + (1||1) + (0||1) + (0&&1))+
	(-13-+5----4+---(---(---3)));
`
	answers := []int{36348, -17, 3, 35860, 138, -64, 0, 428, 0}

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
	vm.Run()
	i := 0
	for vm.IsEmpty() == false {
		obj := vm.Pop()
		if obj != nil {
			fmt.Printf("result:%s\n", obj.ToString())
		}
		if obj.(*NkInteger).Value != answers[i] {
			t.Fatalf("%d not equal %d", obj.(*NkInteger).Value, answers[i])
		}
		i++
	}

}
