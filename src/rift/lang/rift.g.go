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
						position7 := position
						depth++
						{
							add(ruleAction0, position)
						}
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
						add(ruleRift, position7)
					}
					if !_rules[rulesp]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				{
					position10, tokenIndex10, depth10 := position, tokenIndex, depth
					if !matchDot() {
						goto l10
					}
					goto l0
				l10:
					position, tokenIndex, depth = position10, tokenIndex10, depth10
				}
				depth--
				add(ruleSource, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Rift <- <(Action0 LocalRef sp ('=' '>') sp Block Action1)> */
		nil,
		/* 2 Block <- <('{' sp (Line msp)* '}')> */
		func() bool {
			position12, tokenIndex12, depth12 := position, tokenIndex, depth
			{
				position13 := position
				depth++
				if buffer[position] != rune('{') {
					goto l12
				}
				position++
				if !_rules[rulesp]() {
					goto l12
				}
			l14:
				{
					position15, tokenIndex15, depth15 := position, tokenIndex, depth
					{
						position16 := position
						depth++
						{
							position17, tokenIndex17, depth17 := position, tokenIndex, depth
							{
								position19 := position
								depth++
								{
									add(ruleAction7, position)
								}
								if !_rules[ruleLocalRef]() {
									goto l18
								}
								if !_rules[rulesp]() {
									goto l18
								}
								if buffer[position] != rune('=') {
									goto l18
								}
								position++
								if !_rules[rulesp]() {
									goto l18
								}
								if !_rules[ruleExpr]() {
									goto l18
								}
								{
									add(ruleAction8, position)
								}
								depth--
								add(ruleAssignment, position19)
							}
							goto l17
						l18:
							position, tokenIndex, depth = position17, tokenIndex17, depth17
							if !_rules[ruleExpr]() {
								goto l15
							}
						}
					l17:
						depth--
						add(ruleLine, position16)
					}
					{
						position22 := position
						depth++
						{
							position25, tokenIndex25, depth25 := position, tokenIndex, depth
							if !_rules[rulews]() {
								goto l26
							}
							goto l25
						l26:
							position, tokenIndex, depth = position25, tokenIndex25, depth25
							if !_rules[rulecomment]() {
								goto l15
							}
						}
					l25:
					l23:
						{
							position24, tokenIndex24, depth24 := position, tokenIndex, depth
							{
								position27, tokenIndex27, depth27 := position, tokenIndex, depth
								if !_rules[rulews]() {
									goto l28
								}
								goto l27
							l28:
								position, tokenIndex, depth = position27, tokenIndex27, depth27
								if !_rules[rulecomment]() {
									goto l24
								}
							}
						l27:
							goto l23
						l24:
							position, tokenIndex, depth = position24, tokenIndex24, depth24
						}
						depth--
						add(rulemsp, position22)
					}
					goto l14
				l15:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
				}
				if buffer[position] != rune('}') {
					goto l12
				}
				position++
				depth--
				add(ruleBlock, position13)
			}
			return true
		l12:
			position, tokenIndex, depth = position12, tokenIndex12, depth12
			return false
		},
		/* 3 Line <- <(Assignment / Expr)> */
		nil,
		/* 4 Expr <- <(Op / FuncApply / Value)> */
		func() bool {
			position30, tokenIndex30, depth30 := position, tokenIndex, depth
			{
				position31 := position
				depth++
				{
					position32, tokenIndex32, depth32 := position, tokenIndex, depth
					{
						position34 := position
						depth++
						{
							add(ruleAction2, position)
						}
						if !_rules[ruleValue]() {
							goto l33
						}
						if !_rules[rulesp]() {
							goto l33
						}
						{
							position38 := position
							depth++
							{
								add(ruleAction4, position)
							}
							{
								position40 := position
								depth++
								{
									position41, tokenIndex41, depth41 := position, tokenIndex, depth
									if buffer[position] != rune('*') {
										goto l42
									}
									position++
									goto l41
								l42:
									position, tokenIndex, depth = position41, tokenIndex41, depth41
									{
										switch buffer[position] {
										case '%':
											if buffer[position] != rune('%') {
												goto l33
											}
											position++
											break
										case '*':
											if buffer[position] != rune('*') {
												goto l33
											}
											position++
											if buffer[position] != rune('*') {
												goto l33
											}
											position++
											break
										case '/':
											if buffer[position] != rune('/') {
												goto l33
											}
											position++
											break
										case '-':
											if buffer[position] != rune('-') {
												goto l33
											}
											position++
											break
										default:
											if buffer[position] != rune('+') {
												goto l33
											}
											position++
											break
										}
									}

								}
							l41:
								depth--
								add(rulePegText, position40)
							}
							{
								add(ruleAction5, position)
							}
							{
								add(ruleAction6, position)
							}
							depth--
							add(ruleBinaryOp, position38)
						}
						if !_rules[rulesp]() {
							goto l33
						}
						if !_rules[ruleValue]() {
							goto l33
						}
					l36:
						{
							position37, tokenIndex37, depth37 := position, tokenIndex, depth
							if !_rules[rulesp]() {
								goto l37
							}
							{
								position46 := position
								depth++
								{
									add(ruleAction4, position)
								}
								{
									position48 := position
									depth++
									{
										position49, tokenIndex49, depth49 := position, tokenIndex, depth
										if buffer[position] != rune('*') {
											goto l50
										}
										position++
										goto l49
									l50:
										position, tokenIndex, depth = position49, tokenIndex49, depth49
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
								l49:
									depth--
									add(rulePegText, position48)
								}
								{
									add(ruleAction5, position)
								}
								{
									add(ruleAction6, position)
								}
								depth--
								add(ruleBinaryOp, position46)
							}
							if !_rules[rulesp]() {
								goto l37
							}
							if !_rules[ruleValue]() {
								goto l37
							}
							goto l36
						l37:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
						}
						{
							add(ruleAction3, position)
						}
						depth--
						add(ruleOp, position34)
					}
					goto l32
				l33:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					{
						position56 := position
						depth++
						{
							add(ruleAction28, position)
						}
						if !_rules[ruleRef]() {
							goto l55
						}
						if !_rules[ruleTuple]() {
							goto l55
						}
						{
							add(ruleAction29, position)
						}
						depth--
						add(ruleFuncApply, position56)
					}
					goto l32
				l55:
					position, tokenIndex, depth = position32, tokenIndex32, depth32
					if !_rules[ruleValue]() {
						goto l30
					}
				}
			l32:
				depth--
				add(ruleExpr, position31)
			}
			return true
		l30:
			position, tokenIndex, depth = position30, tokenIndex30, depth30
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
			position62, tokenIndex62, depth62 := position, tokenIndex, depth
			{
				position63 := position
				depth++
				{
					add(ruleAction9, position)
				}
				{
					position65, tokenIndex65, depth65 := position, tokenIndex, depth
					{
						position67 := position
						depth++
						{
							position68 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l66
							}
						l69:
							{
								position70, tokenIndex70, depth70 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l70
								}
								goto l69
							l70:
								position, tokenIndex, depth = position70, tokenIndex70, depth70
							}
							depth--
							add(rulePegText, position68)
						}
						{
							add(ruleAction11, position)
						}
						if buffer[position] != rune(':') {
							goto l66
						}
						position++
						{
							position72 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l66
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
							add(ruleAction12, position)
						}
						depth--
						add(ruleFullRef, position67)
					}
					goto l65
				l66:
					position, tokenIndex, depth = position65, tokenIndex65, depth65
					if !_rules[ruleLocalRef]() {
						goto l62
					}
				}
			l65:
				{
					add(ruleAction10, position)
				}
				depth--
				add(ruleRef, position63)
			}
			return true
		l62:
			position, tokenIndex, depth = position62, tokenIndex62, depth62
			return false
		},
		/* 9 FullRef <- <(<RefChar+> Action11 ':' <RefChar+> Action12)> */
		nil,
		/* 10 LocalRef <- <(<RefChar+> Action13)> */
		func() bool {
			position78, tokenIndex78, depth78 := position, tokenIndex, depth
			{
				position79 := position
				depth++
				{
					position80 := position
					depth++
					if !_rules[ruleRefChar]() {
						goto l78
					}
				l81:
					{
						position82, tokenIndex82, depth82 := position, tokenIndex, depth
						if !_rules[ruleRefChar]() {
							goto l82
						}
						goto l81
					l82:
						position, tokenIndex, depth = position82, tokenIndex82, depth82
					}
					depth--
					add(rulePegText, position80)
				}
				{
					add(ruleAction13, position)
				}
				depth--
				add(ruleLocalRef, position79)
			}
			return true
		l78:
			position, tokenIndex, depth = position78, tokenIndex78, depth78
			return false
		},
		/* 11 RefChar <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position84, tokenIndex84, depth84 := position, tokenIndex, depth
			{
				position85 := position
				depth++
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l84
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l84
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l84
						}
						position++
						break
					}
				}

				depth--
				add(ruleRefChar, position85)
			}
			return true
		l84:
			position, tokenIndex, depth = position84, tokenIndex84, depth84
			return false
		},
		/* 12 Value <- <(Ref / Literal)> */
		func() bool {
			position87, tokenIndex87, depth87 := position, tokenIndex, depth
			{
				position88 := position
				depth++
				{
					position89, tokenIndex89, depth89 := position, tokenIndex, depth
					if !_rules[ruleRef]() {
						goto l90
					}
					goto l89
				l90:
					position, tokenIndex, depth = position89, tokenIndex89, depth89
					{
						position91 := position
						depth++
						{
							position92, tokenIndex92, depth92 := position, tokenIndex, depth
							{
								position94 := position
								depth++
								{
									add(ruleAction24, position)
								}
								{
									position96 := position
									depth++
									{
										add(ruleAction26, position)
									}
									if buffer[position] != rune('(') {
										goto l93
									}
									position++
									if !_rules[rulesp]() {
										goto l93
									}
									{
										position98, tokenIndex98, depth98 := position, tokenIndex, depth
										if !_rules[ruleLocalRef]() {
											goto l98
										}
									l100:
										{
											position101, tokenIndex101, depth101 := position, tokenIndex, depth
											if !_rules[rulesp]() {
												goto l101
											}
											if buffer[position] != rune(',') {
												goto l101
											}
											position++
											if !_rules[rulesp]() {
												goto l101
											}
											if !_rules[ruleLocalRef]() {
												goto l101
											}
											goto l100
										l101:
											position, tokenIndex, depth = position101, tokenIndex101, depth101
										}
										if !_rules[rulesp]() {
											goto l98
										}
										goto l99
									l98:
										position, tokenIndex, depth = position98, tokenIndex98, depth98
									}
								l99:
									if buffer[position] != rune(')') {
										goto l93
									}
									position++
									{
										add(ruleAction27, position)
									}
									depth--
									add(ruleFuncArgs, position96)
								}
								if !_rules[rulesp]() {
									goto l93
								}
								if buffer[position] != rune('-') {
									goto l93
								}
								position++
								if buffer[position] != rune('>') {
									goto l93
								}
								position++
								if !_rules[rulesp]() {
									goto l93
								}
								{
									position103, tokenIndex103, depth103 := position, tokenIndex, depth
									if !_rules[ruleBlock]() {
										goto l104
									}
									goto l103
								l104:
									position, tokenIndex, depth = position103, tokenIndex103, depth103
									if !_rules[ruleExpr]() {
										goto l93
									}
								}
							l103:
								{
									add(ruleAction25, position)
								}
								depth--
								add(ruleFunc, position94)
							}
							goto l92
						l93:
							position, tokenIndex, depth = position92, tokenIndex92, depth92
							{
								position107 := position
								depth++
								{
									switch buffer[position] {
									case 'f', 't':
										{
											position109 := position
											depth++
											{
												add(ruleAction21, position)
											}
											{
												position111 := position
												depth++
												{
													position112, tokenIndex112, depth112 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l113
													}
													position++
													if buffer[position] != rune('r') {
														goto l113
													}
													position++
													if buffer[position] != rune('u') {
														goto l113
													}
													position++
													if buffer[position] != rune('e') {
														goto l113
													}
													position++
													goto l112
												l113:
													position, tokenIndex, depth = position112, tokenIndex112, depth112
													if buffer[position] != rune('f') {
														goto l106
													}
													position++
													if buffer[position] != rune('a') {
														goto l106
													}
													position++
													if buffer[position] != rune('l') {
														goto l106
													}
													position++
													if buffer[position] != rune('s') {
														goto l106
													}
													position++
													if buffer[position] != rune('e') {
														goto l106
													}
													position++
												}
											l112:
												depth--
												add(rulePegText, position111)
											}
											{
												add(ruleAction22, position)
											}
											{
												add(ruleAction23, position)
											}
											depth--
											add(ruleBoolean, position109)
										}
										break
									case '"':
										{
											position116 := position
											depth++
											{
												add(ruleAction14, position)
											}
											if buffer[position] != rune('"') {
												goto l106
											}
											position++
											{
												position118 := position
												depth++
											l119:
												{
													position120, tokenIndex120, depth120 := position, tokenIndex, depth
													{
														position121 := position
														depth++
														{
															position122, tokenIndex122, depth122 := position, tokenIndex, depth
															{
																position124 := position
																depth++
																{
																	position125 := position
																	depth++
																	if buffer[position] != rune('\\') {
																		goto l123
																	}
																	position++
																	{
																		switch buffer[position] {
																		case 'v':
																			if buffer[position] != rune('v') {
																				goto l123
																			}
																			position++
																			break
																		case 't':
																			if buffer[position] != rune('t') {
																				goto l123
																			}
																			position++
																			break
																		case 'r':
																			if buffer[position] != rune('r') {
																				goto l123
																			}
																			position++
																			break
																		case 'n':
																			if buffer[position] != rune('n') {
																				goto l123
																			}
																			position++
																			break
																		case 'f':
																			if buffer[position] != rune('f') {
																				goto l123
																			}
																			position++
																			break
																		case 'b':
																			if buffer[position] != rune('b') {
																				goto l123
																			}
																			position++
																			break
																		case 'a':
																			if buffer[position] != rune('a') {
																				goto l123
																			}
																			position++
																			break
																		case '\\':
																			if buffer[position] != rune('\\') {
																				goto l123
																			}
																			position++
																			break
																		case '?':
																			if buffer[position] != rune('?') {
																				goto l123
																			}
																			position++
																			break
																		case '"':
																			if buffer[position] != rune('"') {
																				goto l123
																			}
																			position++
																			break
																		default:
																			if buffer[position] != rune('\'') {
																				goto l123
																			}
																			position++
																			break
																		}
																	}

																	depth--
																	add(ruleSimpleEsc, position125)
																}
																depth--
																add(ruleStringEsc, position124)
															}
															goto l122
														l123:
															position, tokenIndex, depth = position122, tokenIndex122, depth122
															{
																position127, tokenIndex127, depth127 := position, tokenIndex, depth
																{
																	switch buffer[position] {
																	case '\\':
																		if buffer[position] != rune('\\') {
																			goto l127
																		}
																		position++
																		break
																	case '\n':
																		if buffer[position] != rune('\n') {
																			goto l127
																		}
																		position++
																		break
																	default:
																		if buffer[position] != rune('"') {
																			goto l127
																		}
																		position++
																		break
																	}
																}

																goto l120
															l127:
																position, tokenIndex, depth = position127, tokenIndex127, depth127
															}
															if !matchDot() {
																goto l120
															}
														}
													l122:
														depth--
														add(ruleStringChar, position121)
													}
													goto l119
												l120:
													position, tokenIndex, depth = position120, tokenIndex120, depth120
												}
												depth--
												add(rulePegText, position118)
											}
											if buffer[position] != rune('"') {
												goto l106
											}
											position++
											{
												add(ruleAction15, position)
											}
											{
												add(ruleAction16, position)
											}
											depth--
											add(ruleString, position116)
										}
										break
									default:
										{
											position131 := position
											depth++
											{
												add(ruleAction17, position)
											}
											{
												position133, tokenIndex133, depth133 := position, tokenIndex, depth
												{
													position135 := position
													depth++
													if !_rules[ruleDecimal]() {
														goto l134
													}
													{
														position136, tokenIndex136, depth136 := position, tokenIndex, depth
														if buffer[position] != rune('e') {
															goto l137
														}
														position++
														goto l136
													l137:
														position, tokenIndex, depth = position136, tokenIndex136, depth136
														if buffer[position] != rune('E') {
															goto l134
														}
														position++
													}
												l136:
													if !_rules[ruleInteger]() {
														goto l134
													}
													depth--
													add(ruleSciNum, position135)
												}
												goto l133
											l134:
												position, tokenIndex, depth = position133, tokenIndex133, depth133
												if !_rules[ruleDecimal]() {
													goto l138
												}
												goto l133
											l138:
												position, tokenIndex, depth = position133, tokenIndex133, depth133
												if !_rules[ruleInteger]() {
													goto l106
												}
											}
										l133:
											{
												add(ruleAction18, position)
											}
											depth--
											add(ruleNumeric, position131)
										}
										break
									}
								}

								depth--
								add(ruleScalar, position107)
							}
							goto l92
						l106:
							position, tokenIndex, depth = position92, tokenIndex92, depth92
							{
								position140 := position
								depth++
								{
									switch buffer[position] {
									case '{':
										{
											position142 := position
											depth++
											{
												add(ruleAction34, position)
											}
											if buffer[position] != rune('{') {
												goto l87
											}
											position++
											if !_rules[rulesp]() {
												goto l87
											}
											{
												position144, tokenIndex144, depth144 := position, tokenIndex, depth
												if !_rules[ruleExpr]() {
													goto l144
												}
												if !_rules[rulesp]() {
													goto l144
												}
												if buffer[position] != rune(':') {
													goto l144
												}
												position++
												if !_rules[rulesp]() {
													goto l144
												}
												if !_rules[ruleExpr]() {
													goto l144
												}
											l146:
												{
													position147, tokenIndex147, depth147 := position, tokenIndex, depth
													if !_rules[rulesp]() {
														goto l147
													}
													if buffer[position] != rune(',') {
														goto l147
													}
													position++
													if !_rules[rulesp]() {
														goto l147
													}
													if !_rules[ruleExpr]() {
														goto l147
													}
													if !_rules[rulesp]() {
														goto l147
													}
													if buffer[position] != rune(':') {
														goto l147
													}
													position++
													if !_rules[rulesp]() {
														goto l147
													}
													if !_rules[ruleExpr]() {
														goto l147
													}
													goto l146
												l147:
													position, tokenIndex, depth = position147, tokenIndex147, depth147
												}
												if !_rules[rulesp]() {
													goto l144
												}
												goto l145
											l144:
												position, tokenIndex, depth = position144, tokenIndex144, depth144
											}
										l145:
											if buffer[position] != rune('}') {
												goto l87
											}
											position++
											{
												add(ruleAction35, position)
											}
											depth--
											add(ruleMap, position142)
										}
										break
									case '(':
										if !_rules[ruleTuple]() {
											goto l87
										}
										break
									default:
										{
											position149 := position
											depth++
											{
												add(ruleAction30, position)
											}
											if buffer[position] != rune('[') {
												goto l87
											}
											position++
											if !_rules[rulesp]() {
												goto l87
											}
											{
												position151, tokenIndex151, depth151 := position, tokenIndex, depth
												if !_rules[ruleExpr]() {
													goto l151
												}
											l153:
												{
													position154, tokenIndex154, depth154 := position, tokenIndex, depth
													if !_rules[rulesp]() {
														goto l154
													}
													if buffer[position] != rune(',') {
														goto l154
													}
													position++
													if !_rules[rulesp]() {
														goto l154
													}
													if !_rules[ruleExpr]() {
														goto l154
													}
													goto l153
												l154:
													position, tokenIndex, depth = position154, tokenIndex154, depth154
												}
												if !_rules[rulesp]() {
													goto l151
												}
												goto l152
											l151:
												position, tokenIndex, depth = position151, tokenIndex151, depth151
											}
										l152:
											if buffer[position] != rune(']') {
												goto l87
											}
											position++
											{
												add(ruleAction31, position)
											}
											depth--
											add(ruleList, position149)
										}
										break
									}
								}

								depth--
								add(ruleVector, position140)
							}
						}
					l92:
						depth--
						add(ruleLiteral, position91)
					}
				}
			l89:
				depth--
				add(ruleValue, position88)
			}
			return true
		l87:
			position, tokenIndex, depth = position87, tokenIndex87, depth87
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
			position165, tokenIndex165, depth165 := position, tokenIndex, depth
			{
				position166 := position
				depth++
				if !_rules[ruleInteger]() {
					goto l165
				}
				if buffer[position] != rune('.') {
					goto l165
				}
				position++
				{
					position167 := position
					depth++
				l168:
					{
						position169, tokenIndex169, depth169 := position, tokenIndex, depth
						if !_rules[ruleDigit]() {
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
					add(ruleAction19, position)
				}
				depth--
				add(ruleDecimal, position166)
			}
			return true
		l165:
			position, tokenIndex, depth = position165, tokenIndex165, depth165
			return false
		},
		/* 23 Integer <- <(<WholeNum> Action20)> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
				{
					position173 := position
					depth++
					{
						position174 := position
						depth++
						{
							position175, tokenIndex175, depth175 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l176
							}
							position++
							goto l175
						l176:
							position, tokenIndex, depth = position175, tokenIndex175, depth175
							{
								position177, tokenIndex177, depth177 := position, tokenIndex, depth
								if buffer[position] != rune('-') {
									goto l177
								}
								position++
								goto l178
							l177:
								position, tokenIndex, depth = position177, tokenIndex177, depth177
							}
						l178:
							if c := buffer[position]; c < rune('1') || c > rune('9') {
								goto l171
							}
							position++
						l179:
							{
								position180, tokenIndex180, depth180 := position, tokenIndex, depth
								if !_rules[ruleDigit]() {
									goto l180
								}
								goto l179
							l180:
								position, tokenIndex, depth = position180, tokenIndex180, depth180
							}
						}
					l175:
						depth--
						add(ruleWholeNum, position174)
					}
					depth--
					add(rulePegText, position173)
				}
				{
					add(ruleAction20, position)
				}
				depth--
				add(ruleInteger, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 24 WholeNum <- <('0' / ('-'? [1-9] Digit*))> */
		nil,
		/* 25 Digit <- <[0-9]> */
		func() bool {
			position183, tokenIndex183, depth183 := position, tokenIndex, depth
			{
				position184 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l183
				}
				position++
				depth--
				add(ruleDigit, position184)
			}
			return true
		l183:
			position, tokenIndex, depth = position183, tokenIndex183, depth183
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
			position190, tokenIndex190, depth190 := position, tokenIndex, depth
			{
				position191 := position
				depth++
				{
					add(ruleAction32, position)
				}
				if buffer[position] != rune('(') {
					goto l190
				}
				position++
				if !_rules[rulesp]() {
					goto l190
				}
				{
					position193, tokenIndex193, depth193 := position, tokenIndex, depth
					if !_rules[ruleExpr]() {
						goto l193
					}
				l195:
					{
						position196, tokenIndex196, depth196 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l196
						}
						if buffer[position] != rune(',') {
							goto l196
						}
						position++
						if !_rules[rulesp]() {
							goto l196
						}
						if !_rules[ruleExpr]() {
							goto l196
						}
						goto l195
					l196:
						position, tokenIndex, depth = position196, tokenIndex196, depth196
					}
					if !_rules[rulesp]() {
						goto l193
					}
					goto l194
				l193:
					position, tokenIndex, depth = position193, tokenIndex193, depth193
				}
			l194:
				if buffer[position] != rune(')') {
					goto l190
				}
				position++
				{
					add(ruleAction33, position)
				}
				depth--
				add(ruleTuple, position191)
			}
			return true
		l190:
			position, tokenIndex, depth = position190, tokenIndex190, depth190
			return false
		},
		/* 32 Map <- <(Action34 '{' sp (Expr sp ':' sp Expr (sp ',' sp Expr sp ':' sp Expr)* sp)? '}' Action35)> */
		nil,
		/* 33 msp <- <(ws / comment)+> */
		nil,
		/* 34 sp <- <(ws / comment)*> */
		func() bool {
			{
				position201 := position
				depth++
			l202:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					{
						position204, tokenIndex204, depth204 := position, tokenIndex, depth
						if !_rules[rulews]() {
							goto l205
						}
						goto l204
					l205:
						position, tokenIndex, depth = position204, tokenIndex204, depth204
						if !_rules[rulecomment]() {
							goto l203
						}
					}
				l204:
					goto l202
				l203:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
				}
				depth--
				add(rulesp, position201)
			}
			return true
		},
		/* 35 comment <- <('#' (!'\n' .)*)> */
		func() bool {
			position206, tokenIndex206, depth206 := position, tokenIndex, depth
			{
				position207 := position
				depth++
				if buffer[position] != rune('#') {
					goto l206
				}
				position++
			l208:
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					{
						position210, tokenIndex210, depth210 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l210
						}
						position++
						goto l209
					l210:
						position, tokenIndex, depth = position210, tokenIndex210, depth210
					}
					if !matchDot() {
						goto l209
					}
					goto l208
				l209:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
				}
				depth--
				add(rulecomment, position207)
			}
			return true
		l206:
			position, tokenIndex, depth = position206, tokenIndex206, depth206
			return false
		},
		/* 36 ws <- <((&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))> */
		func() bool {
			position211, tokenIndex211, depth211 := position, tokenIndex, depth
			{
				position212 := position
				depth++
				{
					switch buffer[position] {
					case '\r':
						if buffer[position] != rune('\r') {
							goto l211
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l211
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l211
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l211
						}
						position++
						break
					}
				}

				depth--
				add(rulews, position212)
			}
			return true
		l211:
			position, tokenIndex, depth = position211, tokenIndex211, depth211
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
