package lang

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 4

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleSource
	ruleRift
	ruleBlock
	ruleLine
	ruleSingle
	ruleExpr
	ruleOp
	ruleBinaryOp
	ruleStatement
	ruleAssignment
	ruleIf
	ruleElse
	ruleRef
	ruleFullRef
	ruleLocalRef
	ruleRefChar
	ruleValue
	ruleLiteral
	ruleScalar
	ruleVector
	ruleString
	ruleStringChar
	ruleStringEsc
	ruleSimpleEsc
	ruleNumeric
	ruleSciNum
	ruleDecimal
	ruleInteger
	ruleWholeNum
	ruleDigit
	ruleBoolean
	ruleFunc
	ruleFuncArgs
	ruleFuncApply
	ruleList
	ruleTuple
	ruleMap
	ruleGravitasse
	rulemsp
	rulesp
	rulecomment
	rulews
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	rulePegText
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35
	ruleAction36
	ruleAction37
	ruleAction38
	ruleAction39
	ruleAction40
	ruleAction41

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"Source",
	"Rift",
	"Block",
	"Line",
	"Single",
	"Expr",
	"Op",
	"BinaryOp",
	"Statement",
	"Assignment",
	"If",
	"Else",
	"Ref",
	"FullRef",
	"LocalRef",
	"RefChar",
	"Value",
	"Literal",
	"Scalar",
	"Vector",
	"String",
	"StringChar",
	"StringEsc",
	"SimpleEsc",
	"Numeric",
	"SciNum",
	"Decimal",
	"Integer",
	"WholeNum",
	"Digit",
	"Boolean",
	"Func",
	"FuncArgs",
	"FuncApply",
	"List",
	"Tuple",
	"Map",
	"Gravitasse",
	"msp",
	"sp",
	"comment",
	"ws",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"PegText",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",
	"Action36",
	"Action37",
	"Action38",
	"Action39",
	"Action40",
	"Action41",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(buffer[node.begin:node.end]))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token16 struct {
	pegRule
	begin, end, next int16
}

func (t *token16) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token16) isParentOf(u token16) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token16) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token16) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens16 struct {
	tree    []token16
	ordered [][]token16
}

func (t *tokens16) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens16) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens16) Order() [][]token16 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int16, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token16, len(depths)), make([]token16, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int16(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state16 struct {
	token16
	depths []int16
	leaf   bool
}

func (t *tokens16) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens16) PreOrder() (<-chan state16, [][]token16) {
	s, ordered := make(chan state16, 6), t.Order()
	go func() {
		var states [8]state16
		for i, _ := range states {
			states[i].depths = make([]int16, len(ordered))
		}
		depths, state, depth := make([]int16, len(ordered)), 0, 1
		write := func(t token16, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int16(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token16 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token16{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token16{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token16{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens16) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens16) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens16) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token16{pegRule: rule, begin: int16(begin), end: int16(end), next: int16(depth)}
}

func (t *tokens16) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens16) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next int32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token32{pegRule: rule, begin: int32(begin), end: int32(end), next: int32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type riftParser struct {
	parseStack

	Buffer string
	buffer []rune
	rules  [86]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer[0:] {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p *riftParser
}

func (e *parseError) Error() string {
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *riftParser) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *riftParser) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *riftParser) Execute() {
	buffer, begin, end := p.Buffer, 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {
		case rulePegText:
			begin, end = int(token.begin), int(token.end)
		case ruleAction0:
			p.Start(RIFT)
		case ruleAction1:
			p.End()
		case ruleAction2:
			p.Start(OP)
		case ruleAction3:
			p.End()
		case ruleAction4:
			p.Start(BINOP)
		case ruleAction5:
			p.Emit(string(buffer[begin:end]))
		case ruleAction6:
			p.End()
		case ruleAction7:
			p.Start(ASSIGNMENT)
		case ruleAction8:
			p.End()
		case ruleAction9:
			p.Start(IF)
		case ruleAction10:
			p.End()
		case ruleAction11:
			p.Start(ELSE)
		case ruleAction12:
			p.End()
		case ruleAction13:
			p.Start(REF)
		case ruleAction14:
			p.Emit(string(buffer[begin:end]))
		case ruleAction15:
			p.Emit(string(buffer[begin:end]))
		case ruleAction16:
			p.End()
		case ruleAction17:
			p.Start(REF)
		case ruleAction18:
			p.Emit(string(buffer[begin:end]))
		case ruleAction19:
			p.End()
		case ruleAction20:
			p.Start(STRING)
		case ruleAction21:
			p.Emit(string(buffer[begin:end]))
		case ruleAction22:
			p.End()
		case ruleAction23:
			p.Start(NUM)
		case ruleAction24:
			p.End()
		case ruleAction25:
			p.Emit(string(buffer[begin:end]))
		case ruleAction26:
			p.Emit(string(buffer[begin:end]))
		case ruleAction27:
			p.Start(BOOL)
		case ruleAction28:
			p.Emit(string(buffer[begin:end]))
		case ruleAction29:
			p.End()
		case ruleAction30:
			p.Start(FUNC)
		case ruleAction31:
			p.End()
		case ruleAction32:
			p.Start(ARGS)
		case ruleAction33:
			p.End()
		case ruleAction34:
			p.Start(FUNCAPPLY)
		case ruleAction35:
			p.End()
		case ruleAction36:
			p.Start(LIST)
		case ruleAction37:
			p.End()
		case ruleAction38:
			p.Start(TUPLE)
		case ruleAction39:
			p.End()
		case ruleAction40:
			p.Start("map")
		case ruleAction41:
			p.End()

		}
	}
}

func (p *riftParser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens16{tree: make([]token16, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := 0, 0, 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin int) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	matchDot := func() bool {
		if buffer[position] != end_symbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Source <- <(sp (Rift sp)+ !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[rulesp]() {
					goto l0
				}
				{
					position4 := position
					depth++
					{
						add(ruleAction0, position)
					}
					{
						position6, tokenIndex6, depth6 := position, tokenIndex, depth
						{
							position8 := position
							depth++
							if buffer[position] != rune('@') {
								goto l6
							}
							position++
							depth--
							add(ruleGravitasse, position8)
						}
						goto l7
					l6:
						position, tokenIndex, depth = position6, tokenIndex6, depth6
					}
				l7:
					if !_rules[ruleLocalRef]() {
						goto l0
					}
					if !_rules[rulesp]() {
						goto l0
					}
					if buffer[position] != rune('=') {
						goto l0
					}
					position++
					if buffer[position] != rune('>') {
						goto l0
					}
					position++
					if !_rules[rulesp]() {
						goto l0
					}
					if !_rules[ruleBlock]() {
						goto l0
					}
					{
						add(ruleAction1, position)
					}
					depth--
					add(ruleRift, position4)
				}
				if !_rules[rulesp]() {
					goto l0
				}
			l2:
				{
					position3, tokenIndex3, depth3 := position, tokenIndex, depth
					{
						position10 := position
						depth++
						{
							add(ruleAction0, position)
						}
						{
							position12, tokenIndex12, depth12 := position, tokenIndex, depth
							{
								position14 := position
								depth++
								if buffer[position] != rune('@') {
									goto l12
								}
								position++
								depth--
								add(ruleGravitasse, position14)
							}
							goto l13
						l12:
							position, tokenIndex, depth = position12, tokenIndex12, depth12
						}
					l13:
						if !_rules[ruleLocalRef]() {
							goto l3
						}
						if !_rules[rulesp]() {
							goto l3
						}
						if buffer[position] != rune('=') {
							goto l3
						}
						position++
						if buffer[position] != rune('>') {
							goto l3
						}
						position++
						if !_rules[rulesp]() {
							goto l3
						}
						if !_rules[ruleBlock]() {
							goto l3
						}
						{
							add(ruleAction1, position)
						}
						depth--
						add(ruleRift, position10)
					}
					if !_rules[rulesp]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				{
					position16, tokenIndex16, depth16 := position, tokenIndex, depth
					if !matchDot() {
						goto l16
					}
					goto l0
				l16:
					position, tokenIndex, depth = position16, tokenIndex16, depth16
				}
				depth--
				add(ruleSource, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Rift <- <(Action0 Gravitasse? LocalRef sp ('=' '>') sp Block Action1)> */
		nil,
		/* 2 Block <- <('{' sp (Line msp)* '}')> */
		func() bool {
			position18, tokenIndex18, depth18 := position, tokenIndex, depth
			{
				position19 := position
				depth++
				if buffer[position] != rune('{') {
					goto l18
				}
				position++
				if !_rules[rulesp]() {
					goto l18
				}
			l20:
				{
					position21, tokenIndex21, depth21 := position, tokenIndex, depth
					{
						position22 := position
						depth++
						{
							position23, tokenIndex23, depth23 := position, tokenIndex, depth
							{
								position25 := position
								depth++
								{
									position26, tokenIndex26, depth26 := position, tokenIndex, depth
									{
										position28 := position
										depth++
										{
											add(ruleAction7, position)
										}
										if !_rules[ruleLocalRef]() {
											goto l27
										}
										if !_rules[rulesp]() {
											goto l27
										}
										if buffer[position] != rune('=') {
											goto l27
										}
										position++
										if !_rules[rulesp]() {
											goto l27
										}
										if !_rules[ruleExpr]() {
											goto l27
										}
										{
											add(ruleAction8, position)
										}
										depth--
										add(ruleAssignment, position28)
									}
									goto l26
								l27:
									position, tokenIndex, depth = position26, tokenIndex26, depth26
									{
										position31 := position
										depth++
										{
											add(ruleAction9, position)
										}
										if buffer[position] != rune('i') {
											goto l24
										}
										position++
										if buffer[position] != rune('f') {
											goto l24
										}
										position++
										if !_rules[rulesp]() {
											goto l24
										}
										if !_rules[ruleExpr]() {
											goto l24
										}
										if !_rules[rulesp]() {
											goto l24
										}
										if !_rules[ruleBlock]() {
											goto l24
										}
										{
											position33, tokenIndex33, depth33 := position, tokenIndex, depth
											if !_rules[rulesp]() {
												goto l33
											}
											{
												position35 := position
												depth++
												{
													add(ruleAction11, position)
												}
												if buffer[position] != rune('e') {
													goto l33
												}
												position++
												if buffer[position] != rune('l') {
													goto l33
												}
												position++
												if buffer[position] != rune('s') {
													goto l33
												}
												position++
												if buffer[position] != rune('e') {
													goto l33
												}
												position++
												if !_rules[rulesp]() {
													goto l33
												}
												if !_rules[ruleBlock]() {
													goto l33
												}
												{
													add(ruleAction12, position)
												}
												depth--
												add(ruleElse, position35)
											}
											goto l34
										l33:
											position, tokenIndex, depth = position33, tokenIndex33, depth33
										}
									l34:
										{
											add(ruleAction10, position)
										}
										depth--
										add(ruleIf, position31)
									}
								}
							l26:
								depth--
								add(ruleStatement, position25)
							}
							goto l23
						l24:
							position, tokenIndex, depth = position23, tokenIndex23, depth23
							if !_rules[ruleExpr]() {
								goto l21
							}
						}
					l23:
						depth--
						add(ruleLine, position22)
					}
					{
						position39 := position
						depth++
						{
							position42, tokenIndex42, depth42 := position, tokenIndex, depth
							if !_rules[rulews]() {
								goto l43
							}
							goto l42
						l43:
							position, tokenIndex, depth = position42, tokenIndex42, depth42
							if !_rules[rulecomment]() {
								goto l21
							}
						}
					l42:
					l40:
						{
							position41, tokenIndex41, depth41 := position, tokenIndex, depth
							{
								position44, tokenIndex44, depth44 := position, tokenIndex, depth
								if !_rules[rulews]() {
									goto l45
								}
								goto l44
							l45:
								position, tokenIndex, depth = position44, tokenIndex44, depth44
								if !_rules[rulecomment]() {
									goto l41
								}
							}
						l44:
							goto l40
						l41:
							position, tokenIndex, depth = position41, tokenIndex41, depth41
						}
						depth--
						add(rulemsp, position39)
					}
					goto l20
				l21:
					position, tokenIndex, depth = position21, tokenIndex21, depth21
				}
				if buffer[position] != rune('}') {
					goto l18
				}
				position++
				depth--
				add(ruleBlock, position19)
			}
			return true
		l18:
			position, tokenIndex, depth = position18, tokenIndex18, depth18
			return false
		},
		/* 3 Line <- <(Statement / Expr)> */
		nil,
		/* 4 Single <- <(FuncApply / Value)> */
		func() bool {
			position47, tokenIndex47, depth47 := position, tokenIndex, depth
			{
				position48 := position
				depth++
				{
					position49, tokenIndex49, depth49 := position, tokenIndex, depth
					{
						position51 := position
						depth++
						{
							add(ruleAction34, position)
						}
						if !_rules[ruleRef]() {
							goto l50
						}
						if !_rules[ruleTuple]() {
							goto l50
						}
						{
							add(ruleAction35, position)
						}
						depth--
						add(ruleFuncApply, position51)
					}
					goto l49
				l50:
					position, tokenIndex, depth = position49, tokenIndex49, depth49
					{
						position54 := position
						depth++
						{
							position55, tokenIndex55, depth55 := position, tokenIndex, depth
							if !_rules[ruleRef]() {
								goto l56
							}
							goto l55
						l56:
							position, tokenIndex, depth = position55, tokenIndex55, depth55
							{
								position57 := position
								depth++
								{
									position58, tokenIndex58, depth58 := position, tokenIndex, depth
									{
										position60 := position
										depth++
										{
											add(ruleAction30, position)
										}
										{
											position62 := position
											depth++
											{
												add(ruleAction32, position)
											}
											if buffer[position] != rune('(') {
												goto l59
											}
											position++
											if !_rules[rulesp]() {
												goto l59
											}
											{
												position64, tokenIndex64, depth64 := position, tokenIndex, depth
												if !_rules[ruleLocalRef]() {
													goto l64
												}
											l66:
												{
													position67, tokenIndex67, depth67 := position, tokenIndex, depth
													if !_rules[rulesp]() {
														goto l67
													}
													if buffer[position] != rune(',') {
														goto l67
													}
													position++
													if !_rules[rulesp]() {
														goto l67
													}
													if !_rules[ruleLocalRef]() {
														goto l67
													}
													goto l66
												l67:
													position, tokenIndex, depth = position67, tokenIndex67, depth67
												}
												if !_rules[rulesp]() {
													goto l64
												}
												goto l65
											l64:
												position, tokenIndex, depth = position64, tokenIndex64, depth64
											}
										l65:
											if buffer[position] != rune(')') {
												goto l59
											}
											position++
											{
												add(ruleAction33, position)
											}
											depth--
											add(ruleFuncArgs, position62)
										}
										if !_rules[rulesp]() {
											goto l59
										}
										if buffer[position] != rune('-') {
											goto l59
										}
										position++
										if buffer[position] != rune('>') {
											goto l59
										}
										position++
										if !_rules[rulesp]() {
											goto l59
										}
										{
											position69, tokenIndex69, depth69 := position, tokenIndex, depth
											if !_rules[ruleBlock]() {
												goto l70
											}
											goto l69
										l70:
											position, tokenIndex, depth = position69, tokenIndex69, depth69
											if !_rules[ruleExpr]() {
												goto l59
											}
										}
									l69:
										{
											add(ruleAction31, position)
										}
										depth--
										add(ruleFunc, position60)
									}
									goto l58
								l59:
									position, tokenIndex, depth = position58, tokenIndex58, depth58
									{
										position73 := position
										depth++
										{
											switch buffer[position] {
											case 'f', 't':
												{
													position75 := position
													depth++
													{
														add(ruleAction27, position)
													}
													{
														position77 := position
														depth++
														{
															position78, tokenIndex78, depth78 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l79
															}
															position++
															if buffer[position] != rune('r') {
																goto l79
															}
															position++
															if buffer[position] != rune('u') {
																goto l79
															}
															position++
															if buffer[position] != rune('e') {
																goto l79
															}
															position++
															goto l78
														l79:
															position, tokenIndex, depth = position78, tokenIndex78, depth78
															if buffer[position] != rune('f') {
																goto l72
															}
															position++
															if buffer[position] != rune('a') {
																goto l72
															}
															position++
															if buffer[position] != rune('l') {
																goto l72
															}
															position++
															if buffer[position] != rune('s') {
																goto l72
															}
															position++
															if buffer[position] != rune('e') {
																goto l72
															}
															position++
														}
													l78:
														depth--
														add(rulePegText, position77)
													}
													{
														add(ruleAction28, position)
													}
													{
														add(ruleAction29, position)
													}
													depth--
													add(ruleBoolean, position75)
												}
												break
											case '"':
												{
													position82 := position
													depth++
													{
														add(ruleAction20, position)
													}
													if buffer[position] != rune('"') {
														goto l72
													}
													position++
													{
														position84 := position
														depth++
													l85:
														{
															position86, tokenIndex86, depth86 := position, tokenIndex, depth
															{
																position87 := position
																depth++
																{
																	position88, tokenIndex88, depth88 := position, tokenIndex, depth
																	{
																		position90 := position
																		depth++
																		{
																			position91 := position
																			depth++
																			if buffer[position] != rune('\\') {
																				goto l89
																			}
																			position++
																			{
																				switch buffer[position] {
																				case 'v':
																					if buffer[position] != rune('v') {
																						goto l89
																					}
																					position++
																					break
																				case 't':
																					if buffer[position] != rune('t') {
																						goto l89
																					}
																					position++
																					break
																				case 'r':
																					if buffer[position] != rune('r') {
																						goto l89
																					}
																					position++
																					break
																				case 'n':
																					if buffer[position] != rune('n') {
																						goto l89
																					}
																					position++
																					break
																				case 'f':
																					if buffer[position] != rune('f') {
																						goto l89
																					}
																					position++
																					break
																				case 'b':
																					if buffer[position] != rune('b') {
																						goto l89
																					}
																					position++
																					break
																				case 'a':
																					if buffer[position] != rune('a') {
																						goto l89
																					}
																					position++
																					break
																				case '\\':
																					if buffer[position] != rune('\\') {
																						goto l89
																					}
																					position++
																					break
																				case '?':
																					if buffer[position] != rune('?') {
																						goto l89
																					}
																					position++
																					break
																				case '"':
																					if buffer[position] != rune('"') {
																						goto l89
																					}
																					position++
																					break
																				default:
																					if buffer[position] != rune('\'') {
																						goto l89
																					}
																					position++
																					break
																				}
																			}

																			depth--
																			add(ruleSimpleEsc, position91)
																		}
																		depth--
																		add(ruleStringEsc, position90)
																	}
																	goto l88
																l89:
																	position, tokenIndex, depth = position88, tokenIndex88, depth88
																	{
																		position93, tokenIndex93, depth93 := position, tokenIndex, depth
																		{
																			switch buffer[position] {
																			case '\\':
																				if buffer[position] != rune('\\') {
																					goto l93
																				}
																				position++
																				break
																			case '\n':
																				if buffer[position] != rune('\n') {
																					goto l93
																				}
																				position++
																				break
																			default:
																				if buffer[position] != rune('"') {
																					goto l93
																				}
																				position++
																				break
																			}
																		}

																		goto l86
																	l93:
																		position, tokenIndex, depth = position93, tokenIndex93, depth93
																	}
																	if !matchDot() {
																		goto l86
																	}
																}
															l88:
																depth--
																add(ruleStringChar, position87)
															}
															goto l85
														l86:
															position, tokenIndex, depth = position86, tokenIndex86, depth86
														}
														depth--
														add(rulePegText, position84)
													}
													if buffer[position] != rune('"') {
														goto l72
													}
													position++
													{
														add(ruleAction21, position)
													}
													{
														add(ruleAction22, position)
													}
													depth--
													add(ruleString, position82)
												}
												break
											default:
												{
													position97 := position
													depth++
													{
														add(ruleAction23, position)
													}
													{
														position99, tokenIndex99, depth99 := position, tokenIndex, depth
														{
															position101 := position
															depth++
															if !_rules[ruleDecimal]() {
																goto l100
															}
															{
																position102, tokenIndex102, depth102 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l103
																}
																position++
																goto l102
															l103:
																position, tokenIndex, depth = position102, tokenIndex102, depth102
																if buffer[position] != rune('E') {
																	goto l100
																}
																position++
															}
														l102:
															if !_rules[ruleInteger]() {
																goto l100
															}
															depth--
															add(ruleSciNum, position101)
														}
														goto l99
													l100:
														position, tokenIndex, depth = position99, tokenIndex99, depth99
														if !_rules[ruleDecimal]() {
															goto l104
														}
														goto l99
													l104:
														position, tokenIndex, depth = position99, tokenIndex99, depth99
														if !_rules[ruleInteger]() {
															goto l72
														}
													}
												l99:
													{
														add(ruleAction24, position)
													}
													depth--
													add(ruleNumeric, position97)
												}
												break
											}
										}

										depth--
										add(ruleScalar, position73)
									}
									goto l58
								l72:
									position, tokenIndex, depth = position58, tokenIndex58, depth58
									{
										position106 := position
										depth++
										{
											switch buffer[position] {
											case '{':
												{
													position108 := position
													depth++
													{
														add(ruleAction40, position)
													}
													if buffer[position] != rune('{') {
														goto l47
													}
													position++
													if !_rules[rulesp]() {
														goto l47
													}
													{
														position110, tokenIndex110, depth110 := position, tokenIndex, depth
														if !_rules[ruleExpr]() {
															goto l110
														}
														if !_rules[rulesp]() {
															goto l110
														}
														if buffer[position] != rune(':') {
															goto l110
														}
														position++
														if !_rules[rulesp]() {
															goto l110
														}
														if !_rules[ruleExpr]() {
															goto l110
														}
													l112:
														{
															position113, tokenIndex113, depth113 := position, tokenIndex, depth
															if !_rules[rulesp]() {
																goto l113
															}
															if buffer[position] != rune(',') {
																goto l113
															}
															position++
															if !_rules[rulesp]() {
																goto l113
															}
															if !_rules[ruleExpr]() {
																goto l113
															}
															if !_rules[rulesp]() {
																goto l113
															}
															if buffer[position] != rune(':') {
																goto l113
															}
															position++
															if !_rules[rulesp]() {
																goto l113
															}
															if !_rules[ruleExpr]() {
																goto l113
															}
															goto l112
														l113:
															position, tokenIndex, depth = position113, tokenIndex113, depth113
														}
														if !_rules[rulesp]() {
															goto l110
														}
														goto l111
													l110:
														position, tokenIndex, depth = position110, tokenIndex110, depth110
													}
												l111:
													if buffer[position] != rune('}') {
														goto l47
													}
													position++
													{
														add(ruleAction41, position)
													}
													depth--
													add(ruleMap, position108)
												}
												break
											case '(':
												if !_rules[ruleTuple]() {
													goto l47
												}
												break
											default:
												{
													position115 := position
													depth++
													{
														add(ruleAction36, position)
													}
													if buffer[position] != rune('[') {
														goto l47
													}
													position++
													if !_rules[rulesp]() {
														goto l47
													}
													{
														position117, tokenIndex117, depth117 := position, tokenIndex, depth
														if !_rules[ruleExpr]() {
															goto l117
														}
													l119:
														{
															position120, tokenIndex120, depth120 := position, tokenIndex, depth
															if !_rules[rulesp]() {
																goto l120
															}
															if buffer[position] != rune(',') {
																goto l120
															}
															position++
															if !_rules[rulesp]() {
																goto l120
															}
															if !_rules[ruleExpr]() {
																goto l120
															}
															goto l119
														l120:
															position, tokenIndex, depth = position120, tokenIndex120, depth120
														}
														if !_rules[rulesp]() {
															goto l117
														}
														goto l118
													l117:
														position, tokenIndex, depth = position117, tokenIndex117, depth117
													}
												l118:
													if buffer[position] != rune(']') {
														goto l47
													}
													position++
													{
														add(ruleAction37, position)
													}
													depth--
													add(ruleList, position115)
												}
												break
											}
										}

										depth--
										add(ruleVector, position106)
									}
								}
							l58:
								depth--
								add(ruleLiteral, position57)
							}
						}
					l55:
						depth--
						add(ruleValue, position54)
					}
				}
			l49:
				depth--
				add(ruleSingle, position48)
			}
			return true
		l47:
			position, tokenIndex, depth = position47, tokenIndex47, depth47
			return false
		},
		/* 5 Expr <- <(Op / Single)> */
		func() bool {
			position122, tokenIndex122, depth122 := position, tokenIndex, depth
			{
				position123 := position
				depth++
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					{
						position126 := position
						depth++
						{
							add(ruleAction2, position)
						}
						if !_rules[ruleSingle]() {
							goto l125
						}
						if !_rules[rulesp]() {
							goto l125
						}
						{
							position130 := position
							depth++
							{
								add(ruleAction4, position)
							}
							{
								position132 := position
								depth++
								{
									position133, tokenIndex133, depth133 := position, tokenIndex, depth
									if buffer[position] != rune('*') {
										goto l134
									}
									position++
									if buffer[position] != rune('*') {
										goto l134
									}
									position++
									goto l133
								l134:
									position, tokenIndex, depth = position133, tokenIndex133, depth133
									if buffer[position] != rune('>') {
										goto l135
									}
									position++
									if buffer[position] != rune('=') {
										goto l135
									}
									position++
									goto l133
								l135:
									position, tokenIndex, depth = position133, tokenIndex133, depth133
									if buffer[position] != rune('<') {
										goto l136
									}
									position++
									if buffer[position] != rune('=') {
										goto l136
									}
									position++
									goto l133
								l136:
									position, tokenIndex, depth = position133, tokenIndex133, depth133
									{
										switch buffer[position] {
										case '<':
											if buffer[position] != rune('<') {
												goto l125
											}
											position++
											break
										case '>':
											if buffer[position] != rune('>') {
												goto l125
											}
											position++
											break
										case '%':
											if buffer[position] != rune('%') {
												goto l125
											}
											position++
											break
										case '/':
											if buffer[position] != rune('/') {
												goto l125
											}
											position++
											break
										case '*':
											if buffer[position] != rune('*') {
												goto l125
											}
											position++
											break
										case '-':
											if buffer[position] != rune('-') {
												goto l125
											}
											position++
											break
										case '+':
											if buffer[position] != rune('+') {
												goto l125
											}
											position++
											break
										default:
											if buffer[position] != rune('=') {
												goto l125
											}
											position++
											if buffer[position] != rune('=') {
												goto l125
											}
											position++
											break
										}
									}

								}
							l133:
								depth--
								add(rulePegText, position132)
							}
							{
								add(ruleAction5, position)
							}
							{
								add(ruleAction6, position)
							}
							depth--
							add(ruleBinaryOp, position130)
						}
						if !_rules[rulesp]() {
							goto l125
						}
						if !_rules[ruleSingle]() {
							goto l125
						}
					l128:
						{
							position129, tokenIndex129, depth129 := position, tokenIndex, depth
							if !_rules[rulesp]() {
								goto l129
							}
							{
								position140 := position
								depth++
								{
									add(ruleAction4, position)
								}
								{
									position142 := position
									depth++
									{
										position143, tokenIndex143, depth143 := position, tokenIndex, depth
										if buffer[position] != rune('*') {
											goto l144
										}
										position++
										if buffer[position] != rune('*') {
											goto l144
										}
										position++
										goto l143
									l144:
										position, tokenIndex, depth = position143, tokenIndex143, depth143
										if buffer[position] != rune('>') {
											goto l145
										}
										position++
										if buffer[position] != rune('=') {
											goto l145
										}
										position++
										goto l143
									l145:
										position, tokenIndex, depth = position143, tokenIndex143, depth143
										if buffer[position] != rune('<') {
											goto l146
										}
										position++
										if buffer[position] != rune('=') {
											goto l146
										}
										position++
										goto l143
									l146:
										position, tokenIndex, depth = position143, tokenIndex143, depth143
										{
											switch buffer[position] {
											case '<':
												if buffer[position] != rune('<') {
													goto l129
												}
												position++
												break
											case '>':
												if buffer[position] != rune('>') {
													goto l129
												}
												position++
												break
											case '%':
												if buffer[position] != rune('%') {
													goto l129
												}
												position++
												break
											case '/':
												if buffer[position] != rune('/') {
													goto l129
												}
												position++
												break
											case '*':
												if buffer[position] != rune('*') {
													goto l129
												}
												position++
												break
											case '-':
												if buffer[position] != rune('-') {
													goto l129
												}
												position++
												break
											case '+':
												if buffer[position] != rune('+') {
													goto l129
												}
												position++
												break
											default:
												if buffer[position] != rune('=') {
													goto l129
												}
												position++
												if buffer[position] != rune('=') {
													goto l129
												}
												position++
												break
											}
										}

									}
								l143:
									depth--
									add(rulePegText, position142)
								}
								{
									add(ruleAction5, position)
								}
								{
									add(ruleAction6, position)
								}
								depth--
								add(ruleBinaryOp, position140)
							}
							if !_rules[rulesp]() {
								goto l129
							}
							if !_rules[ruleSingle]() {
								goto l129
							}
							goto l128
						l129:
							position, tokenIndex, depth = position129, tokenIndex129, depth129
						}
						{
							add(ruleAction3, position)
						}
						depth--
						add(ruleOp, position126)
					}
					goto l124
				l125:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
					if !_rules[ruleSingle]() {
						goto l122
					}
				}
			l124:
				depth--
				add(ruleExpr, position123)
			}
			return true
		l122:
			position, tokenIndex, depth = position122, tokenIndex122, depth122
			return false
		},
		/* 6 Op <- <(Action2 Single (sp BinaryOp sp Single)+ Action3)> */
		nil,
		/* 7 BinaryOp <- <(Action4 <(('*' '*') / ('>' '=') / ('<' '=') / ((&('<') '<') | (&('>') '>') | (&('%') '%') | (&('/') '/') | (&('*') '*') | (&('-') '-') | (&('+') '+') | (&('=') ('=' '='))))> Action5 Action6)> */
		nil,
		/* 8 Statement <- <(Assignment / If)> */
		nil,
		/* 9 Assignment <- <(Action7 LocalRef sp '=' sp Expr Action8)> */
		nil,
		/* 10 If <- <(Action9 ('i' 'f') sp Expr sp Block (sp Else)? Action10)> */
		nil,
		/* 11 Else <- <(Action11 ('e' 'l' 's' 'e') sp Block Action12)> */
		nil,
		/* 12 Ref <- <(FullRef / LocalRef)> */
		func() bool {
			position157, tokenIndex157, depth157 := position, tokenIndex, depth
			{
				position158 := position
				depth++
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					{
						position161 := position
						depth++
						{
							add(ruleAction13, position)
						}
						{
							position163 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l160
							}
						l164:
							{
								position165, tokenIndex165, depth165 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l165
								}
								goto l164
							l165:
								position, tokenIndex, depth = position165, tokenIndex165, depth165
							}
							depth--
							add(rulePegText, position163)
						}
						{
							add(ruleAction14, position)
						}
						if buffer[position] != rune(':') {
							goto l160
						}
						position++
						{
							position167 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l160
							}
						l168:
							{
								position169, tokenIndex169, depth169 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l169
								}
								goto l168
							l169:
								position, tokenIndex, depth = position169, tokenIndex169, depth169
							}
							depth--
							add(rulePegText, position167)
						}
						{
							add(ruleAction15, position)
						}
						{
							add(ruleAction16, position)
						}
						depth--
						add(ruleFullRef, position161)
					}
					goto l159
				l160:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
					if !_rules[ruleLocalRef]() {
						goto l157
					}
				}
			l159:
				depth--
				add(ruleRef, position158)
			}
			return true
		l157:
			position, tokenIndex, depth = position157, tokenIndex157, depth157
			return false
		},
		/* 13 FullRef <- <(Action13 <RefChar+> Action14 ':' <RefChar+> Action15 Action16)> */
		nil,
		/* 14 LocalRef <- <(Action17 <RefChar+> Action18 Action19)> */
		func() bool {
			position173, tokenIndex173, depth173 := position, tokenIndex, depth
			{
				position174 := position
				depth++
				{
					add(ruleAction17, position)
				}
				{
					position176 := position
					depth++
					if !_rules[ruleRefChar]() {
						goto l173
					}
				l177:
					{
						position178, tokenIndex178, depth178 := position, tokenIndex, depth
						if !_rules[ruleRefChar]() {
							goto l178
						}
						goto l177
					l178:
						position, tokenIndex, depth = position178, tokenIndex178, depth178
					}
					depth--
					add(rulePegText, position176)
				}
				{
					add(ruleAction18, position)
				}
				{
					add(ruleAction19, position)
				}
				depth--
				add(ruleLocalRef, position174)
			}
			return true
		l173:
			position, tokenIndex, depth = position173, tokenIndex173, depth173
			return false
		},
		/* 15 RefChar <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position181, tokenIndex181, depth181 := position, tokenIndex, depth
			{
				position182 := position
				depth++
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l181
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l181
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l181
						}
						position++
						break
					}
				}

				depth--
				add(ruleRefChar, position182)
			}
			return true
		l181:
			position, tokenIndex, depth = position181, tokenIndex181, depth181
			return false
		},
		/* 16 Value <- <(Ref / Literal)> */
		nil,
		/* 17 Literal <- <(Func / Scalar / Vector)> */
		nil,
		/* 18 Scalar <- <((&('f' | 't') Boolean) | (&('"') String) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric))> */
		nil,
		/* 19 Vector <- <((&('{') Map) | (&('(') Tuple) | (&('[') List))> */
		nil,
		/* 20 String <- <(Action20 '"' <StringChar*> '"' Action21 Action22)> */
		nil,
		/* 21 StringChar <- <(StringEsc / (!((&('\\') '\\') | (&('\n') '\n') | (&('"') '"')) .))> */
		nil,
		/* 22 StringEsc <- <SimpleEsc> */
		nil,
		/* 23 SimpleEsc <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 24 Numeric <- <(Action23 (SciNum / Decimal / Integer) Action24)> */
		nil,
		/* 25 SciNum <- <(Decimal ('e' / 'E') Integer)> */
		nil,
		/* 26 Decimal <- <(Integer '.' <Digit*> Action25)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				if !_rules[ruleInteger]() {
					goto l194
				}
				if buffer[position] != rune('.') {
					goto l194
				}
				position++
				{
					position196 := position
					depth++
				l197:
					{
						position198, tokenIndex198, depth198 := position, tokenIndex, depth
						if !_rules[ruleDigit]() {
							goto l198
						}
						goto l197
					l198:
						position, tokenIndex, depth = position198, tokenIndex198, depth198
					}
					depth--
					add(rulePegText, position196)
				}
				{
					add(ruleAction25, position)
				}
				depth--
				add(ruleDecimal, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 27 Integer <- <(<WholeNum> Action26)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				{
					position202 := position
					depth++
					{
						position203 := position
						depth++
						{
							position204, tokenIndex204, depth204 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l205
							}
							position++
							goto l204
						l205:
							position, tokenIndex, depth = position204, tokenIndex204, depth204
							{
								position206, tokenIndex206, depth206 := position, tokenIndex, depth
								if buffer[position] != rune('-') {
									goto l206
								}
								position++
								goto l207
							l206:
								position, tokenIndex, depth = position206, tokenIndex206, depth206
							}
						l207:
							if c := buffer[position]; c < rune('1') || c > rune('9') {
								goto l200
							}
							position++
						l208:
							{
								position209, tokenIndex209, depth209 := position, tokenIndex, depth
								if !_rules[ruleDigit]() {
									goto l209
								}
								goto l208
							l209:
								position, tokenIndex, depth = position209, tokenIndex209, depth209
							}
						}
					l204:
						depth--
						add(ruleWholeNum, position203)
					}
					depth--
					add(rulePegText, position202)
				}
				{
					add(ruleAction26, position)
				}
				depth--
				add(ruleInteger, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 28 WholeNum <- <('0' / ('-'? [1-9] Digit*))> */
		nil,
		/* 29 Digit <- <[0-9]> */
		func() bool {
			position212, tokenIndex212, depth212 := position, tokenIndex, depth
			{
				position213 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l212
				}
				position++
				depth--
				add(ruleDigit, position213)
			}
			return true
		l212:
			position, tokenIndex, depth = position212, tokenIndex212, depth212
			return false
		},
		/* 30 Boolean <- <(Action27 <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> Action28 Action29)> */
		nil,
		/* 31 Func <- <(Action30 FuncArgs sp ('-' '>') sp (Block / Expr) Action31)> */
		nil,
		/* 32 FuncArgs <- <(Action32 '(' sp (LocalRef (sp ',' sp LocalRef)* sp)? ')' Action33)> */
		nil,
		/* 33 FuncApply <- <(Action34 Ref Tuple Action35)> */
		nil,
		/* 34 List <- <(Action36 '[' sp (Expr (sp ',' sp Expr)* sp)? ']' Action37)> */
		nil,
		/* 35 Tuple <- <(Action38 '(' sp (Expr (sp ',' sp Expr)* sp)? ')' Action39)> */
		func() bool {
			position219, tokenIndex219, depth219 := position, tokenIndex, depth
			{
				position220 := position
				depth++
				{
					add(ruleAction38, position)
				}
				if buffer[position] != rune('(') {
					goto l219
				}
				position++
				if !_rules[rulesp]() {
					goto l219
				}
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if !_rules[ruleExpr]() {
						goto l222
					}
				l224:
					{
						position225, tokenIndex225, depth225 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l225
						}
						if buffer[position] != rune(',') {
							goto l225
						}
						position++
						if !_rules[rulesp]() {
							goto l225
						}
						if !_rules[ruleExpr]() {
							goto l225
						}
						goto l224
					l225:
						position, tokenIndex, depth = position225, tokenIndex225, depth225
					}
					if !_rules[rulesp]() {
						goto l222
					}
					goto l223
				l222:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
				}
			l223:
				if buffer[position] != rune(')') {
					goto l219
				}
				position++
				{
					add(ruleAction39, position)
				}
				depth--
				add(ruleTuple, position220)
			}
			return true
		l219:
			position, tokenIndex, depth = position219, tokenIndex219, depth219
			return false
		},
		/* 36 Map <- <(Action40 '{' sp (Expr sp ':' sp Expr (sp ',' sp Expr sp ':' sp Expr)* sp)? '}' Action41)> */
		nil,
		/* 37 Gravitasse <- <'@'> */
		nil,
		/* 38 msp <- <(ws / comment)+> */
		nil,
		/* 39 sp <- <(ws / comment)*> */
		func() bool {
			{
				position231 := position
				depth++
			l232:
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					{
						position234, tokenIndex234, depth234 := position, tokenIndex, depth
						if !_rules[rulews]() {
							goto l235
						}
						goto l234
					l235:
						position, tokenIndex, depth = position234, tokenIndex234, depth234
						if !_rules[rulecomment]() {
							goto l233
						}
					}
				l234:
					goto l232
				l233:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
				}
				depth--
				add(rulesp, position231)
			}
			return true
		},
		/* 40 comment <- <('#' (!'\n' .)*)> */
		func() bool {
			position236, tokenIndex236, depth236 := position, tokenIndex, depth
			{
				position237 := position
				depth++
				if buffer[position] != rune('#') {
					goto l236
				}
				position++
			l238:
				{
					position239, tokenIndex239, depth239 := position, tokenIndex, depth
					{
						position240, tokenIndex240, depth240 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l240
						}
						position++
						goto l239
					l240:
						position, tokenIndex, depth = position240, tokenIndex240, depth240
					}
					if !matchDot() {
						goto l239
					}
					goto l238
				l239:
					position, tokenIndex, depth = position239, tokenIndex239, depth239
				}
				depth--
				add(rulecomment, position237)
			}
			return true
		l236:
			position, tokenIndex, depth = position236, tokenIndex236, depth236
			return false
		},
		/* 41 ws <- <((&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))> */
		func() bool {
			position241, tokenIndex241, depth241 := position, tokenIndex, depth
			{
				position242 := position
				depth++
				{
					switch buffer[position] {
					case '\r':
						if buffer[position] != rune('\r') {
							goto l241
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l241
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l241
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l241
						}
						position++
						break
					}
				}

				depth--
				add(rulews, position242)
			}
			return true
		l241:
			position, tokenIndex, depth = position241, tokenIndex241, depth241
			return false
		},
		/* 43 Action0 <- <{ p.Start(RIFT) }> */
		nil,
		/* 44 Action1 <- <{ p.End() }> */
		nil,
		/* 45 Action2 <- <{ p.Start(OP) }> */
		nil,
		/* 46 Action3 <- <{ p.End() }> */
		nil,
		/* 47 Action4 <- <{ p.Start(BINOP) }> */
		nil,
		nil,
		/* 49 Action5 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 50 Action6 <- <{ p.End() }> */
		nil,
		/* 51 Action7 <- <{ p.Start(ASSIGNMENT) }> */
		nil,
		/* 52 Action8 <- <{ p.End() }> */
		nil,
		/* 53 Action9 <- <{ p.Start(IF) }> */
		nil,
		/* 54 Action10 <- <{ p.End() }> */
		nil,
		/* 55 Action11 <- <{ p.Start(ELSE) }> */
		nil,
		/* 56 Action12 <- <{ p.End() }> */
		nil,
		/* 57 Action13 <- <{ p.Start(REF) }> */
		nil,
		/* 58 Action14 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 59 Action15 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 60 Action16 <- <{ p.End() }> */
		nil,
		/* 61 Action17 <- <{ p.Start(REF) }> */
		nil,
		/* 62 Action18 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 63 Action19 <- <{ p.End() }> */
		nil,
		/* 64 Action20 <- <{ p.Start(STRING) }> */
		nil,
		/* 65 Action21 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 66 Action22 <- <{ p.End() }> */
		nil,
		/* 67 Action23 <- <{ p.Start(NUM) }> */
		nil,
		/* 68 Action24 <- <{ p.End() }> */
		nil,
		/* 69 Action25 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 70 Action26 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 71 Action27 <- <{ p.Start(BOOL) }> */
		nil,
		/* 72 Action28 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 73 Action29 <- <{ p.End() }> */
		nil,
		/* 74 Action30 <- <{ p.Start(FUNC) }> */
		nil,
		/* 75 Action31 <- <{ p.End() }> */
		nil,
		/* 76 Action32 <- <{ p.Start(ARGS) }> */
		nil,
		/* 77 Action33 <- <{ p.End() }> */
		nil,
		/* 78 Action34 <- <{ p.Start(FUNCAPPLY) }> */
		nil,
		/* 79 Action35 <- <{ p.End() }> */
		nil,
		/* 80 Action36 <- <{ p.Start(LIST) }> */
		nil,
		/* 81 Action37 <- <{ p.End() }> */
		nil,
		/* 82 Action38 <- <{ p.Start(TUPLE) }> */
		nil,
		/* 83 Action39 <- <{ p.End() }> */
		nil,
		/* 84 Action40 <- <{ p.Start("map") }> */
		nil,
		/* 85 Action41 <- <{ p.End() }> */
		nil,
	}
	p.rules = _rules
}
