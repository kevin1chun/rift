package runtime

import (
	"rift/lang"
	"rift/support/collections"
	"rift/support/logging"
)

func mainRift(riftDefs []*lang.Node) *lang.Rift {
	for _, riftDef := range riftDefs {
		rift := riftDef.Rift()
		if rift.Name() == "main" {
			return rift
		}
	}

	return nil
}

func doAssignment(env collections.PersistentMap, assignment *lang.Assignment) {
	var rhs interface{}
	switch assignment.Value().Type {
	default:
		rhs = assignment.Value()
	case lang.STRING:
		rhs = assignment.Value().Str()
	case lang.NUM:
		rhs = assignment.Value().Int()
	case lang.BOOL:
		rhs = assignment.Value().Bool()
	}

	env.Set(assignment.Ref().String(), rhs)
}

func Run(rifts []*lang.Node) {
	// TODO: Oops this only supports one rift :)
	if main := mainRift(rifts); main != nil {
		env := collections.NewPersistentMap()
		for _, line := range main.Lines() {
			switch line.Type{
			case lang.ASSIGNMENT:
				assignment := line.Assignment()
				logging.Debug("Assigning to ref [%s] the value [%s]", assignment.Ref(), assignment.Value())
				doAssignment(env, assignment)
			case lang.FUNCAPPLY:
				funcApply := line.FuncApply()
				logging.Debug("Apply func [%s] with args [%s]", funcApply.Ref(), funcApply.Args())
			}
		}
		logging.Debug("Final environment: %+v", env.Freeze())
	} else {
		logging.Warn("No such rift [main]")
	}
}