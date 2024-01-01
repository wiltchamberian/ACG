package parser

import (
	"bufio"
	"fmt"
	"os"
)

type FileWriter struct {
	outputPath string
	file       *os.File
	writer     *bufio.Writer
	tabs       string
}

func (s *FileWriter) SetOutputPath(path string) {
	s.outputPath = path
}

func (s *FileWriter) ResetTab() {
	s.tabs = ""
}

func (s *FileWriter) OpenFile(path string) error {
	var err error
	s.file, err = os.Create(path)
	if err != nil {
		return err
	}
	s.writer = bufio.NewWriter(s.file)
	return nil
}

func (s *FileWriter) CloseFile() {
	s.writer.Flush()
	s.file.Close()
}

func (s *FileWriter) Printf(format string, a ...any) (n int, err error) {
	fmt.Fprint(s.writer, s.tabs)
	return fmt.Fprintf(s.writer, format, a...)
}

// print without tag
func (s *FileWriter) FPrintf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(s.writer, format, a...)
}

func (s *FileWriter) Print(a ...any) (n int, err error) {
	fmt.Fprint(s.writer, s.tabs)
	return fmt.Fprint(s.writer, a...)
}

func (s *FileWriter) FPrint(a ...any) (n int, err error) {
	return fmt.Fprint(s.writer, a...)
}

func (s *FileWriter) AddTab() {
	s.tabs = s.tabs + "\t"
}

func (s *FileWriter) SubTab() {
	s.tabs = s.tabs[0 : len(s.tabs)-1]
}

func (s *FileWriter) printTabs(level int) {
	for i := 0; i < level; i++ {
		fmt.Fprint(s.writer, "\t")
	}
}

func (s *FileWriter) Flush() error {
	return s.writer.Flush()
}

func (s *FileWriter) PrintStructHead(name string) {
	s.Printf("type %s struct{\n", name)
	s.AddTab()
}

func (s *FileWriter) PrintStructEnd() {
	s.SubTab()
	s.Print("}\n\n")
}

func (s *FileWriter) PrintInterfaceHead(name string) {
	s.Printf("type %s interface{\n", name)
	s.AddTab()
}

func (s *FileWriter) PrintInterfaceEnd() {
	s.SubTab()
	s.Print("}\n\n")
}

func (s *FileWriter) PrintEnd() {
	s.Print("}\n\n")
}

func (s *FileWriter) PrintPackage(name string) {
	s.Printf("package %s\n\n", name)
}

// we can't input two variable-length parameter
// so we use lambda calculus
func (s *FileWriter) PrintFuncHead(name string) func(...string) func(...string) {
	return func(parameters ...string) func(...string) {
		return func(returnTypes ...string) {
			s.Print("func %s (", name)
			for i := 0; i < len(parameters)-2; i += 2 {
				s.Print(parameters[i])
				s.Print(" ")
				s.Print(parameters[i+1])
				s.Print(",")
			}
			for i := len(parameters) - 2; i < len(parameters); i += 2 {
				s.Print(parameters[i])
				s.Print(" ")
				s.Print(parameters[i+1])
			}
			s.Print(") ")
			for i := 0; i < len(returnTypes)-1; i++ {
				s.Print(returnTypes[i])
				s.Print(",")
			}
			for i := len(returnTypes) - 1; i < len(returnTypes); i++ {
				s.Print(returnTypes[i])
			}
			s.Print("{\n")
		}
	}
}
