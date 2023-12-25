package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type Generator2 struct {
	id int
	FileWriter
}

// item can be terminator or nontermator
func (s *Generator2) generateItem(nodesName string, item *Token) {
	literal := string(item.Literal)
	//senmantic
	if item.Type.isType(TkString) {
		literal = fmt.Sprintf("\"%s\"", literal[1:len(literal)-1])
		s.Printf("Append(&%s,s.ExpectValueW(%s))", nodesName, literal)
	} else if unicode.IsUpper(rune(literal[0])) {
		s.Printf("Append(&%s,s.ExpectW(Tk%s))", nodesName, literal)
	} else {
		s.Printf("Append(&%s,s.%s())", nodesName, strings.ToTitle(literal))
	}
}

func (s *Generator2) PrintKleen(name string, id int, items []*Token) bool {
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

func (s *Generator2) PrintPara(name string, id int, items []*Token) {
	s.Printf("func (s *%s)Para_%d(nodes *[]INode) bool {\n", name, id)
	s.AddTab()
	s.Printf("pos := s.Mark()\n")
	s.Printf("var tmpNodes []INode\n")
	s.Printf("if(true ")
	for k := len(items) - 1; k >= 0; k-- {
		item := items[k]
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

func (s *Generator2) printKleens(name string, rules []Rule) {
	id := 0
	for _, rule := range rules {
		for i := len(rule.Alts) - 1; i >= 0; i-- {
			alt := rule.Alts[i]
			for j := len(alt.Groups) - 1; j >= 0; j-- {
				group := alt.Groups[j]
				if group.Type == GroupType_Kleen {
					s.PrintPara(name, id, group.Tokens)
					s.PrintKleen(name, id, group.Tokens)
					id++
				}
			}
		}
	}
}

func (s *Generator2) Generate_rparser(name string, bnf BNFRules) error {
	rules := bnf.Rules
	var err error
	err = s.OpenFile(s.outputPath)
	if err != nil {
		return err
	}
	defer s.CloseFile()

	s.Print("package parser\n\n")
	//s.Print("import \"fmt\"\n")
	s.Print("import \"errors\"\n")
	s.Print("import \"slices\"\n\n")

	s.Printf("type %s struct {\n", name)
	s.Printf("\tRBasicParser\n")
	s.Printf("}\n\n")

	s.printKleens(name, rules)

	id := 0
	for _, rule := range rules {
		s.Printf("func (s *%s) %s() Ret {\n", name, strings.ToTitle(rule.Name))
		s.AddTab()
		//for debug
		//s.Printf("fmt.Printf(\"%s\\n\")\n", rule.Name)
		s.Printf("var nodes []INode\n")
		s.Printf("pos := s.Mark()\n")
		for i := len(rule.Alts) - 1; i >= 0; i-- {
			alt := rule.Alts[i]
			s.Printf("if(true ")
			s.AddTab()
			for j := len(alt.Groups) - 1; j >= 0; j-- {
				group := alt.Groups[j]
				s.FPrintf(" &&\n")
				if group.Type == GroupType_Item {
					s.generateItem("nodes", alt.Groups[j].Tokens[0])
				} else if group.Type == GroupType_Kleen {
					s.Printf("s.Kleen_%d(&nodes)", id)
					id++
				}
			}
			s.SubTab()
			s.Printf("){\n")
			s.Printf("\tslices.Reverse(nodes)\n")
			s.Printf("\treturn Ret{&Node{\"%s\",nodes,nil,%d},nil}\n", string(rule.Name), i)
			s.Printf("} else {\n")
			s.Printf("\ts.Reset(pos)\n")
			s.Printf("\tnodes = []INode{}\n")
			s.Printf("}\n")
		}
		s.Printf("return Ret{nil, errors.New(\"None\")}\n")
		s.SubTab()
		s.Printf("}\n\n")
	}
	s.Flush()

	return nil
}
