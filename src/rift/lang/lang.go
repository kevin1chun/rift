package lang

import (
	"strings"
)

const (
	RIFT  = "rift"
	FUNC = "function-definition"
	FUNCAPPLY = "function-apply"
	ARGS = "arguments"
	TUPLE = "tuple"
	LIST = "list"
	ASSIGNMENT = "assignment"
	STRING = "string"
	NUM = "numeric"
	BOOL = "boolean"
	REF = "reference"
	OP = "operation"
	BINOP = "binary-operator"
)

// TODO: Add node type assertions

type Rift struct{
	node *Node
	context map[string]interface{}
}

func NewRift(node interface{}) *Rift {
	return &Rift{node.(*Node), make(map[string]interface{})}
}

func (r *Rift) RawName() string {
	return r.node.Values[0].(string)
}

func (r *Rift) Name() string {
	rawName := r.RawName()
	if r.HasGravity() {
		return string(rawName[1:])
	} else {
		return rawName
	}
}

func (r *Rift) Lines() []*Node {
	var lines []*Node
	for _, line := range r.node.Values[1:] {
		lines = append(lines, line.(*Node))
	}
	return lines
}

func (r *Rift) Assignments() []*Assignment {
	var assignments []*Assignment
	for _, line := range r.Lines() {
		if line.Type == ASSIGNMENT {
			assignments = append(assignments, NewAssignment(line))
		}
	}
	return assignments
}

// TODO: Should this somehow be separate from the main code?
func (r *Rift) Protocol() map[string]*Node {
	proto := make(map[string]*Node)
	for _, assignment := range r.Assignments() {
		value := assignment.Value()
		if value.Type == FUNC {
			proto[assignment.Ref().Name()] = assignment.Value()
		}
	}
	return proto
}

func (r *Rift) HasGravity() bool {
	return strings.HasPrefix(r.RawName(), "@")
}

type Ref struct{
	node *Node
}

func NewRef(ref interface{}) *Ref {
	return &Ref{ref.(*Node)}
}

func (r *Ref) IsLocal() bool {
	return len(r.node.Values) == 1
}

func (r *Ref) Rift() string {
	if r.IsLocal() {
		return "@"
	} else {
		return r.node.Values[0].(string)
	}
}

func (r *Ref) Name() string {
	if r.IsLocal() {
		return r.node.Values[0].(string)
	} else {
		return r.node.Values[1].(string)
	}
}

func (r *Ref) String() string {
	return r.Rift() + ":" + r.Name()
}

type FuncApply struct{
	node *Node
}

func NewFuncApply(funcApply interface{}) *FuncApply {
	return &FuncApply{funcApply.(*Node)}
}

func (fa *FuncApply) Ref() *Ref {
	return NewRef(fa.node.Values[0])
}

func (fa *FuncApply) Args() []*Node {
	var args []*Node
	for _, arg := range fa.node.Values[1].(*Node).Values {
		args = append(args, arg.(*Node))
	}
	return args
}

type Assignment struct{
	node *Node
}

func NewAssignment(assignment interface{}) *Assignment {
	return &Assignment{assignment.(*Node)}
}

func (a *Assignment) Ref() *Ref {
	return NewRef(a.node.Values[0])
}

func (a *Assignment) Value() *Node {
	return a.node.Values[1].(*Node)
}
