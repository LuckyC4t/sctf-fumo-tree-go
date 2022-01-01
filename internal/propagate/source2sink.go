package propagate

import (
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/data"
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/walker/finder"
	"github.com/VKCOM/php-parser/pkg/visitor/traverser"
)

var viewed = make(map[string]struct{})

func Source2Sink(source data.PHPClassMethod, sinkName string, callStack data.CallStack) {
	callStack.Push(source)

	sinkFinder := finder.NewFuncCallFinder(sinkName)
	traverser.NewTraverser(sinkFinder).Traverse(source.Method)
	if sinkFinder.FindFuncCall {
		callStack.Show()
	}
	viewed[source.String()] = struct{}{}

	// 随着 call 边跟出去
	for _, callee := range data.CallGraph[source.String()] {
		if _, has := viewed[callee.String()]; !has {
			Source2Sink(callee, sinkName, callStack)
		}
	}

	// 当前函数跟踪完成
	callStack.Pop()
}
