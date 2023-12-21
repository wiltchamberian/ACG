package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Generator struct {
	id         int
	outputPath string
	file       *os.File
	writer     *bufio.Writer
	tabs       string
}

func (s *Generator) SetOutputPath(path string) {
	s.outputPath = path
}

func (s *Generator) ResetTab() {
	s.tabs = ""
}

func (s *Generator) OpenFile(path string) error {
	var err error
	s.file, err = os.Create(s.outputPath)
	if err != nil {
		return err
	}
	s.writer = bufio.NewWriter(s.file)
	return nil
}

func (s *Generator) CloseFile() {
	s.file.Close()
}

/*
generator作为关键组件，其逻辑本质比较简明，对于rules,首先生成
parser类名，然后每个rule对应一个函数
每个rule的每个alt生成一个if从句，然后alt里的每个item对应 从句中的一个expect语句
在从句结尾处(对应alt)生成parser树的节点，从句失败则不生成对应节点
*/
func (s *Generator) Generate_rparser(name string, rules []Rule) error {
	var err error
	err = s.OpenFile(s.outputPath)
	if err != nil {
		return err
	}
	defer s.CloseFile()
	writer := s.writer

	fmt.Fprint(writer, "package parser\n\n")
	fmt.Fprint(writer, "import \"errors\"\n")
	fmt.Fprint(writer, "import \"slices\"\n\n")
	fmt.Fprintf(writer, "type %s struct {\n", name)
	fmt.Fprint(writer, "\tRBasicParser\n")
	fmt.Fprint(writer, "}\n\n")
	for _, rule := range rules {
		fmt.Fprintf(writer, "func (s *%s) %s() (INode, error){\n", name, strings.ToTitle(rule.Name))
		fmt.Fprint(writer, "\tpos := s.Mark()\n")

		var counter int = 0
		for alt_index, alt := range rule.Alts {
			//var variableNames []string
			//varNamePrefix := "var"
			fmt.Fprint(writer, "\t//alternative\n")
			fmt.Fprint(writer, "\t{\n")
			fmt.Fprint(writer, "\t\tok := true\n")
			// for i, _ := range alt {
			// 	fmt.Fprintf(writer, "\t\tvar var%d INode\n", i)
			// }
			fmt.Fprint(writer, "\t\tvar err error\n")
			fmt.Fprint(writer, "\t\tvar nodes []INode\n")
			fmt.Fprint(writer, "\t\tvar node INode\n\n")

			for i := len(alt.Groups) - 1; i >= 0; i-- {
				if alt.Groups[i].Type == 1 {
					fmt.Fprint(writer, "\t\tpos2 := s.Mark()\n")
					fmt.Fprint(writer, "\t\tvar rpos = pos2\n")
					fmt.Fprint(writer, "\t\tvar arr []INode\n")
					break
				}
			}

			//reverse travel
			for i := len(alt.Groups) - 1; i >= 0; i-- {
				fmt.Fprint(writer, "\t\t//Group\n")
				if alt.Groups[i].Type == GroupType_Item {
					item := alt.Groups[i].Tokens[0]
					// varName := fmt.Sprintf("%s%d", varNamePrefix, len(variableNames))
					// variableNames = append(variableNames, varName)
					s.generateItem(item, "node", 2)
					fmt.Fprint(writer, "\t\tnodes = append(nodes,node)\n")
					fmt.Fprint(writer, "\t\tok = ok && (err == nil)\n")
					fmt.Fprintf(writer, "\t\tif err != nil{\n\t\t\tgoto LABEL%d\n\t\t}\n", counter)
				} else if alt.Groups[i].Type == GroupType_Kleen {
					items := alt.Groups[i].Tokens
					fmt.Fprint(writer, "\t\tpos2 = s.Mark()\n")
					fmt.Fprint(writer, "\t\trpos = pos2\n")
					fmt.Fprint(writer, "\t\tarr = []INode{}\n")
					fmt.Fprint(writer, "\t\tfor{\n")
					fmt.Fprint(writer, "\t\t\trpos = s.Mark()\n")
					for j := len(items) - 1; j >= 0; j-- {
						item := items[j]
						s.generateItem(item, "node", 3)
						fmt.Fprintf(writer, "\t\t\tif err != nil{\n\t\t\t\tbreak\n\t\t\t}else{\n")
						fmt.Fprintf(writer, "\t\t\t\tarr = append(arr,node)\n\t\t\t}\n")
					}
					fmt.Fprint(writer, "\t\t}\n")
					fmt.Fprint(writer, "\t\tif(s.Mark() - rpos > 0){\n")
					fmt.Fprint(writer, "\t\t\ts.Reset(pos2)\n")
					fmt.Fprint(writer, "\t\t}else{\n")

					fmt.Fprint(writer, "\t\t\tnodes = append(nodes,arr...)\n\t\t}\n")
					//fmt.Fprint(writer, "\t\t\tvar newNode INode = &Node{Name:\"__kleene\",Children:arr}\n")
					//fmt.Fprint(writer, "\t\t\tnodes = append(nodes, newNode)\n\t\t}\n")
				} else {

				}
				fmt.Fprint(writer, "\t\t\n")
			}
			//shut up the compiler
			fmt.Fprintf(writer, "\t\tgoto LABEL%d\n", counter)
			fmt.Fprintf(writer, "LABEL%d:\n", counter)
			counter++
			fmt.Fprint(writer, "\t\tif ok == true {\n")

			// slices.Reverse(variableNames)
			// fmt.Fprintf(writer, "\t\t\treturn &Node{\"%s\",[]INode{%s}},nil\n", string(rule.Name), strings.Join(variableNames, ", "))
			fmt.Fprint(writer, "\t\t\tslices.Reverse(nodes)\n")
			fmt.Fprintf(writer, "\t\t\treturn &Node{\"%s\",nodes,nil,%d},nil\n", string(rule.Name), alt_index)

			fmt.Fprint(writer, "\t\t}\n")
			fmt.Fprint(writer, "\t\ts.Reset(pos)\n")

			fmt.Fprint(writer, "\t}\n")
		}
		fmt.Fprint(writer, "\treturn nil, errors.New(\"None\")\n")
		fmt.Fprint(writer, "}\n\n")
	}
	writer.Flush()
	fmt.Println("data written to ", s.outputPath)
	return nil
}

func (s *Generator) generateItem(item *Token, varName string, level int) {
	literal := string(item.Literal)
	//senmantic
	if item.Type.isType(TkString) {
		literal = fmt.Sprintf("\"%s\"", literal[1:len(literal)-1])
		s.printTabs(level)
		fmt.Fprintf(s.writer, "%s,err = s.ExpectValue(%s)\n", varName, literal)
	} else {
		if unicode.IsUpper(rune(literal[0])) {
			s.printTabs(level)
			fmt.Fprintf(s.writer, "%s,err = s.Expect(Tk%s)\n", varName, literal)
		} else {
			s.printTabs(level)
			fmt.Fprintf(s.writer, "%s,err = s.%s()\n", varName, strings.ToTitle(literal))
		}
	}
}

func (s *Generator) Generate_rparser2(name string, rules []Rule) error {

	return nil
}

func (s *Generator) printTabs(level int) {
	for i := 0; i < level; i++ {
		fmt.Fprint(s.writer, "\t")
	}
}

func (s *Generator) Generate_eval(name string, rules []Rule) error {
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

func (s *Generator) Printf(format string, a ...any) (n int, err error) {
	fmt.Fprint(s.writer, s.tabs)
	return fmt.Fprintf(s.writer, format, a...)
}

func (s *Generator) Print(a ...any) (n int, err error) {
	fmt.Fprint(s.writer, s.tabs)
	return fmt.Fprint(s.writer, a...)
}

func (s *Generator) AddTab() {
	s.tabs = s.tabs + "\t"
}

func (s *Generator) SubTab() {
	s.tabs = s.tabs[0 : len(s.tabs)-1]
}

func (s *Generator) GenerateIR(name string, rules []Rule) error {
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
	s.writer.Flush()

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
				s.Printf("%s\n", alt.action)
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
