package main

import (
	. "ACG/parser"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// 递归拷贝文件夹 src 到 dst
func CopyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		} else {
			// 拷贝文件
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			dstFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
			if err != nil {
				return err
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			return err
		}
	})
}

func GenerateSourceCode(bnf string, outputDir string, copyFrom string) error {
	// 拷贝到目标文件夹
	if copyFrom != "" {
		err := CopyDir(copyFrom, outputDir)
		if err != nil {
			fmt.Println("err:", err)
			return err
		}
	}

	path := bnf
	parser := NewBNFParser()
	rules, _ := parser.ParseFile(path)

	var generator Generator

	// parser
	parserFile := filepath.Join(outputDir, "parser/nika_parser.go")
	generator.SetOutputPath(parserFile)
	generator.GenerateLParser("NikaParser", rules)

	// compiler
	compilerFile := filepath.Join(outputDir, "parser/nika_compiler.go")
	generator.SetOutputPath(compilerFile)
	generator.PrintCompiler("NikaCompiler", rules)

	// nodes
	nodesFile := filepath.Join(outputDir, "parser/nika_nodes.go")
	generator.SetOutputPath(nodesFile)
	generator.PrintNodes("NikaCompiler", rules)

	return nil
}

// 编译 Go 文件为 exe
func BuildExecutable(srcDir, exePath string) error {
	cmd := exec.Command("go", "build", "-o", exePath, "main_compiler.go")
	cmd.Dir = srcDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请将 bnf 文件拖到此程序上运行")
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取 exe 路径失败:", err)
		fmt.Scanln()
		return
	}
	exeDir := filepath.Dir(exePath)
	fmt.Println("exeDir:", exeDir)
	bnfPath := os.Args[1]

	copyFrom := filepath.Join(exeDir, "src")

	// 创建临时源码目录
	tempSrcDir := filepath.Join(exeDir, "temp_nika")
	os.MkdirAll(tempSrcDir, 0755)

	fmt.Println("tempSrcDir", tempSrcDir)
	err = GenerateSourceCode(bnfPath, tempSrcDir, copyFrom)
	if err != nil {
		fmt.Println("生成源码失败:", err)
		fmt.Scanln()
		return
	}

	// 输出 exe
	outputExe := filepath.Join(exeDir, "output", "nika.exe")
	os.MkdirAll(filepath.Dir(outputExe), 0755)

	err = BuildExecutable(tempSrcDir, outputExe)
	if err != nil {
		fmt.Println("编译失败:", err)
		fmt.Scanln()
		return
	}

	fmt.Println("生成完成:", outputExe)

	// 可选：删除临时源码
	//os.RemoveAll(tempSrcDir)

	fmt.Scanln()
}
