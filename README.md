# ACG
Automatically code generation, 
自动生成parser和compiler.

## user guide
1. 首先参考nika_vm.gram定义一个BNF文件,假设叫abc.gram。
2. 在main.go所在目录执行 go build -o main.exe main.go 生成main.exe
3. 在项目外建立新文件夹(假设叫Test)，把main.exe拖进去。
4. 在Test里建立src文件夹。把parser文件夹整个拷贝到src里
5. 如果parser文件夹有nika_打头的文件删除，再把main_compiler.go拷贝到src文件夹下
6. 把abc.gram拖到main.exe里，等待生成一个output文件夹，里面有nika.exe，这个就是生成的编译器（支持虚拟机运行）,同时还会生成temp_nika是这个编译器的源码。
7. 把目标文件比如.c拖到这个编译器里就正常执行了，也可也修改源码执行。

## future work
1. 支持更多功能。