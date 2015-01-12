package runtime

import (
	"rift/lang"
)

type VM struct{
	dispatcher Dispatcher
}

func NewVM(dispatcher Dispatcher) *VM {
	return &VM{dispatcher}
}

func (vm *VM) Run() {
	entryPoint, entryPointExists := vm.dispatcher.EntryPoint()
	if entryPointExists {
		for _, line := range entryPoint.Lines() {
			switch line.Type {
			case lang.FUNCAPPLY:
				funcApply := lang.NewFuncApply(line)
				vm.dispatcher.Dispatch(funcApply)
			}
		}
	}
}

