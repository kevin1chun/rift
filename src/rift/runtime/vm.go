package runtime

import (
	"fmt"
	"rift/support/collections"
)

type VM struct{
	contextStack collections.Stack
	dispatcher   Dispatcher
}

func NewVM(dispatcher Dispatcher) *VM {
	return &VM{collections.Stack{}, dispatcher}
}

func (vm *VM) SwitchTo(context *Context) {
	vm.contextStack.Push(context)
}

func (vm *VM) Return() {
	vm.contextStack.Pop()
}

func (vm *VM) Run() {
	initialCtx := vm.contextStack.Pop().(*Context)
	fmt.Printf("Initial context: %+v\n", initialCtx)
	if initialCtx.Exists("@:main") {
		main := initialCtx.Dereference("@:main")
		fmt.Printf("Main: %+v\n", main)
	}
	// entryPoint, entryPointExists := vm.dispatcher.EntryPoint()
	// if entryPointExists {
	// 	for _, line := range entryPoint.Lines() {
	// 		switch line.Type {
	// 		case lang.FUNCAPPLY:
	// 			funcApply := lang.NewFuncApply(line)
	// 			vm.dispatcher.Dispatch(funcApply)
	// 		}
	// 	}
	// }
}

