package vm

import (
	"fmt"
	"rift/lang"
)

type VM struct{
	rifts map[string]*lang.Rift
}

func New(nodes []*lang.Node) *VM {
	rifts := map[string]*lang.Rift{
		"std": nil, // TODO: Prepopulate built-ins
	}
	for _, node := range nodes {
		rift := lang.NewRift(node)
		rifts[rift.Name()] = rift
	}

	for k, _ := range rifts {
		fmt.Printf("Discovered rift[%s]\n", k)
	}
	return &VM{rifts}
}

func (vm *VM) Run() {
	main, mainExists := vm.rifts["main"]
	if mainExists {
		for _, line := range main.Lines() {
			switch line.Type {
			case lang.FUNCAPPLY:
				funcApply := lang.NewFuncApply(line)
				ref := funcApply.Ref()
				args := funcApply.Args()
				fmt.Printf("Calling func [%s] with args [%s]\n", ref, args)
			}
		}
	}
}