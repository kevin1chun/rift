package lang

import (
	"strings"
	"rift/support/sanity"
)

const (
	RIFT  = "rift"
	FUNC = "function-definition"
	FUNCAPPLY = "function-apply"
	ARGS = "arguments"
	TUPLE = "tuple"
	LIST = "list"
	ASSIGNMENT = "assignment"
	IF = "if"
	STRING = "string"
	NUM = "numeric"
	BOOL = "boolean"
	REF = "reference"
	OP = "operation"
	BINOP = "binary-operator"
)

type Source struct{
	rifts []*Node
}

type Node struct{
	Type   string
	Values []interface{}
}

func (n *Node) Add(value interface{}) {
	n.Values = append(n.Values, value)
}

func (n *Node) Rift() *Rift {
	sanity.Ensure(n.Type == RIFT, "Node must be [%s], but was [%s]", RIFT, n.Type)
	return &Rift{n}
}

func (n *Node) Ref() *Ref {
	sanity.Ensure(n.Type == RIFT, "Node must be [%s], but was [%s]", REF, n.Type)
	return &Ref{n}
}

func (n *Node) Assignment() *Assignment {
	sanity.Ensure(n.Type == ASSIGNMENT, "Node must be [%s], but was [%s]", ASSIGNMENT, n.Type)
	return &Assignment{n}
}

type Rift struct{
	node *Node
}

func (r *Rift) RawName() string {
	return r.node.Values[0].(*Node).Values[0].(string)
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
			assignments = append(assignments, line.Assignment())
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

func (r *Ref) IsLocal() bool {
	return len(r.node.Values) == 1
}

func (r *Ref) Rift() string {
	if r.IsLocal() {
		return "_"
	} else {
		return r.node.Values[0].(*Node).String()
	}
}

func (r *Ref) RawName() string {
	if r.IsLocal() {
		return r.node.Values[0].(*Node).String()
	} else {
		return r.node.Values[1].(*Node).String()
	}
}

func (r *Ref) Name() string {
	rawName := r.RawName()
	if r.HasGravity() {
		return rawName[1:]
	} else {
		return rawName
	}
}

func (r *Ref) HasGravity() bool {
	scoping := r.node.Values[0].(*Node).String()
	if r.IsLocal() {
		return strings.HasPrefix(scoping, "@")
	} else {
		return strings.HasPrefix(scoping, "@")
	}
}

func (r *Ref) String() string {
	var nameParts []string
	for _, value := range r.node.Values {
		nameParts = append(nameParts, value.(*Node).String())
	}
	return strings.Join(nameParts, ":")
}

type FuncApply struct{
	node *Node
}

func NewFuncApply(funcApply interface{}) *FuncApply {
	return &FuncApply{funcApply.(*Node)}
}

func (fa *FuncApply) Ref() *Ref {
	return fa.node.Values[0].(*Node).Ref()
}

func (fa *FuncApply) Args() []*Node {
	var values []*Node
	for _, value := range fa.node.Values[1].(*Node).Values {
		values = append(values, value.(*Node))
	}
	return values
}

type Assignment struct{
	node *Node
}

func (a *Assignment) Ref() *Ref {
	return a.node.Values[0].(*Node).Ref()
}

func (a *Assignment) Value() *Node {
	return a.node.Values[1].(*Node)
}

