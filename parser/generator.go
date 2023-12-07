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
func (s *Generator) Generate_rparser(rules []Rule) error {
	file, err := os.Create(s.outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprint(writer, "package parser\n\n")
	fmt.Fprint(writer, "import \"errors\"\n\n")
	fmt.Fprint(writer, "type RParser struct {\n")
	fmt.Fprint(writer, "\tRBasicParser\n")
	fmt.Fprint(writer, "}\n\n")
	for _, rule := range rules {
		fmt.Fprintf(writer, "func (s *RParser) %s() (INode, error){\n", strings.ToTitle(rule.Name))
		fmt.Fprint(writer, "\tpos := s.Mark()\n")

		var counter int = 0
		for _, alt := range rule.Alts {
			var variableNames []string
			varNamePrefix := "var"
			fmt.Fprint(writer, "\t{\n")
			fmt.Fprint(writer, "\t\tok := true\n")
			for i, _ := range alt {
				fmt.Fprintf(writer, "\t\tvar var%d INode\n", i)
			}
			fmt.Fprint(writer, "\t\tvar err error\n")
			//reverse travel
			for i := len(alt) - 1; i >= 0; i-- {
				item := alt[i]
				varName := fmt.Sprintf("%s%d", varNamePrefix, len(variableNames))
				variableNames = append(variableNames, varName)
				literal := string(item.Literal)
				//senmantic
				if item.Type.isType(TkString) {
					literal = fmt.Sprintf("\"%s\"", literal[1:len(literal)-1])
					fmt.Fprintf(writer, "\t\t%s,err = s.ExpectValue(%s)\n", varName, literal)
				} else {
					if unicode.IsUpper(rune(literal[0])) {
						fmt.Fprintf(writer, "\t\t%s,err = s.Expect(Tk%s)\n", varName, literal)
					} else {
						fmt.Fprintf(writer, "\t\t%s,err = s.%s()\n", varName, strings.ToTitle(literal))
					}
				}
				fmt.Fprint(writer, "\t\tok = ok && (err == nil)\n")
				fmt.Fprintf(writer, "\t\tif err != nil{\n\t\t\tgoto END%d\n\t\t}\n", counter)
			}
			fmt.Fprintf(writer, "END%d:\n", counter)
			counter++
			fmt.Fprint(writer, "\t\tif ok == true {\n")

			variableNames = ReverseSlice(variableNames)
			fmt.Fprintf(writer, "\t\t\treturn &Node{\"%s\",[]INode{%s}},nil\n", string(rule.Name), strings.Join(variableNames, ", "))
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
