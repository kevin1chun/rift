package runtime

import (
	"fmt"
	"rift/lang"
	"rift/std"
)

type Dispatcher interface{
	EntryPoint() (*lang.Rift, bool)
	Dispatch(*lang.FuncApply) interface{}
}

type LocalDispatcher struct{
	rifts map[string]*lang.Rift
}

func NewLocalDispatcher(nodes []*lang.Node) *LocalDispatcher {
	rifts := map[string]*lang.Rift{
		// TODO: Prepopulate built-ins
	}
	for _, node := range nodes {
		rift := lang.NewRift(node)
		rifts[rift.Name()] = rift
	}
	return &LocalDispatcher{rifts}
}

func (d *LocalDispatcher) EntryPoint() (*lang.Rift, bool) {
	main, mainExists := d.rifts["main"]
	return main, mainExists
}

func (d *LocalDispatcher) Dispatch(funcApply *lang.FuncApply) interface{} {
	ref := funcApply.Ref()
	args := funcApply.Args()
	var f func(*lang.Ref, ...*lang.Node) interface{}
	var funcExists bool
	// TODO: Create a built-in for the `std` rift
	if ref.Rift() == "std" {
		f = StdApplier
		funcExists = true
	} else {
		rift, riftExists := d.rifts[ref.Rift()]
		if riftExists {
			protocol := rift.Protocol()
			_, funcExists = protocol[ref.Name()]
			funcExists = false
		}
	}

	if funcExists {
		return f(ref, args...)
	} else {
		fmt.Printf("No such function [%s]\n", ref)
	}
	return nil
}

func StdApplier(ref *lang.Ref, args...*lang.Node) interface{} {
	switch ref.Name() {
	default:
		fmt.Printf("No such function [%s:%s] exists\n", ref.Rift(), ref.Name())
		return nil
	case "println":
		var interfaceArgs []interface{}
		for _, arg := range args {
			if arg.Type == lang.STRING {
				interfaceArgs = append(interfaceArgs, arg.Values[0].(string))
			} else {
				interfaceArgs = append(interfaceArgs, fmt.Sprintf("%+v", arg))
			}
		}
		std.Println(interfaceArgs...)
		return nil
	}
}
