package lang

import (
	"rift/support/logging"
)

func mainRift(riftDefs []*Node) *Rift {
	for _, riftDef := range riftDefs {
		rift := riftDef.Rift()
		if rift.Name() == "main" {
			return rift
		}
	}

	return nil
}

func Compile(rifts []*Node) {
	if main := mainRift(rifts); main != nil {
		for _, line := range main.Lines() {
			switch line.Type{
			case ASSIGNMENT:
				assignment := line.Assignment()
				logging.Info("Assigning to ref [%s] the value [%s]", assignment.Ref(), assignment.Value())
			case FUNCAPPLY:
				funcApply := line.FuncApply()
				logging.Info("Apply func [%s] with args [%s]", funcApply.Ref(), funcApply.Args())
			}
		}
	} else {
		logging.Warn("No such rift [main]")
	}
}