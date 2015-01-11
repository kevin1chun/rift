package vm

import (
	"fmt"
	"rift/lang"
	"rift/std"
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
				if ref.Rift() == "std" && ref.Name() == "println" {
					var interfaceArgs []interface{}
					for _, arg := range args {
						if arg.Type == lang.STRING {
							interfaceArgs = append(interfaceArgs, arg.Values[0].(string))
						} else {
							interfaceArgs = append(interfaceArgs, fmt.Sprintf("%+v", arg))
						}
					}
					std.Println(interfaceArgs...)
				} else {
					fmt.Printf("Applying func [%s] with args [%s]\n", ref, args)
				}
			}
		}
	}
}