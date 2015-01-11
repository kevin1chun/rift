package lang

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

func (r *Rift) Name() string {
	return r.node.Values[0].(string)
}

func (r *Rift) Lines() []*Node {
	var lines []*Node
	for _, line := range r.node.Values[1:] {
		lines = append(lines, line.(*Node))
	}
	return lines
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
	for _, arg := range fa.node.Values[1:] {
		args = append(args, arg.(*Node))
	}
	return args
}
