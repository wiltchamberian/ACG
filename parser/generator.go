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
}

func (s *Generator) SetOutputPath(path string) {
	s.outputPath = path
}

/*
generator作为关键组件，其逻辑本质比较简明，对于rules,首先生成
parser类名，然后每个rule对应一个函数
每个rule的每个alt生成一个if从句，然后alt里的每个item对应 从句中的一个expect语句
在从句结尾处(对应alt)生成parser树的节点，从句失败则不生成对应节点
*/
func (s *Generator) Generate_rparser(name string, rules []Rule) error {
	var err error
	s.file, err = os.Create(s.outputPath)
	if err != nil {
		return err
	}
	defer s.file.Close()
	writer := bufio.NewWriter(s.file)
	s.writer = writer
	fmt.Fprint(writer, "package parser\n\n")
	fmt.Fprint(writer, "import \"errors\"\n\n")
	fmt.Fprintf(writer, "type %s struct {\n", name)
	fmt.Fprint(writer, "\tRBasicParser\n")
	fmt.Fprint(writer, "}\n\n")
	for _, rule := range rules {
		fmt.Fprintf(writer, "func (s *%s) %s() (INode, error){\n", name, strings.ToTitle(rule.Name))
		fmt.Fprint(writer, "\tpos := s.Mark()\n")

		var counter int = 0
		for _, alt := range rule.Alts {
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
			//reverse travel
			for i := len(alt) - 1; i >= 0; i-- {
				fmt.Fprint(writer, "\t\t//Group\n")
				if alt[i].Type == 0 {
					item := alt[i].Tokens[0]
					// varName := fmt.Sprintf("%s%d", varNamePrefix, len(variableNames))
					// variableNames = append(variableNames, varName)
					s.generateItem(item, "node", 2)
					fmt.Fprint(writer, "\t\tnodes = append(nodes,node)\n")
					fmt.Fprint(writer, "\t\tok = ok && (err == nil)\n")
					fmt.Fprintf(writer, "\t\tif err != nil{\n\t\t\tgoto LABEL%d\n\t\t}\n", counter)
				} else if alt[i].Type == 1 {
					items := alt[i].Tokens
					fmt.Fprint(writer, "\t\tpos2 := s.Mark()\n")
					fmt.Fprint(writer, "\t\tvar rpos = pos2\n")
					fmt.Fprint(writer, "\t\tvar arr []INode\n")
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
				}
				fmt.Fprint(writer, "\t\t\n")
			}
			//shut up the compiler
			fmt.Fprintf(writer, "\t\tgoto LABEL%d\n", counter)
			fmt.Fprintf(writer, "LABEL%d:\n", counter)
			counter++
			fmt.Fprint(writer, "\t\tif ok == true {\n")

			// variableNames = ReverseSlice(variableNames)
			// fmt.Fprintf(writer, "\t\t\treturn &Node{\"%s\",[]INode{%s}},nil\n", string(rule.Name), strings.Join(variableNames, ", "))
			fmt.Fprint(writer, "\t\t\tReverseSliceInPlace(nodes)\n")
			fmt.Fprintf(writer, "\t\t\treturn &Node{\"%s\",nodes},nil\n", string(rule.Name))

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

func (s *Generator) printTabs(level int) {
	for i := 0; i < level; i++ {
		fmt.Fprint(s.writer, "\t")
	}
}
