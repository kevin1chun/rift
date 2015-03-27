package runtime

import (
	"rift/lang"
	"rift/support/collections"
	"rift/support/logging"
	"rift/support/sanity"
)

// TODO: Better organization
// TODO: Consistency in when things are evaluated
// TODO: Remote dereference/dispatch
// TODO: Gravitasse for remote function dispatch
// TODO: Is `nil` okay for void ops?
// TODO: Support multiple rifts in references
// TODO: Tail-call optimization

func mainRift(riftDefs []*lang.Node) *lang.Rift {
	for _, riftDef := range riftDefs {
		rift := riftDef.Rift()
		if rift.Name() == "main" {
			return rift
		}
	}

	return nil
}

func dereference(env collections.PersistentMap, ref *lang.Ref) interface{} {
	// TODO: Support multi-rift scenario
	sanity.Ensure(env.Contains(ref.String()), "Undefined reference to [%s]", ref.String())
	return env.GetOrNil(ref.String())
}

func doAssignment(env collections.PersistentMap, assignment *lang.Assignment) interface{} {
	// TODO: Should I use lazy assignment here?
	env.Set(assignment.Ref().String(), evaluate(env, assignment.Value()))
	return nil
}

func doOperation(env collections.PersistentMap, op *lang.Operation) interface{} {
	lhsValue := evaluate(env, op.LHS())
	rhsValue := evaluate(env, op.RHS())
	// TODO: Handle boolean logic elsewhere
	// TODO: What to do about operator overloading
	return doMath(lhsValue, rhsValue, op.Operator())
}

func doIf(env collections.PersistentMap, i *lang.If) interface{} {
	cond := evaluate(env, i.Condition()).(bool)
	var lastValue interface{}
	if cond {
		for _, line := range i.Lines() {
			lastValue = evaluate(env, line)
		}
	} else {
		for _, line := range i.ElseLines() {
			lastValue = evaluate(env, line)
		}
	}
	return lastValue
}

func evaluate(env collections.PersistentMap, v interface{}) interface{} {
	if a, isNode := v.(*lang.Node); isNode {
		switch a.Type {
		default:
			return nil
		case lang.IF:
			return doIf(env, a.If())
		case lang.OP:
			return doOperation(env, a.Operation())
		case lang.ASSIGNMENT:
			return doAssignment(env, a.Assignment())
		case lang.FUNCAPPLY:
			return doFuncApply(env, a.FuncApply())
		case lang.REF:
			return dereference(env, a.Ref())
		case lang.FUNC:
			return makeFunc(env, a.Func())
		case lang.STRING:
			return a.Str()
		case lang.NUM:
			return a.Int()
		case lang.BOOL:
			return a.Bool()
		}
	} else {
		return v
	}
}

func Run(rifts []*lang.Node) {
	// TODO: Oops this only supports one rift :)
	if main := mainRift(rifts); main != nil {
		// ctx := collections.Stack{}
		InitPredefs()
		env := collections.ExtendPersistentMap(Predefs)
		for _, line := range main.Lines() {
			evaluate(env, line)
		}
		logging.Debug("Final environment:")
		for k, v := range env.Freeze() {
			logging.Debug(" |- %s = %+v", k, v)
		}
	} else {
		logging.Warn("No such rift [main]")
	}
}
