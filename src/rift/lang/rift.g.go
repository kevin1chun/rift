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
	ruleExpr
	ruleOp
	ruleBinaryOp
	ruleAssignment
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
	"Expr",
	"Op",
	"BinaryOp",
	"Assignment",
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
	rules  [75]func() bool
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
			p.Start(REF)
		case ruleAction10:
			p.End()
		case ruleAction11:
			p.Emit(string(buffer[begin:end]))
		case ruleAction12:
			p.Emit(string(buffer[begin:end]))
		case ruleAction13:
			p.Emit(string(buffer[begin:end]))
		case ruleAction14:
			p.Start(STRING)
		case ruleAction15:
			p.Emit(string(buffer[begin:end]))
		case ruleAction16:
			p.End()
		case ruleAction17:
			p.Start(NUM)
		case ruleAction18:
			p.End()
		case ruleAction19:
			p.Emit(string(buffer[begin:end]))
		case ruleAction20:
			p.Emit(string(buffer[begin:end]))
		case ruleAction21:
			p.Start(BOOL)
		case ruleAction22:
			p.Emit(string(buffer[begin:end]))
		case ruleAction23:
			p.End()
		case ruleAction24:
			p.Start(FUNC)
		case ruleAction25:
			p.End()
		case ruleAction26:
			p.Start(ARGS)
		case ruleAction27:
			p.End()
		case ruleAction28:
			p.Start(FUNCAPPLY)
		case ruleAction29:
			p.End()
		case ruleAction30:
			p.Start(LIST)
		case ruleAction31:
			p.End()
		case ruleAction32:
			p.Start(TUPLE)
		case ruleAction33:
			p.End()
		case ruleAction34:
			p.Start("map")
		case ruleAction35:
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
						if buffer[position] != rune('@') {
							goto l6
						}
						position++
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
						position9 := position
						depth++
						{
							add(ruleAction0, position)
						}
						{
							position11, tokenIndex11, depth11 := position, tokenIndex, depth
							if buffer[position] != rune('@') {
								goto l11
							}
							position++
							goto l12
						l11:
							position, tokenIndex, depth = position11, tokenIndex11, depth11
						}
					l12:
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
						add(ruleRift, position9)
					}
					if !_rules[rulesp]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				{
					position14, tokenIndex14, depth14 := position, tokenIndex, depth
					if !matchDot() {
						goto l14
					}
					goto l0
				l14:
					position, tokenIndex, depth = position14, tokenIndex14, depth14
				}
				depth--
				add(ruleSource, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Rift <- <(Action0 '@'? LocalRef sp ('=' '>') sp Block Action1)> */
		nil,
		/* 2 Block <- <('{' sp (Line msp)* '}')> */
		func() bool {
			position16, tokenIndex16, depth16 := position, tokenIndex, depth
			{
				position17 := position
				depth++
				if buffer[position] != rune('{') {
					goto l16
				}
				position++
				if !_rules[rulesp]() {
					goto l16
				}
			l18:
				{
					position19, tokenIndex19, depth19 := position, tokenIndex, depth
					{
						position20 := position
						depth++
						{
							position21, tokenIndex21, depth21 := position, tokenIndex, depth
							{
								position23 := position
								depth++
								{
									add(ruleAction7, position)
								}
								if !_rules[ruleLocalRef]() {
									goto l22
								}
								if !_rules[rulesp]() {
									goto l22
								}
								if buffer[position] != rune('=') {
									goto l22
								}
								position++
								if !_rules[rulesp]() {
									goto l22
								}
								if !_rules[ruleExpr]() {
									goto l22
								}
								{
									add(ruleAction8, position)
								}
								depth--
								add(ruleAssignment, position23)
							}
							goto l21
						l22:
							position, tokenIndex, depth = position21, tokenIndex21, depth21
							if !_rules[ruleExpr]() {
								goto l19
							}
						}
					l21:
						depth--
						add(ruleLine, position20)
					}
					{
						position26 := position
						depth++
						{
							position29, tokenIndex29, depth29 := position, tokenIndex, depth
							if !_rules[rulews]() {
								goto l30
							}
							goto l29
						l30:
							position, tokenIndex, depth = position29, tokenIndex29, depth29
							if !_rules[rulecomment]() {
								goto l19
							}
						}
					l29:
					l27:
						{
							position28, tokenIndex28, depth28 := position, tokenIndex, depth
							{
								position31, tokenIndex31, depth31 := position, tokenIndex, depth
								if !_rules[rulews]() {
									goto l32
								}
								goto l31
							l32:
								position, tokenIndex, depth = position31, tokenIndex31, depth31
								if !_rules[rulecomment]() {
									goto l28
								}
							}
						l31:
							goto l27
						l28:
							position, tokenIndex, depth = position28, tokenIndex28, depth28
						}
						depth--
						add(rulemsp, position26)
					}
					goto l18
				l19:
					position, tokenIndex, depth = position19, tokenIndex19, depth19
				}
				if buffer[position] != rune('}') {
					goto l16
				}
				position++
				depth--
				add(ruleBlock, position17)
			}
			return true
		l16:
			position, tokenIndex, depth = position16, tokenIndex16, depth16
			return false
		},
		/* 3 Line <- <(Assignment / Expr)> */
		nil,
		/* 4 Expr <- <(Op / FuncApply / Value)> */
		func() bool {
			position34, tokenIndex34, depth34 := position, tokenIndex, depth
			{
				position35 := position
				depth++
				{
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					{
						position38 := position
						depth++
						{
							add(ruleAction2, position)
						}
						if !_rules[ruleValue]() {
							goto l37
						}
						if !_rules[rulesp]() {
							goto l37
						}
						{
							position42 := position
							depth++
							{
								add(ruleAction4, position)
							}
							{
								position44 := position
								depth++
								{
									position45, tokenIndex45, depth45 := position, tokenIndex, depth
									if buffer[position] != rune('*') {
										goto l46
									}
									position++
									goto l45
								l46:
									position, tokenIndex, depth = position45, tokenIndex45, depth45
									{
										switch buffer[position] {
										case '%':
											if buffer[position] != rune('%') {
												goto l37
											}
											position++
											break
										case '*':
											if buffer[position] != rune('*') {
												goto l37
											}
											position++
											if buffer[position] != rune('*') {
												goto l37
											}
											position++
											break
										case '/':
											if buffer[position] != rune('/') {
												goto l37
											}
											position++
											break
										case '-':
											if buffer[position] != rune('-') {
												goto l37
											}
											position++
											break
										default:
											if buffer[position] != rune('+') {
												goto l37
											}
											position++
											break
										}
									}

								}
							l45:
								depth--
								add(rulePegText, position44)
							}
							{
								add(ruleAction5, position)
							}
							{
								add(ruleAction6, position)
							}
							depth--
							add(ruleBinaryOp, position42)
						}
						if !_rules[rulesp]() {
							goto l37
						}
						if !_rules[ruleValue]() {
							goto l37
						}
					l40:
						{
							position41, tokenIndex41, depth41 := position, tokenIndex, depth
							if !_rules[rulesp]() {
								goto l41
							}
							{
								position50 := position
								depth++
								{
									add(ruleAction4, position)
								}
								{
									position52 := position
									depth++
									{
										position53, tokenIndex53, depth53 := position, tokenIndex, depth
										if buffer[position] != rune('*') {
											goto l54
										}
										position++
										goto l53
									l54:
										position, tokenIndex, depth = position53, tokenIndex53, depth53
										{
											switch buffer[position] {
											case '%':
												if buffer[position] != rune('%') {
													goto l41
												}
												position++
												break
											case '*':
												if buffer[position] != rune('*') {
													goto l41
												}
												position++
												if buffer[position] != rune('*') {
													goto l41
												}
												position++
												break
											case '/':
												if buffer[position] != rune('/') {
													goto l41
												}
												position++
												break
											case '-':
												if buffer[position] != rune('-') {
													goto l41
												}
												position++
												break
											default:
												if buffer[position] != rune('+') {
													goto l41
												}
												position++
												break
											}
										}

									}
								l53:
									depth--
									add(rulePegText, position52)
								}
								{
									add(ruleAction5, position)
								}
								{
									add(ruleAction6, position)
								}
								depth--
								add(ruleBinaryOp, position50)
							}
							if !_rules[rulesp]() {
								goto l41
							}
							if !_rules[ruleValue]() {
								goto l41
							}
							goto l40
						l41:
							position, tokenIndex, depth = position41, tokenIndex41, depth41
						}
						{
							add(ruleAction3, position)
						}
						depth--
						add(ruleOp, position38)
					}
					goto l36
				l37:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					{
						position60 := position
						depth++
						{
							add(ruleAction28, position)
						}
						if !_rules[ruleRef]() {
							goto l59
						}
						if !_rules[ruleTuple]() {
							goto l59
						}
						{
							add(ruleAction29, position)
						}
						depth--
						add(ruleFuncApply, position60)
					}
					goto l36
				l59:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					if !_rules[ruleValue]() {
						goto l34
					}
				}
			l36:
				depth--
				add(ruleExpr, position35)
			}
			return true
		l34:
			position, tokenIndex, depth = position34, tokenIndex34, depth34
			return false
		},
		/* 5 Op <- <(Action2 Value (sp BinaryOp sp Value)+ Action3)> */
		nil,
		/* 6 BinaryOp <- <(Action4 <('*' / ((&('%') '%') | (&('*') ('*' '*')) | (&('/') '/') | (&('-') '-') | (&('+') '+')))> Action5 Action6)> */
		nil,
		/* 7 Assignment <- <(Action7 LocalRef sp '=' sp Expr Action8)> */
		nil,
		/* 8 Ref <- <(Action9 (FullRef / LocalRef) Action10)> */
		func() bool {
			position66, tokenIndex66, depth66 := position, tokenIndex, depth
			{
				position67 := position
				depth++
				{
					add(ruleAction9, position)
				}
				{
					position69, tokenIndex69, depth69 := position, tokenIndex, depth
					{
						position71 := position
						depth++
						{
							position72 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l70
							}
						l73:
							{
								position74, tokenIndex74, depth74 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l74
								}
								goto l73
							l74:
								position, tokenIndex, depth = position74, tokenIndex74, depth74
							}
							depth--
							add(rulePegText, position72)
						}
						{
							add(ruleAction11, position)
						}
						if buffer[position] != rune(':') {
							goto l70
						}
						position++
						{
							position76 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l70
							}
						l77:
							{
								position78, tokenIndex78, depth78 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l78
								}
								goto l77
							l78:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
							}
							depth--
							add(rulePegText, position76)
						}
						{
							add(ruleAction12, position)
						}
						depth--
						add(ruleFullRef, position71)
					}
					goto l69
				l70:
					position, tokenIndex, depth = position69, tokenIndex69, depth69
					if !_rules[ruleLocalRef]() {
						goto l66
					}
				}
			l69:
				{
					add(ruleAction10, position)
				}
				depth--
				add(ruleRef, position67)
			}
			return true
		l66:
			position, tokenIndex, depth = position66, tokenIndex66, depth66
			return false
		},
		/* 9 FullRef <- <(<RefChar+> Action11 ':' <RefChar+> Action12)> */
		nil,
		/* 10 LocalRef <- <(<RefChar+> Action13)> */
		func() bool {
			position82, tokenIndex82, depth82 := position, tokenIndex, depth
			{
				position83 := position
				depth++
				{
					position84 := position
					depth++
					if !_rules[ruleRefChar]() {
						goto l82
					}
				l85:
					{
						position86, tokenIndex86, depth86 := position, tokenIndex, depth
						if !_rules[ruleRefChar]() {
							goto l86
						}
						goto l85
					l86:
						position, tokenIndex, depth = position86, tokenIndex86, depth86
					}
					depth--
					add(rulePegText, position84)
				}
				{
					add(ruleAction13, position)
				}
				depth--
				add(ruleLocalRef, position83)
			}
			return true
		l82:
			position, tokenIndex, depth = position82, tokenIndex82, depth82
			return false
		},
		/* 11 RefChar <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position88, tokenIndex88, depth88 := position, tokenIndex, depth
			{
				position89 := position
				depth++
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l88
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l88
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l88
						}
						position++
						break
					}
				}

				depth--
				add(ruleRefChar, position89)
			}
			return true
		l88:
			position, tokenIndex, depth = position88, tokenIndex88, depth88
			return false
		},
		/* 12 Value <- <(Ref / Literal)> */
		func() bool {
			position91, tokenIndex91, depth91 := position, tokenIndex, depth
			{
				position92 := position
				depth++
				{
					position93, tokenIndex93, depth93 := position, tokenIndex, depth
					if !_rules[ruleRef]() {
						goto l94
					}
					goto l93
				l94:
					position, tokenIndex, depth = position93, tokenIndex93, depth93
					{
						position95 := position
						depth++
						{
							position96, tokenIndex96, depth96 := position, tokenIndex, depth
							{
								position98 := position
								depth++
								{
									add(ruleAction24, position)
								}
								{
									position100 := position
									depth++
									{
										add(ruleAction26, position)
									}
									if buffer[position] != rune('(') {
										goto l97
									}
									position++
									if !_rules[rulesp]() {
										goto l97
									}
									{
										position102, tokenIndex102, depth102 := position, tokenIndex, depth
										if !_rules[ruleLocalRef]() {
											goto l102
										}
									l104:
										{
											position105, tokenIndex105, depth105 := position, tokenIndex, depth
											if !_rules[rulesp]() {
												goto l105
											}
											if buffer[position] != rune(',') {
												goto l105
											}
											position++
											if !_rules[rulesp]() {
												goto l105
											}
											if !_rules[ruleLocalRef]() {
												goto l105
											}
											goto l104
										l105:
											position, tokenIndex, depth = position105, tokenIndex105, depth105
										}
										if !_rules[rulesp]() {
											goto l102
										}
										goto l103
									l102:
										position, tokenIndex, depth = position102, tokenIndex102, depth102
									}
								l103:
									if buffer[position] != rune(')') {
										goto l97
									}
									position++
									{
										add(ruleAction27, position)
									}
									depth--
									add(ruleFuncArgs, position100)
								}
								if !_rules[rulesp]() {
									goto l97
								}
								if buffer[position] != rune('-') {
									goto l97
								}
								position++
								if buffer[position] != rune('>') {
									goto l97
								}
								position++
								if !_rules[rulesp]() {
									goto l97
								}
								{
									position107, tokenIndex107, depth107 := position, tokenIndex, depth
									if !_rules[ruleBlock]() {
										goto l108
									}
									goto l107
								l108:
									position, tokenIndex, depth = position107, tokenIndex107, depth107
									if !_rules[ruleExpr]() {
										goto l97
									}
								}
							l107:
								{
									add(ruleAction25, position)
								}
								depth--
								add(ruleFunc, position98)
							}
							goto l96
						l97:
							position, tokenIndex, depth = position96, tokenIndex96, depth96
							{
								position111 := position
								depth++
								{
									switch buffer[position] {
									case 'f', 't':
										{
											position113 := position
											depth++
											{
												add(ruleAction21, position)
											}
											{
												position115 := position
												depth++
												{
													position116, tokenIndex116, depth116 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l117
													}
													position++
													if buffer[position] != rune('r') {
														goto l117
													}
													position++
													if buffer[position] != rune('u') {
														goto l117
													}
													position++
													if buffer[position] != rune('e') {
														goto l117
													}
													position++
													goto l116
												l117:
													position, tokenIndex, depth = position116, tokenIndex116, depth116
													if buffer[position] != rune('f') {
														goto l110
													}
													position++
													if buffer[position] != rune('a') {
														goto l110
													}
													position++
													if buffer[position] != rune('l') {
														goto l110
													}
													position++
													if buffer[position] != rune('s') {
														goto l110
													}
													position++
													if buffer[position] != rune('e') {
														goto l110
													}
													position++
												}
											l116:
												depth--
												add(rulePegText, position115)
											}
											{
												add(ruleAction22, position)
											}
											{
												add(ruleAction23, position)
											}
											depth--
											add(ruleBoolean, position113)
										}
										break
									case '"':
										{
											position120 := position
											depth++
											{
												add(ruleAction14, position)
											}
											if buffer[position] != rune('"') {
												goto l110
											}
											position++
											{
												position122 := position
												depth++
											l123:
												{
													position124, tokenIndex124, depth124 := position, tokenIndex, depth
													{
														position125 := position
														depth++
														{
															position126, tokenIndex126, depth126 := position, tokenIndex, depth
															{
																position128 := position
																depth++
																{
																	position129 := position
																	depth++
																	if buffer[position] != rune('\\') {
																		goto l127
																	}
																	position++
																	{
																		switch buffer[position] {
																		case 'v':
																			if buffer[position] != rune('v') {
																				goto l127
																			}
																			position++
																			break
																		case 't':
																			if buffer[position] != rune('t') {
																				goto l127
																			}
																			position++
																			break
																		case 'r':
																			if buffer[position] != rune('r') {
																				goto l127
																			}
																			position++
																			break
																		case 'n':
																			if buffer[position] != rune('n') {
																				goto l127
																			}
																			position++
																			break
																		case 'f':
																			if buffer[position] != rune('f') {
																				goto l127
																			}
																			position++
																			break
																		case 'b':
																			if buffer[position] != rune('b') {
																				goto l127
																			}
																			position++
																			break
																		case 'a':
																			if buffer[position] != rune('a') {
																				goto l127
																			}
																			position++
																			break
																		case '\\':
																			if buffer[position] != rune('\\') {
																				goto l127
																			}
																			position++
																			break
																		case '?':
																			if buffer[position] != rune('?') {
																				goto l127
																			}
																			position++
																			break
																		case '"':
																			if buffer[position] != rune('"') {
																				goto l127
																			}
																			position++
																			break
																		default:
																			if buffer[position] != rune('\'') {
																				goto l127
																			}
																			position++
																			break
																		}
																	}

																	depth--
																	add(ruleSimpleEsc, position129)
																}
																depth--
																add(ruleStringEsc, position128)
															}
															goto l126
														l127:
															position, tokenIndex, depth = position126, tokenIndex126, depth126
															{
																position131, tokenIndex131, depth131 := position, tokenIndex, depth
																{
																	switch buffer[position] {
																	case '\\':
																		if buffer[position] != rune('\\') {
																			goto l131
																		}
																		position++
																		break
																	case '\n':
																		if buffer[position] != rune('\n') {
																			goto l131
																		}
																		position++
																		break
																	default:
																		if buffer[position] != rune('"') {
																			goto l131
																		}
																		position++
																		break
																	}
																}

																goto l124
															l131:
																position, tokenIndex, depth = position131, tokenIndex131, depth131
															}
															if !matchDot() {
																goto l124
															}
														}
													l126:
														depth--
														add(ruleStringChar, position125)
													}
													goto l123
												l124:
													position, tokenIndex, depth = position124, tokenIndex124, depth124
												}
												depth--
												add(rulePegText, position122)
											}
											if buffer[position] != rune('"') {
												goto l110
											}
											position++
											{
												add(ruleAction15, position)
											}
											{
												add(ruleAction16, position)
											}
											depth--
											add(ruleString, position120)
										}
										break
									default:
										{
											position135 := position
											depth++
											{
												add(ruleAction17, position)
											}
											{
												position137, tokenIndex137, depth137 := position, tokenIndex, depth
												{
													position139 := position
													depth++
													if !_rules[ruleDecimal]() {
														goto l138
													}
													{
														position140, tokenIndex140, depth140 := position, tokenIndex, depth
														if buffer[position] != rune('e') {
															goto l141
														}
														position++
														goto l140
													l141:
														position, tokenIndex, depth = position140, tokenIndex140, depth140
														if buffer[position] != rune('E') {
															goto l138
														}
														position++
													}
												l140:
													if !_rules[ruleInteger]() {
														goto l138
													}
													depth--
													add(ruleSciNum, position139)
												}
												goto l137
											l138:
												position, tokenIndex, depth = position137, tokenIndex137, depth137
												if !_rules[ruleDecimal]() {
													goto l142
												}
												goto l137
											l142:
												position, tokenIndex, depth = position137, tokenIndex137, depth137
												if !_rules[ruleInteger]() {
													goto l110
												}
											}
										l137:
											{
												add(ruleAction18, position)
											}
											depth--
											add(ruleNumeric, position135)
										}
										break
									}
								}

								depth--
								add(ruleScalar, position111)
							}
							goto l96
						l110:
							position, tokenIndex, depth = position96, tokenIndex96, depth96
							{
								position144 := position
								depth++
								{
									switch buffer[position] {
									case '{':
										{
											position146 := position
											depth++
											{
												add(ruleAction34, position)
											}
											if buffer[position] != rune('{') {
												goto l91
											}
											position++
											if !_rules[rulesp]() {
												goto l91
											}
											{
												position148, tokenIndex148, depth148 := position, tokenIndex, depth
												if !_rules[ruleExpr]() {
													goto l148
												}
												if !_rules[rulesp]() {
													goto l148
												}
												if buffer[position] != rune(':') {
													goto l148
												}
												position++
												if !_rules[rulesp]() {
													goto l148
												}
												if !_rules[ruleExpr]() {
													goto l148
												}
											l150:
												{
													position151, tokenIndex151, depth151 := position, tokenIndex, depth
													if !_rules[rulesp]() {
														goto l151
													}
													if buffer[position] != rune(',') {
														goto l151
													}
													position++
													if !_rules[rulesp]() {
														goto l151
													}
													if !_rules[ruleExpr]() {
														goto l151
													}
													if !_rules[rulesp]() {
														goto l151
													}
													if buffer[position] != rune(':') {
														goto l151
													}
													position++
													if !_rules[rulesp]() {
														goto l151
													}
													if !_rules[ruleExpr]() {
														goto l151
													}
													goto l150
												l151:
													position, tokenIndex, depth = position151, tokenIndex151, depth151
												}
												if !_rules[rulesp]() {
													goto l148
												}
												goto l149
											l148:
												position, tokenIndex, depth = position148, tokenIndex148, depth148
											}
										l149:
											if buffer[position] != rune('}') {
												goto l91
											}
											position++
											{
												add(ruleAction35, position)
											}
											depth--
											add(ruleMap, position146)
										}
										break
									case '(':
										if !_rules[ruleTuple]() {
											goto l91
										}
										break
									default:
										{
											position153 := position
											depth++
											{
												add(ruleAction30, position)
											}
											if buffer[position] != rune('[') {
												goto l91
											}
											position++
											if !_rules[rulesp]() {
												goto l91
											}
											{
												position155, tokenIndex155, depth155 := position, tokenIndex, depth
												if !_rules[ruleExpr]() {
													goto l155
												}
											l157:
												{
													position158, tokenIndex158, depth158 := position, tokenIndex, depth
													if !_rules[rulesp]() {
														goto l158
													}
													if buffer[position] != rune(',') {
														goto l158
													}
													position++
													if !_rules[rulesp]() {
														goto l158
													}
													if !_rules[ruleExpr]() {
														goto l158
													}
													goto l157
												l158:
													position, tokenIndex, depth = position158, tokenIndex158, depth158
												}
												if !_rules[rulesp]() {
													goto l155
												}
												goto l156
											l155:
												position, tokenIndex, depth = position155, tokenIndex155, depth155
											}
										l156:
											if buffer[position] != rune(']') {
												goto l91
											}
											position++
											{
												add(ruleAction31, position)
											}
											depth--
											add(ruleList, position153)
										}
										break
									}
								}

								depth--
								add(ruleVector, position144)
							}
						}
					l96:
						depth--
						add(ruleLiteral, position95)
					}
				}
			l93:
				depth--
				add(ruleValue, position92)
			}
			return true
		l91:
			position, tokenIndex, depth = position91, tokenIndex91, depth91
			return false
		},
		/* 13 Literal <- <(Func / Scalar / Vector)> */
		nil,
		/* 14 Scalar <- <((&('f' | 't') Boolean) | (&('"') String) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric))> */
		nil,
		/* 15 Vector <- <((&('{') Map) | (&('(') Tuple) | (&('[') List))> */
		nil,
		/* 16 String <- <(Action14 '"' <StringChar*> '"' Action15 Action16)> */
		nil,
		/* 17 StringChar <- <(StringEsc / (!((&('\\') '\\') | (&('\n') '\n') | (&('"') '"')) .))> */
		nil,
		/* 18 StringEsc <- <SimpleEsc> */
		nil,
		/* 19 SimpleEsc <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 20 Numeric <- <(Action17 (SciNum / Decimal / Integer) Action18)> */
		nil,
		/* 21 SciNum <- <(Decimal ('e' / 'E') Integer)> */
		nil,
		/* 22 Decimal <- <(Integer '.' <Digit*> Action19)> */
		func() bool {
			position169, tokenIndex169, depth169 := position, tokenIndex, depth
			{
				position170 := position
				depth++
				if !_rules[ruleInteger]() {
					goto l169
				}
				if buffer[position] != rune('.') {
					goto l169
				}
				position++
				{
					position171 := position
					depth++
				l172:
					{
						position173, tokenIndex173, depth173 := position, tokenIndex, depth
						if !_rules[ruleDigit]() {
							goto l173
						}
						goto l172
					l173:
						position, tokenIndex, depth = position173, tokenIndex173, depth173
					}
					depth--
					add(rulePegText, position171)
				}
				{
					add(ruleAction19, position)
				}
				depth--
				add(ruleDecimal, position170)
			}
			return true
		l169:
			position, tokenIndex, depth = position169, tokenIndex169, depth169
			return false
		},
		/* 23 Integer <- <(<WholeNum> Action20)> */
		func() bool {
			position175, tokenIndex175, depth175 := position, tokenIndex, depth
			{
				position176 := position
				depth++
				{
					position177 := position
					depth++
					{
						position178 := position
						depth++
						{
							position179, tokenIndex179, depth179 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l180
							}
							position++
							goto l179
						l180:
							position, tokenIndex, depth = position179, tokenIndex179, depth179
							{
								position181, tokenIndex181, depth181 := position, tokenIndex, depth
								if buffer[position] != rune('-') {
									goto l181
								}
								position++
								goto l182
							l181:
								position, tokenIndex, depth = position181, tokenIndex181, depth181
							}
						l182:
							if c := buffer[position]; c < rune('1') || c > rune('9') {
								goto l175
							}
							position++
						l183:
							{
								position184, tokenIndex184, depth184 := position, tokenIndex, depth
								if !_rules[ruleDigit]() {
									goto l184
								}
								goto l183
							l184:
								position, tokenIndex, depth = position184, tokenIndex184, depth184
							}
						}
					l179:
						depth--
						add(ruleWholeNum, position178)
					}
					depth--
					add(rulePegText, position177)
				}
				{
					add(ruleAction20, position)
				}
				depth--
				add(ruleInteger, position176)
			}
			return true
		l175:
			position, tokenIndex, depth = position175, tokenIndex175, depth175
			return false
		},
		/* 24 WholeNum <- <('0' / ('-'? [1-9] Digit*))> */
		nil,
		/* 25 Digit <- <[0-9]> */
		func() bool {
			position187, tokenIndex187, depth187 := position, tokenIndex, depth
			{
				position188 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l187
				}
				position++
				depth--
				add(ruleDigit, position188)
			}
			return true
		l187:
			position, tokenIndex, depth = position187, tokenIndex187, depth187
			return false
		},
		/* 26 Boolean <- <(Action21 <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> Action22 Action23)> */
		nil,
		/* 27 Func <- <(Action24 FuncArgs sp ('-' '>') sp (Block / Expr) Action25)> */
		nil,
		/* 28 FuncArgs <- <(Action26 '(' sp (LocalRef (sp ',' sp LocalRef)* sp)? ')' Action27)> */
		nil,
		/* 29 FuncApply <- <(Action28 Ref Tuple Action29)> */
		nil,
		/* 30 List <- <(Action30 '[' sp (Expr (sp ',' sp Expr)* sp)? ']' Action31)> */
		nil,
		/* 31 Tuple <- <(Action32 '(' sp (Expr (sp ',' sp Expr)* sp)? ')' Action33)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				{
					add(ruleAction32, position)
				}
				if buffer[position] != rune('(') {
					goto l194
				}
				position++
				if !_rules[rulesp]() {
					goto l194
				}
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if !_rules[ruleExpr]() {
						goto l197
					}
				l199:
					{
						position200, tokenIndex200, depth200 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l200
						}
						if buffer[position] != rune(',') {
							goto l200
						}
						position++
						if !_rules[rulesp]() {
							goto l200
						}
						if !_rules[ruleExpr]() {
							goto l200
						}
						goto l199
					l200:
						position, tokenIndex, depth = position200, tokenIndex200, depth200
					}
					if !_rules[rulesp]() {
						goto l197
					}
					goto l198
				l197:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
				}
			l198:
				if buffer[position] != rune(')') {
					goto l194
				}
				position++
				{
					add(ruleAction33, position)
				}
				depth--
				add(ruleTuple, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 32 Map <- <(Action34 '{' sp (Expr sp ':' sp Expr (sp ',' sp Expr sp ':' sp Expr)* sp)? '}' Action35)> */
		nil,
		/* 33 msp <- <(ws / comment)+> */
		nil,
		/* 34 sp <- <(ws / comment)*> */
		func() bool {
			{
				position205 := position
				depth++
			l206:
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					{
						position208, tokenIndex208, depth208 := position, tokenIndex, depth
						if !_rules[rulews]() {
							goto l209
						}
						goto l208
					l209:
						position, tokenIndex, depth = position208, tokenIndex208, depth208
						if !_rules[rulecomment]() {
							goto l207
						}
					}
				l208:
					goto l206
				l207:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
				}
				depth--
				add(rulesp, position205)
			}
			return true
		},
		/* 35 comment <- <('#' (!'\n' .)*)> */
		func() bool {
			position210, tokenIndex210, depth210 := position, tokenIndex, depth
			{
				position211 := position
				depth++
				if buffer[position] != rune('#') {
					goto l210
				}
				position++
			l212:
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					{
						position214, tokenIndex214, depth214 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l214
						}
						position++
						goto l213
					l214:
						position, tokenIndex, depth = position214, tokenIndex214, depth214
					}
					if !matchDot() {
						goto l213
					}
					goto l212
				l213:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
				}
				depth--
				add(rulecomment, position211)
			}
			return true
		l210:
			position, tokenIndex, depth = position210, tokenIndex210, depth210
			return false
		},
		/* 36 ws <- <((&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))> */
		func() bool {
			position215, tokenIndex215, depth215 := position, tokenIndex, depth
			{
				position216 := position
				depth++
				{
					switch buffer[position] {
					case '\r':
						if buffer[position] != rune('\r') {
							goto l215
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l215
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l215
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l215
						}
						position++
						break
					}
				}

				depth--
				add(rulews, position216)
			}
			return true
		l215:
			position, tokenIndex, depth = position215, tokenIndex215, depth215
			return false
		},
		/* 38 Action0 <- <{ p.Start(RIFT) }> */
		nil,
		/* 39 Action1 <- <{ p.End() }> */
		nil,
		/* 40 Action2 <- <{ p.Start(OP) }> */
		nil,
		/* 41 Action3 <- <{ p.End() }> */
		nil,
		/* 42 Action4 <- <{ p.Start(BINOP) }> */
		nil,
		nil,
		/* 44 Action5 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 45 Action6 <- <{ p.End() }> */
		nil,
		/* 46 Action7 <- <{ p.Start(ASSIGNMENT) }> */
		nil,
		/* 47 Action8 <- <{ p.End() }> */
		nil,
		/* 48 Action9 <- <{ p.Start(REF) }> */
		nil,
		/* 49 Action10 <- <{ p.End() }> */
		nil,
		/* 50 Action11 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 51 Action12 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 52 Action13 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 53 Action14 <- <{ p.Start(STRING) }> */
		nil,
		/* 54 Action15 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 55 Action16 <- <{ p.End() }> */
		nil,
		/* 56 Action17 <- <{ p.Start(NUM) }> */
		nil,
		/* 57 Action18 <- <{ p.End() }> */
		nil,
		/* 58 Action19 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 59 Action20 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 60 Action21 <- <{ p.Start(BOOL) }> */
		nil,
		/* 61 Action22 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 62 Action23 <- <{ p.End() }> */
		nil,
		/* 63 Action24 <- <{ p.Start(FUNC) }> */
		nil,
		/* 64 Action25 <- <{ p.End() }> */
		nil,
		/* 65 Action26 <- <{ p.Start(ARGS) }> */
		nil,
		/* 66 Action27 <- <{ p.End() }> */
		nil,
		/* 67 Action28 <- <{ p.Start(FUNCAPPLY) }> */
		nil,
		/* 68 Action29 <- <{ p.End() }> */
		nil,
		/* 69 Action30 <- <{ p.Start(LIST) }> */
		nil,
		/* 70 Action31 <- <{ p.End() }> */
		nil,
		/* 71 Action32 <- <{ p.Start(TUPLE) }> */
		nil,
		/* 72 Action33 <- <{ p.End() }> */
		nil,
		/* 73 Action34 <- <{ p.Start("map") }> */
		nil,
		/* 74 Action35 <- <{ p.End() }> */
		nil,
	}
	p.rules = _rules
}
