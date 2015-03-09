package lang

import (
	"fmt"
	"io"
	"io/ioutil"
	"rift/support/collections"
	"strings"
)

func Parse(source io.Reader) (*riftParser, error) {
	readSource, sourceErr := ioutil.ReadAll(source)
	if sourceErr != nil {
		return nil, sourceErr
	}

	parser := &riftParser{Buffer: string(readSource[:])}
	parser.Init()
	err := parser.Parse()
	if err != nil {
		return parser, err
	}
	parser.Execute()

	return parser, nil
}

// TODO: Can this work any better?
func GetSyntaxErrors(p *riftParser) string {
	var errors []string
	for _, err := range p.Error() {
		pos := translatePositions(p.Buffer, []int{int(err.begin), int(err.end)})[0]
		errors = append(errors, fmt.Sprintf("Line %d, character %d", pos.line, pos.symbol))
	}
	return strings.Join(errors, "\n")
}

type parseStack struct{
	source Node
	stack collections.Stack
}

func (s *parseStack) Start(Type string) {
	s.stack.Push(&Node{Type: Type})
}

func (s *parseStack) Emit(value string) {
	var top *Node
	if s.stack.Len() > 0 {
		top = s.stack.Peek().(*Node)
	} else {
		top = &s.source
	}
	top.Add(value)
}

func (s *parseStack) EmitNode(value *Node) {
	var top *Node
	if s.stack.Len() > 0 {
		top = s.stack.Peek().(*Node)
	} else {
		top = &s.source
	}
	top.Add(value)
}

func (s *parseStack) End() {
	popped := s.stack.Pop()
	s.EmitNode(popped.(*Node))
}

func (s *parseStack) String() {
	ToString(s.source)
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	default:
		return fmt.Sprintf("%s", value)
	case *Node:
		return fmt.Sprintf("Node{Type:%s Values:%s}", v.Type, ToString(v.Values))
	// TODO: This sucks
	case []*Node:
		var nodeValues []string
		for _, nodeValue := range v {
			nodeValues = append(nodeValues, ToString(nodeValue))
		}
		return "[" + strings.Join(nodeValues, ", ") + "]"
	case []interface{}:
		var nodeValues []string
		for _, nodeValue := range v {
			nodeValues = append(nodeValues, ToString(nodeValue))
		}
		return "[" + strings.Join(nodeValues, ", ") + "]"
	}
}

func ToLisp(value interface{}) string {
	switch v := value.(type) {
	default:
		return fmt.Sprintf("%+v", value)
	case *Node:
		return fmt.Sprintf("(%s %s)", v.Type, ToLisp(v.Values))
	// TODO: This sucks
	case []*Node:
		var nodeValues []string
		for _, nodeValue := range v {
			nodeValues = append(nodeValues, ToLisp(nodeValue))
		}
		return strings.Join(nodeValues, " ")
	case []interface{}:
		var nodeValues []string
		for _, nodeValue := range v {
			nodeValues = append(nodeValues, ToLisp(nodeValue))
		}
		return strings.Join(nodeValues, " ")
	}
}

func (s *parseStack) Rifts() []*Node {
	var rifts []*Node
	for _, rift := range s.source.Values {
		rifts = append(rifts, rift.(*Node))
	}
	return rifts
}

func (s *parseStack) Lisp() string {
	return fmt.Sprintf(ToLisp(s))
}