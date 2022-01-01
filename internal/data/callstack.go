package data

import "fmt"

type CallStack struct {
	methodStack []PHPClassMethod
}

func (cs *CallStack) Pop() PHPClassMethod {
	topItem := cs.methodStack[len(cs.methodStack)-1]
	cs.methodStack = cs.methodStack[:len(cs.methodStack)-1]
	return topItem
}

func (cs *CallStack) Push(method PHPClassMethod) {
	cs.methodStack = append(cs.methodStack, method)
}

func (cs CallStack) Show() {
	for _, s := range cs.methodStack {
		fmt.Println(s)
	}
	fmt.Println("-----------------------------")
}
