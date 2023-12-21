package main

import (
	. "ACG/parser"
)

func run() {
	//parse ebnf
	path := "./nika.gram"
	var parser = NewEBNFParser()
	rules := parser.Parse(path)

	//generate parser
	var generator Generator
	generator.SetOutputPath("./parser/nika_parser.go")
	generator.Generate_rparser("NikaParser", rules)

}
func test1() {
	run()
	//parse code
	var nika NikaParser
	nika.ReadFile("./nika.nk")
	nika.TokenStream()
	tree, _ := nika.PROG()

	//travel tree
	var travel Travel
	var printer NodePrinter
	printer.Init("./nika_parser_tree.txt")
	travel.DepthFirstTravel(tree, &printer)
	printer.Close()

	//Parse Tree To AST
	ToAST(tree)
	printer.Init("./nika_ast.txt")
	travel.DepthFirstTravel(tree, &printer)
	printer.Close()
}

func test2() {
	path := "./nika.gram"
	var parser = NewEBNFParser()
	rules := parser.Parse(path)

	//generate parser
	var generator Generator
	generator.SetOutputPath("./parser/nika_eval.go")
	generator.Generate_eval("NikaEval", rules)
}

func main() {
	//testAll()
	//Repl()
	TestRunVM()

	return
}
