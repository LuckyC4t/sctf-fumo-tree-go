package main

import (
	"fmt"
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/data"
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/parser"
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/propagate"
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/walker"
	"github.com/VKCOM/php-parser/pkg/visitor/traverser"
	"log"
	"os"
)

func main() {
	classFilePath := os.Args[1]
	fileInfo, err := os.Stat(classFilePath)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(fmt.Errorf("%s is not exist", classFilePath))
		}
	}
	if fileInfo.IsDir() {
		log.Fatal(fmt.Errorf("class file %s is not file", classFilePath))
	}

	// 解析 class 文件
	rootNode, errs := parser.ParsePHP(classFilePath)
	if len(errs) != 0 {
		for _, e := range errs {
			log.Println(e)
		}
		os.Exit(-1)
	}

	// 提取类信息并进行索引
	classesInfo := walker.NewClassInfoWalker()
	traverser.NewTraverser(classesInfo).Traverse(rootNode)
	// 索引类信息
	for _, cls := range classesInfo.PHPClasses {
		data.ClassIndex[cls.Name] = cls
		for _, m := range cls.Methods {
			data.ClassMethodIndex[m.Name] =
				append(data.ClassMethodIndex[m.Name], m)
		}
	}

	// 构建 call graph
	for _, methods := range data.ClassMethodIndex {
		for _, method := range methods {
			key := method.String()
			callGraph := walker.NewCallGraphWalker()
			traverser.NewTraverser(callGraph).Traverse(method.Method)
			data.CallGraph[key] = append(data.CallGraph[key], callGraph.CallOut...)
		}
	}

	// 查找路径
	sinkFuncName := "readfile"
	entry := data.ClassMethodIndex["__destruct"]
	var callStack data.CallStack
	propagate.Source2Sink(entry[0], sinkFuncName, callStack)
}
