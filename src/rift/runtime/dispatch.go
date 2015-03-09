package runtime

import (
	"fmt"
	"rift/lang"
)

type Dispatcher interface{
	Dispatch(*Context, *lang.FuncApply) interface{}
}

type LocalDispatcher struct{
}

func (d *LocalDispatcher) Dispatch(context *Context, funcApply *lang.FuncApply) interface{} {
	ref := funcApply.Ref()
	// args := funcApply.Args()
	funcExists := context.Exists(ref.String())
	if funcExists {
		fmt.Printf("Found function [%s]\n", ref)
	} else {
		fmt.Printf("No such function [%s]\n", ref)
	}
	return nil
}
