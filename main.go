package main

import (
	. "ACG/parser"
	"fmt"
)

func run() {
	//parse ebnf
	path := "./nika.gram"
	var parser EBNFParser
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

	//generate parser
	var generator Generator
	generator.SetOutputPath("./parser/nika_parser.go")
	name := "NikaParser"
	generator.Generate_rparser(name, rules)

}
func main() {
	run()

	//parse code
	var nika NikaParser
	nika.ReadFile("./nika.nk")
	nika.TokenStream()
	root, _ := nika.PROG()

	//travel tree
	var travel Travel
	var printer NodePrinter
	printer.Init("./nika_parser_tree.txt")
	travel.DepthFirstTravel(root, &printer)
	printer.Close()
}
