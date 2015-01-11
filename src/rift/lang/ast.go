package lang

import (
	"fmt"
	"strings"
)

const DEBUG = false

type parseStack struct{
	source Node
	stack Stack
}

func (s *parseStack) Start(Type string) {
	pushed := &Node{Type: Type}
	if DEBUG {
		fmt.Printf("Starting: %s\n", valueAsString(pushed))
	}
	s.stack.Push(pushed)
}

func (s *parseStack) Emit(value interface{}) {
	var top *Node
	if s.stack.Len() > 0 {
		top = s.stack.Peek().(*Node)
	} else {
		if DEBUG {
			fmt.Printf("Emitting value to source: %s\n", valueAsString(value))
		}
		top = &s.source
	}
	top.Add(value)
}

func (s *parseStack) End() {
	popped := s.stack.Pop()
	if DEBUG {
		fmt.Printf("Ending: %s\n", valueAsString(popped))
	}
	s.Emit(popped)
}

func valueAsString(value interface{}) string {
	switch v := value.(type) {
	default:
		return fmt.Sprintf("%+v", value)
	case *Node:
		return fmt.Sprintf("(%s %s)", v.Type, valueAsString(v.Values))
	// TODO: This sucks
	case []*Node:
		var nodeValues []string
		for _, nodeValue := range v {
			nodeValues = append(nodeValues, valueAsString(nodeValue))
		}
		return strings.Join(nodeValues, " ")
	case []interface{}:
		var nodeValues []string
		for _, nodeValue := range v {
			nodeValues = append(nodeValues, valueAsString(nodeValue))
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
	return fmt.Sprintf(valueAsString(s.Rifts()))
}

type Node struct{
	Type   string
	Values []interface{}
}

func (n *Node) Add(value interface{}) {
	n.Values = append(n.Values, value)
}
