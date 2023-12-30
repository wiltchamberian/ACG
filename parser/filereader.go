package parser

import (
	"fmt"
	"os"
)

type FileReader struct {
	path string
	file *os.File
	tabs string
}

func (s *FileReader) Read(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("can't read file", err)
		return nil
	}
	return content
}
