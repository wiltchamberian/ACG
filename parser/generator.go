package parser

import (
	"bufio"
	"fmt"
	"os"
)

type Generator struct {
	id         int
	outputPath string
}

func (s *Generator) setOutputPath(path string) {
	s.outputPath = path
}

/*
generator作为关键组件，其逻辑本质比较简明，对于rules,首先生成
parser类名，然后每个rule对应一个函数
每个rule的每个alt生成一个if从句，然后alt里的每个item对应 从句中的一个expect语句
在从句结尾处(对应alt)生成parser树的节点，从句失败则不生成对应节点
*/
func (s *Generator) generate_rparser(rules []Rule) error {
	file, err := os.Create(s.outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprint(writer, "type RParser struct {\n")
	fmt.Fprint(writer, "\tBasicParser\n")
	fmt.Fprint(writer, "}\n\n")
	for _, rule := range rules {
		fmt.Fprintf(writer, "func (s *RParser) %s() {\n", rule.Name)
		fmt.Fprint(writer, "\tpos := s.Mark()\n")

		for _, alt := range rule.Alts {
			var variableNames []string
			varNamePrefix := "var"
			fmt.Fprint(writer, "\tok := true\n")
			fmt.Fprint(writer, "\tvar err error\n")
			fmt.Fprint(writer, "\tfor {\n")
			for _, item := range alt {
				varName := fmt.Sprintf("%s%d", varNamePrefix, len(variableNames))
				variableNames = append(variableNames, varName)
				//senmantic
				if item.Type.isType(TkString) || item.Type.isType(TkTerminator) {
					fmt.Fprintf(writer, "\t\t%s,err = s.Expect(Tk%s)\n", varName, string(item.Literal))
				} else {
					fmt.Fprintf(writer, "\t\t%s,err = s.%s()\n", varName, string(item.Literal))
				}
				fmt.Fprint(writer, "\t\tok = ok && (err == nil)\n")
				fmt.Fprint(writer, "\t\tif err != nil{\n\t\t\tbreak\n\t\t}\n")
			}
			fmt.Fprint(writer, "\tif ok == true {\n")
			//TODO
			fmt.Fprint(writer, "\t\t")
		}
	}
	writer.Flush()
	fmt.Println("data written to ", s.outputPath)
	return nil
}
