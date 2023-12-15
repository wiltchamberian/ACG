package main

import (
	. "ACG/parser"
	"bufio"
	"fmt"
	"os"
)

func run() {
	//parse ebnf
	path := "./nika.gram"
	var parser = NewEBNFParser()
	rules := parser.Parse(path)

	//generate parser
	var generator Generator
	generator.SetOutputPath("./parser/nika_parser.go")
	name := "NikaParser"
	generator.Generate_rparser(name, rules)

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
	generator.SetOutputPath("./parser/nika_eval2.go")
	generator.Generate_eval3("NikaEval2", rules)
}

func Repl() {
	var userInput string
	var userLine string
	var eval NikaEval2
	var nika NikaParser
	reader := bufio.NewReader(os.Stdin)
	//fmt.Print(">>")
	for {
		userLine, _ = reader.ReadString('\n')
		if userLine != "\r\n" {
			// fmt.Printf("length:%d\n", len(userLine))
			// fmt.Println(userLine)
			// fmt.Println("<--not equal-->")
			userInput = userInput + userLine[0:len(userLine)-2]
		} else {
			// fmt.Println("<--start parsing-->")
			nika.ReadString(userInput)
			nika.TokenStream()
			inode, err := nika.PROG()
			if err != nil {
				fmt.Println("<--parser.PROG fail-->")
				userInput = ""
				//fmt.Print(">>")
				continue
			}
			obj := eval.Eval_nonterminal(inode)
			if obj == nil {
				fmt.Println("<--eval fail-->")
				userInput = ""
				//fmt.Print(">>")
				continue
			}
			fmt.Println(obj.ToString())
			userInput = ""
			//fmt.Print(">>")
		}

	}
}

func main() {

	//test2()
	Repl()

	var nika NikaParser
	var eval NikaEval2
	nika.ReadString("3+4;")
	nika.TokenStream()
	inode, err := nika.PROG()
	if err != nil {
		fmt.Println("<--parser.PROG fail-->")
	}
	//travel tree
	var travel Travel
	var printer NodePrinter
	printer.Init("./test_tree.txt")
	travel.DepthFirstTravel(inode, &printer)
	printer.Close()

	obj := eval.Eval_nonterminal(inode)
	if obj == nil {
		fmt.Println("<--eval fail-->")
	}
}
