package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type Generator struct {
	id int
	FileWriter
	ctl ParserControllor
}

func NewLGenerator() *Generator {
	var gen Generator
	gen.ctl = &ParserControllorL{}
	return &gen
}

func NewRGenerator() *Generator {
	var gen Generator
	gen.ctl = &ParserControllorR{}
	return &gen
}

// item can be terminator or nontermator
func (s *Generator) generateItem(nodesName string, item *Token) {
	literal := string(item.Literal)
	//senmantic
	if isType(item.Type, TkString) { //字符串终结符
		literal = fmt.Sprintf("\"%s\"", literal[1:len(literal)-1])
		s.Printf("Append(&%s,s.ExpectValueW(%s))", nodesName, literal)
	} else if literal == "Empty" { //空终结符
		//s.Printf("Append(&%s, NewEmpty())", nodesName)
		s.Printf("true") //不增加节点，因为后续也不直接用到该节点
	} else if unicode.IsUpper(rune(literal[0])) { //普通终结符
		s.Printf("Append(&%s,s.ExpectW(Tk%s))", nodesName, literal)
	} else {
		s.Printf("Append(&%s,s.%s())", nodesName, strings.ToTitle(literal))
	}
}

func (s *Generator) PrintKleen(name string, id int, items []*Token) bool {
	s.Printf("func (s *%s)Kleen_%d(nodes *[]INode) bool {\n", name, id)
	s.AddTab()
	s.Printf("for {\n")
	s.AddTab()
	s.Printf("ok := s.Para_%d(nodes)\n", id)
	s.Printf("if ok==false {\n")
	s.Printf("\tbreak\n")
	s.Printf("}\n")
	s.SubTab()
	s.Printf("}\n")
	s.Printf("return true\n")
	s.SubTab()
	s.Printf("}\n\n")
	return true
}

func (s *Generator) PrintPara(name string, id int, iter Iterator[*Token]) {
	s.Printf("func (s *%s)Para_%d(nodes *[]INode) bool {\n", name, id)
	s.AddTab()
	s.Printf("pos := s.Mark()\n")
	s.Printf("var tmpNodes []INode\n")
	s.Printf("if(true ")
	for k := iter.Start(); iter.Legal(k); k = iter.Advance(k) {
		item := iter.Get(k)
		literal := string(item.Literal)
		if literal == "|" {
			s.FPrintf("){\n")
			s.Printf("\t*nodes = append(*nodes, tmpNodes...)\n")
			s.Printf("\treturn true\n")
			s.Printf("} else{\n")
			s.Printf("\ts.Reset(pos)\n")
			s.Printf("\ttmpNodes = []INode{}\n")
			s.Printf("}\n\n")
			s.Printf("if(true ")
		} else {
			s.FPrintf("&&\n")
			s.generateItem("tmpNodes", item)
		}
	}
	s.FPrintf("){\n")
	s.Printf("\t*nodes = append(*nodes, tmpNodes...)\n")
	s.Printf("\treturn true\n")
	s.Printf("} else{\n")
	s.Printf("\ts.Reset(pos)\n")
	s.Printf("\ttmpNodes = []INode{}\n")
	s.Printf("}\n\n")

	s.Printf("return false\n")
	s.SubTab()
	s.Printf("}\n\n")
}

func (s *Generator) printKleensR(name string, rules []Rule) {
	id := 0
	for _, rule := range rules {
		for i := len(rule.Alts) - 1; i >= 0; i-- {
			alt := rule.Alts[i]
			for j := len(alt.Groups) - 1; j >= 0; j-- {
				group := alt.Groups[j]
				if group.Type == GroupType_Kleen {
					iter := ArrayRIter[*Token]{group.Tokens}
					s.PrintPara(name, id, iter)
					s.PrintKleen(name, id, group.Tokens)
					id++
				}
			}
		}
	}
}

func (s *Generator) printKleensL(name string, rules []Rule) {
	id := 0
	for _, rule := range rules {
		for i := 0; i < len(rule.Alts); i++ {
			alt := rule.Alts[i]
			for j := 0; j < len(alt.Groups); j++ {
				group := alt.Groups[j]
				if group.Type == GroupType_Kleen {
					iter := ArrayIter[*Token]{group.Tokens}
					s.PrintPara(name, id, iter)
					s.PrintKleen(name, id, group.Tokens)
					id++
				}
			}
		}
	}
}

func (s *Generator) GenerateRParser(name string, bnf BNFRules) error {
	s.ctl = &ParserControllorR{}
	return s.Generate_rparser(name, bnf)
}

func (s *Generator) GenerateLParser(name string, bnf BNFRules) error {
	s.ctl = &ParserControllorL{}
	return s.Generate_lparser(name, bnf)
}

func (s *Generator) PrintRuleStart(name string, rule Rule) {
	s.Printf("func (s *%s) %s() Ret {\n", name, strings.ToTitle(rule.Name))
	s.AddTab()
	//for debug
	//s.Printf("fmt.Printf(\"%s\\n\")\n", rule.Name)
	s.Printf("var nodes []INode\n")
	s.Printf("pos := s.Mark()\n")
}

func (s *Generator) PrintRuleEnd() {
	s.Printf("return Ret{nil, errors.New(\"None\")}\n")
	s.SubTab()
	s.Printf("}\n\n")
}

func (s *Generator) PrintGroup(group Group, id *int) {
	s.FPrintf(" &&\n")
	if group.Type == GroupType_Item {
		s.generateItem("nodes", group.Tokens[0])
	} else if group.Type == GroupType_Kleen {
		s.Printf("s.Kleen_%d(&nodes)", *id)
		(*id)++
	}
}

func (s *Generator) PrintStructHead(name string) {
	s.Print("package parser\n\n")
	//s.Print("import \"fmt\"\n")
	s.Print("import \"errors\"\n")
	s.ctl.PrintImportSlice(s)

	s.Printf("type %s struct {\n", name)
	s.Printf("\tBasicParser\n")
	s.Printf("}\n\n")

	s.ctl.PrintNewParser(s, name)
}

func (s *Generator) PrintAltEnd(rule Rule, i int) {
	s.ctl.Print(s)
	s.Printf("\tif(len(nodes)==1 && (!nodes[0].IsTerminal())){\n")
	s.Printf("\t\treturn Ret{nodes[0],nil}\n")
	s.Printf("\t} else{\n")
	s.Printf("\t\treturn Ret{&Node{\"%s\",nodes,nil,%d,\"%s\"},nil}\n", string(rule.Name), i, rule.Alts[i].action)

	s.Printf("\t}\n")
	s.Printf("} else {\n")
	s.Printf("\ts.Reset(pos)\n")
	s.Printf("\tnodes = []INode{}\n")
	s.Printf("}\n")

	// //reverse nodes
	// s.ctl.Print(s)
	// s.Printf("\t\treturn Ret{&Node{\"%s\",nodes,nil,%d,\"%s\"},nil}\n", string(rule.Name), i, rule.Alts[i].action)
	// s.Printf("} else {\n")
	// s.Printf("\ts.Reset(pos)\n")
	// s.Printf("\tnodes = []INode{}\n")
	// s.Printf("}\n")

}

func (s *Generator) Generate_rparser(name string, bnf BNFRules) error {
	rules := bnf.Rules
	err := s.OpenFile(s.outputPath)
	if err != nil {
		return err
	}
	defer s.CloseFile()

	s.PrintStructHead(name)
	s.printKleensR(name, rules)
	id := 0
	for _, rule := range rules {
		s.PrintRuleStart(name, rule)
		for i := len(rule.Alts) - 1; i >= 0; i-- {
			alt := rule.Alts[i]
			s.Printf("if(true ")
			s.AddTab()
			for j := len(alt.Groups) - 1; j >= 0; j-- {
				s.PrintGroup(alt.Groups[j], &id)
			}
			s.SubTab()
			s.Printf("){\n")
			s.PrintAltEnd(rule, i)
		}
		s.PrintRuleEnd()
	}
	s.Flush()

	return nil
}

func (s *Generator) Generate_lparser(name string, bnf BNFRules) error {
	rules := bnf.Rules
	err := s.OpenFile(s.outputPath)
	if err != nil {
		return err
	}
	defer s.CloseFile()

	s.PrintStructHead(name)
	s.printKleensL(name, rules)
	id := 0
	for _, rule := range rules {
		s.PrintRuleStart(name, rule)
		for i := 0; i < len(rule.Alts); i++ {
			alt := rule.Alts[i]
			s.Printf("if(true ")
			s.AddTab()
			for j := 0; j < len(alt.Groups); j++ {
				s.PrintGroup(alt.Groups[j], &id)
			}
			s.SubTab()
			s.Printf("){\n")
			s.PrintAltEnd(rule, i)
		}
		s.PrintRuleEnd()
	}
	s.Flush()

	return nil
}

func (s *Generator) Generate_eval(name string, bnf BNFRules) error {
	rules := bnf.Rules
	var err error
	err = s.OpenFile(s.outputPath)
	if err != nil {
		return err
	}
	defer s.CloseFile()
	writer := s.writer

	fmt.Fprint(writer, "package parser\n\n")
	fmt.Fprint(writer, "import \"strconv\"\n\n")
	//fmt.Fprint(writer, "import \"errors\"\n")
	//fmt.Fprint(writer, "import \"slices\"\n\n")
	fmt.Fprintf(writer, "type %s struct {\n", name)
	fmt.Fprint(writer, "}\n\n")

	s.ctl.PrintNewParser(s, name)

	fmt.Fprintf(writer, "func (s *%s) GetActionResult(nd INode, v []NkObject) (result NkObject){\n", name)
	fmt.Fprint(writer, "\tname := nd.GetName()\n")
	fmt.Fprint(writer, "\tswitch name {\n")
	for _, rule := range rules {
		fmt.Fprintf(writer, "\t\tcase \"%s\":{\n", rule.Name)
		fmt.Fprint(writer, "\t\t\tswitch nd.Select() {\n")
		for i, alt := range rule.Alts {
			fmt.Fprintf(writer, "\t\t\t\tcase %d:{\n", i)
			if alt.action != "" {
				fmt.Fprintf(writer, "\t\t\t\t\tresult = %s\n", alt.action)
			} else {
				fmt.Fprintf(writer, "\t\t\t\t\tresult = v[0]\n")
			}
			fmt.Fprintf(writer, "\t\t\t\t}\n")
		}
		fmt.Fprint(writer, "\t\t\t}\n")
		fmt.Fprint(writer, "\t\t}\n")
	}
	fmt.Fprint(writer, "\t\tdefault:{\n")
	fmt.Fprint(writer, "\t\t\tresult = v[0]\n")
	fmt.Fprint(writer, "\t\t}\n")
	fmt.Fprint(writer, "\t}\n")
	fmt.Fprint(writer, "\treturn\n")
	fmt.Fprint(writer, "}\n\n")

	fmt.Fprintf(writer, "func (s *%s) Eval_terminal(nd INode) (result NkObject){\n", name)
	fmt.Fprintf(writer, "\ttoken, ok := nd.(*Token)\n")
	fmt.Fprintf(writer, "\tif ok == false {\n")
	fmt.Fprintf(writer, "\t\treturn\n")
	fmt.Fprintf(writer, "\t}\n")
	fmt.Fprintf(writer, "\tif token.Type.isType(TkNumber) {\n")
	fmt.Fprintf(writer, "\t\tv,e := strconv.Atoi(string(token.Literal))\n")
	fmt.Fprintf(writer, "\t\tif e== nil {\n")
	fmt.Fprintf(writer, "\t\t\treturn &NkInteger{Value: v}\n")
	fmt.Fprintf(writer, "\t\t}\n")
	fmt.Fprintf(writer, "\t}\n")
	fmt.Fprintf(writer, "\treturn nil\n")
	fmt.Fprintf(writer, "}\n\n")

	fmt.Fprintf(writer, "func (s *%s) Eval_nonterminal(nd INode) (result NkObject){\n", name)
	fmt.Fprint(writer, "\tchildren:=nd.GetChildren()\n")
	fmt.Fprint(writer, "\tvar v []NkObject\n")
	fmt.Fprintf(writer, "\tfor _, child := range children {\n")
	fmt.Fprintf(writer, "\t\tif child.IsTerminal(){\n")
	fmt.Fprintf(writer, "\t\t\t v = append(v, s.Eval_terminal(child))\n")
	fmt.Fprintf(writer, "\t\t} else {\n")
	fmt.Fprintf(writer, "\t\t\t v = append(v,s.Eval_nonterminal(child))\n")
	fmt.Fprintf(writer, "\t\t}\n")
	fmt.Fprintf(writer, "\t}\n")

	fmt.Fprintf(writer, "\tresult = s.GetActionResult(nd,v)\n")
	fmt.Fprintf(writer, "\treturn\n")
	fmt.Fprint(writer, "}\n\n")
	writer.Flush()

	return nil
}

func (s *Generator) PrintCompiler(name string, bnf BNFRules) error {
	rules := bnf.Rules
	var err error
	err = s.OpenFile(s.outputPath)
	if err != nil {
		return err
	}
	defer s.CloseFile()

	s.Printf("package parser\n\n")
	//s.Printf("import \"strconv\"\n\n")
	//s.Printf("import \"errors\"\n")
	s.Printf("type %s struct {\n", name)
	s.Printf("\tCompiler\n")
	s.Print("}\n\n")

	s.Printf("func New%s() *%s{\n", name, name)
	s.Printf("\tvar result %s\n", name)
	s.Printf("\tresult.Compiler.InitCompiler()\n")
	s.Printf("\treturn &result\n")
	s.Printf("}\n\n")
	s.writer.Flush()

	s.Printf("func (s *%s) C(node INode) bool {\n", name)
	s.Printf("\terr := s.Compile(node)\n")
	s.Printf("\treturn err == nil\n")
	s.Print("}\n\n")

	s.Printf("func (s *%s) Compile(node INode) (err error) {\n", name)
	s.AddTab()
	s.Printf("switch node.GetName() {\n")
	s.AddTab()
	for _, rule := range rules {
		s.Printf("case \"%s\":{\n", rule.Name)
		s.AddTab()
		s.Printf("switch node.Select(){\n")
		s.AddTab()
		for i, alt := range rule.Alts {
			s.Printf("case %d:{\n", i)
			s.AddTab()
			if alt.action != "" {
				s.Printf("v := node.GetChildren()\n")
				s.Printf("l := len(v)\n")
				s.Printf("%s\n", alt.action)
				s.Printf("Do(v, l)\n") //shut up the compiler!
			}
			s.SubTab()
			s.Printf("}\n")
		}
		s.SubTab()
		s.Printf("}\n")
		s.SubTab()
		s.Printf("}\n")
	}
	s.SubTab()
	s.Printf("}\n")
	s.Printf("return nil\n")
	s.SubTab()
	s.Printf("}\n")

	s.writer.Flush()

	return nil
}
