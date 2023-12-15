package main

import (
	. "ACG/parser"
	"fmt"
	"slices"
)

type ITest interface {
	run()
	roll()
}

type ATest struct {
}

func (s *ATest) run() {
	fmt.Println("ATest")
	s.roll()
}

func (s *ATest) roll() {
	fmt.Println("ATest_roll")
}

type BTest struct {
	ATest
}

func (s *BTest) roll() {
	fmt.Println("BTest_roll")
}

func tt() {
	var b ITest = &BTest{}
	r, er := b.(*ATest)
	if er != true {
		panic("")
	}
	r.run()

	var arr = []int{1, 2, 3, 4}
	var bb = arr
	bb = slices.Delete(bb, 1, 2)
	for i, _ := range arr {
		fmt.Printf("%d,", arr[i])
	}
}

func test() {
	var abc []byte = []byte("1234")
	q := abc[0:2]
	abc[0] = 'x'

	fmt.Println(string(abc))
	fmt.Println(string(q))
}

func testLexer() {
	var content string = `x = 13;
	y = 7 + 4;
	z = x * 3;
	q = 'abc123'`

	var lexer Lexer
	err := lexer.ReadString(content)
	if err != nil {
		panic("readstring error")
	}

	tokens, err := lexer.TokenStream()

	for _, token := range tokens {
		fmt.Println(string(token.Literal))
	}
}

func testEbnfParser() ([]Rule, error) {
	path := "./nika_simple.gram"
	var parser = NewEBNFParser()
	parser.ReadFile(path)
	parser.TokenStream()
	tokens := parser.GetTokens()
	for _, token := range tokens {
		fmt.Println(string(token.Literal))
	}
	rules := parser.Grammar()
	for _, rule := range rules {
		fmt.Println(rule.Name)
	}
	return rules, nil
}

func testGenerator() {
	rules, _ := testEbnfParser()
	var generator Generator
	generator.SetOutputPath("./parser/nika_parser.go")
	name := "NikaParser"
	generator.Generate_rparser(name, rules)

}
