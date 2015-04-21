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
		if rift.IsMain() {
			return rift
		}
	}

	return nil
}

func dereference(rift *lang.Rift, env collections.PersistentMap, ref *lang.Ref) interface{} {
	// TODO: Support gravity
	sanity.Ensure(env.Contains(ref.String()), "Undefined reference to [%s]", ref.String())
	return env.GetOrNil(ref.String())
}

func doAssignment(rift *lang.Rift, env collections.PersistentMap, assignment *lang.Assignment) interface{} {
	// TODO: Should I use lazy assignment here?
	var name string
	if rift.Name() == "main"{
		name = assignment.Ref().String()
	} else {
		name = rift.Name() + ":" + assignment.Ref().String()
	}
	env.Set(name, evaluate(rift, env, assignment.Value()))
	return nil
}

func doOperation(rift *lang.Rift, env collections.PersistentMap, op *lang.Operation) interface{} {
	lhsValue := evaluate(rift, env, op.LHS())
	rhsValue := evaluate(rift, env, op.RHS())
	// TODO: Handle boolean logic elsewhere
	// TODO: What to do about operator overloading
	return doMath(lhsValue, rhsValue, op.Operator())
}

func doIf(rift *lang.Rift, env collections.PersistentMap, i *lang.If) interface{} {
	cond := evaluate(rift, env, i.Condition()).(bool)
	var lastValue interface{}
	if cond {
		for _, line := range i.Lines() {
			lastValue = evaluate(rift, env, line)
		}
	} else {
		for _, line := range i.ElseLines() {
			lastValue = evaluate(rift, env, line)
		}
	}
	return lastValue
}

func evaluate(rift *lang.Rift, env collections.PersistentMap, v interface{}) interface{} {
	if a, isNode := v.(*lang.Node); isNode {
		switch a.Type {
		default:
			return nil
		case lang.IF:
			return doIf(rift, env, a.If())
		case lang.OP:
			return doOperation(rift, env, a.Operation())
		case lang.ASSIGNMENT:
			return doAssignment(rift, env, a.Assignment())
		case lang.FUNCAPPLY:
			return doFuncApply(rift, env, a.FuncApply())
		case lang.REF:
			return dereference(rift, env, a.Ref())
		case lang.FUNC:
			return makeFunc(rift, env, a.Func())
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

func evalRift(rift *lang.Rift, env collections.PersistentMap) {
	logging.Debug("Evaluating rift [%s]", rift.Name())
	for _, line := range rift.Lines() {
		evaluate(rift, env, line)
	}
}

func Run(rifts []*lang.Node) {
	InitPredefs()
	env := collections.ExtendPersistentMap(Predefs)
	for _, riftNode := range rifts {
		rift := riftNode.Rift()
		if !rift.IsMain() {
			evalRift(rift, env)
		}
	}
	if main := mainRift(rifts); main != nil {
		evalRift(main, env)
	} else {
		// TODO: Serve functionality
	}
	logging.Debug("Final environment:")
	for k, v := range env.Freeze() {
		logging.Debug(" |- %s = %+v", k, v)
	}
}
