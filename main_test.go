package main

import (
	. "ACG/parser"
	"fmt"
	"testing"
)

func TestHelloEmpty(t *testing.T) {
	var prog = `(3+ 4*(8-2)) < 0+13;
2*((3*(4-5+13)/2 *13-7)-13);
77*13==15;`

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

}
