package lang

import (
	"fmt"
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
	sanity.Ensure(n.Type == REF, "Node must be [%s], but was [%s]", REF, n.Type)
	return &Ref{n}
}

func (n *Node) Assignment() *Assignment {
	sanity.Ensure(n.Type == ASSIGNMENT, "Node must be [%s], but was [%s]", ASSIGNMENT, n.Type)
	return &Assignment{n}
}

func (n *Node) FuncApply() *FuncApply {
	sanity.Ensure(n.Type == FUNCAPPLY, "Node must be [%s], but was [%s]", FUNCAPPLY, n.Type)
	return &FuncApply{n}
}

func (n *Node) Func() *Func {
	sanity.Ensure(n.Type == FUNC, "Node must be [%s], but was [%s]", FUNC, n.Type)
	return &Func{n}
}

func (n *Node) Operation() *Operation {
	sanity.Ensure(n.Type == OP, "Node must be [%s], but was [%s]", OP, n.Type)
	return &Operation{n}
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

func (r *Rift) HasGravity() bool {
	return strings.HasPrefix(r.RawName(), "@")
}

func (r *Rift) String() string {
	return ToLisp(r.node)
}

type Func struct{
	node *Node
}

func (f *Func) Args() []*Ref {
	var args []*Ref
	for _, arg := range f.node.Values[0].(*Node).Values {
		args = append(args, arg.(*Node).Ref())
	}
	return args
}

func (f *Func) Lines() []*Node {
	var lines []*Node
	for _, line := range f.node.Values[1:] {
		lines = append(lines, line.(*Node))
	}
	return lines
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
		return r.node.Values[0].(*Node).Str()
	}
}

func (r *Ref) RawName() string {
	if r.IsLocal() {
		return r.node.Values[0].(*Node).Str()
	} else {
		return r.node.Values[1].(*Node).Str()
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
	scoping := r.node.Values[0].(*Node).Str()
	if r.IsLocal() {
		return strings.HasPrefix(scoping, "@")
	} else {
		return strings.HasPrefix(scoping, "@")
	}
}

func (r *Ref) String() string {
	var nameParts []string
	for _, value := range r.node.Values {
		nameParts = append(nameParts, value.(string))
	}
	return strings.Join(nameParts, ":")
}

type FuncApply struct{
	node *Node
}

func (fa *FuncApply) Ref() *Ref {
	return &Ref{fa.node.Values[0].(*Node)}
}

func (fa *FuncApply) Args() *Tuple {
	return &Tuple{fa.node.Values[1].(*Node)}
}

type Tuple struct{
	node *Node
}

func (t *Tuple) Arity() int {
	return len(t.node.Values)
}

func (t *Tuple) Values() []interface{} {
	return t.node.Values
}

func (t *Tuple) String() string {
	var values []string
	for _, value := range t.Values() {
		values = append(values, fmt.Sprintf("%s", value))
	}
	return "(" + strings.Join(values, ", ") + ")"
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


type Operation struct{
	node *Node
}

func (o *Operation) Operator() string {
	return o.node.Values[1].(*Node).Values[0].(string)
}

func (o *Operation) LHS() *Node {
	return o.node.Values[0].(*Node)
}

func (o *Operation) RHS() *Node {
	return o.node.Values[2].(*Node)
}
