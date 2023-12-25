package main

import (
	. "ACG/parser"
)

func run() {
	//parse ebnf
	path := "./nika.gram"
	var parser = NewBNFParser()
	rules, _ := parser.ParseFile(path)

	//generate parser
	var generator Generator
	generator.SetOutputPath("./parser/nika_parser.go")
	generator.Generate_rparser("NikaParser", rules)

}

// func test1() {
// 	run()
// 	//parse code
// 	var nika NikaParser
// 	nika.ReadFile("./nika.nk")
// 	nika.TokenStream()
// 	tree, _ := nika.PROG()

// 	//travel tree
// 	var travel Travel
// 	var printer NodePrinter
// 	printer.Init("./nika_parser_tree.txt")
// 	travel.DepthFirstTravel(tree, &printer)
// 	printer.Close()

// 	//Parse Tree To AST
// 	ToAST(tree)
// 	printer.Init("./nika_ast.txt")
// 	travel.DepthFirstTravel(tree, &printer)
// 	printer.Close()
// }

func test2() {
	path := "./nika.gram"
	var parser = NewBNFParser()
	bnf, _ := parser.ParseFile(path)

	//generate parser
	var generator Generator
	generator.SetOutputPath("./parser/nika_eval.go")
	generator.Generate_eval("NikaEval", bnf)
}

func abc() (int, bool) {
	return 0, true
}

func main() {
	//testAll()
	//Repl()

	TestNewGenerator()
	//TestRunVM()

	return
}
