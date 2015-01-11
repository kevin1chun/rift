package lang

import (
	"fmt"
	"strings"
)

const DEBUG = false

type parseStack struct{
	source astNode
	stack Stack
}

func (s *parseStack) Start(nodeType string) {
	pushed := &astNode{nodeType: nodeType}
	if DEBUG {
		fmt.Printf("Starting: %s\n", valueAsString(pushed))
	}
	s.stack.Push(pushed)
}

func (s *parseStack) Emit(value interface{}) {
	var top *astNode
	if s.stack.Len() > 0 {
		top = s.stack.Peek().(*astNode)
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
	case *astNode:
		return fmt.Sprintf("(%s %s)", v.nodeType, valueAsString(v.nodeValues))
	case []interface{}:
		var nodeValues []string
		for _, nodeValue := range v {
			nodeValues = append(nodeValues, valueAsString(nodeValue))
		}
		return strings.Join(nodeValues, " ")
	}
}

func (s *parseStack) Lisp() string {
	return fmt.Sprintf(valueAsString(s.source.nodeValues))
}

type astNode struct{
	nodeType   string
	nodeValues []interface{}
}

func (n *astNode) Add(value interface{}) {
	n.nodeValues = append(n.nodeValues, value)
}
