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
	var parser = NewBNFParser()
	parser.ReadFile(path)
	parser.TokenStream()
	tokens := parser.GetTokens()
	for _, token := range tokens {
		fmt.Println(string(token.Literal))
	}
	rules, err := parser.Parse()
	for _, rule := range rules {
		fmt.Println(rule.Name)
	}
	return rules, err
}

// func Repl() {
// 	var userInput string
// 	var userLine string
// 	var eval NikaEval2
// 	var nika NikaParser
// 	reader := bufio.NewReader(os.Stdin)
// 	//fmt.Print(">>")
// 	for {
// 		userLine, _ = reader.ReadString('\n')
// 		if userLine != "\r\n" {
// 			// fmt.Printf("length:%d\n", len(userLine))
// 			// fmt.Println(userLine)
// 			// fmt.Println("<--not equal-->")
// 			userInput = userInput + userLine[0:len(userLine)-2]
// 		} else {
// 			// fmt.Println("<--start parsing-->")
// 			nika.ReadString(userInput)
// 			nika.TokenStream()
// 			inode, err := nika.PROG()
// 			if err != nil {
// 				fmt.Println("<--parser.PROG fail-->")
// 				userInput = ""
// 				//fmt.Print(">>")
// 				nika.BasicParser.Clear()
// 				continue
// 			}
// 			obj := eval.Eval_nonterminal(inode)
// 			if obj == nil {
// 				fmt.Println("<--eval fail-->")
// 				userInput = ""
// 				//fmt.Print(">>")
// 				nika.BasicParser.Clear()
// 				continue
// 			}
// 			fmt.Println(obj.ToString())
// 			userInput = ""
// 			//fmt.Print(">>")
// 			nika.RBasicParser.Clear()
// 		}

// 	}
// }

// func TestNikaProgram() {
// 	var nika NikaParser
// 	var eval NikaEval2
// 	nika.ReadString("(3+4*5-7)*(10-2*3-3);")
// 	nika.TokenStream()
// 	inode, err := nika.PROG()
// 	if err != nil {
// 		fmt.Println("<--parser.PROG fail-->")
// 	}
// 	//travel tree
// 	var travel Travel
// 	var printer NodePrinter
// 	printer.Init("./test_tree.txt")
// 	travel.DepthFirstTravel(inode, &printer)
// 	printer.Close()

// 	obj := eval.Eval_nonterminal(inode)
// 	if obj == nil {
// 		fmt.Println("<--eval fail-->")
// 	}
// }

// func ReplVM() {
// 	var userInput string
// 	var userLine string

// 	reader := bufio.NewReader(os.Stdin)
// 	//fmt.Print(">>")
// 	for {
// 		userLine, _ = reader.ReadString('\n')
// 		if userLine != "\r\n" {
// 			userInput = userInput + userLine[0:len(userLine)-2]
// 		} else {
// 			var nika NikaParser
// 			nika.ReadString(userInput)
// 			nika.TokenStream()
// 			inode, err := nika.PROG()
// 			if err != nil {
// 				fmt.Printf("<--parse error-->\n")
// 				continue
// 			}
// 			var compiler NikaCompiler
// 			compiler.Compile(inode)
// 			vm := NewVM(compiler)
// 			vm.Run()
// 		}

// 	}
// }
