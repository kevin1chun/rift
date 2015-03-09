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
	ruleStatement
	ruleAssignment
	ruleIf
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
	ruleAction5
	ruleAction6
	rulePegText
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
	"Statement",
	"Assignment",
	"If",
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
	"Action5",
	"Action6",
	"PegText",
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
			p.Start(ASSIGNMENT)
		case ruleAction3:
			p.End()
		case ruleAction4:
			p.Start(IF)
		case ruleAction5:
			p.End()
		case ruleAction6:
			p.Start(REF)
		case ruleAction7:
			p.Emit(string(buffer[begin:end]))
		case ruleAction8:
			p.Emit(string(buffer[begin:end]))
		case ruleAction9:
			p.End()
		case ruleAction10:
			p.Start(REF)
		case ruleAction11:
			p.Emit(string(buffer[begin:end]))
		case ruleAction12:
			p.End()
		case ruleAction13:
			p.Start(STRING)
		case ruleAction14:
			p.Emit(string(buffer[begin:end]))
		case ruleAction15:
			p.End()
		case ruleAction16:
			p.Start(NUM)
		case ruleAction17:
			p.End()
		case ruleAction18:
			p.Emit(string(buffer[begin:end]))
		case ruleAction19:
			p.Emit(string(buffer[begin:end]))
		case ruleAction20:
			p.Start(BOOL)
		case ruleAction21:
			p.Emit(string(buffer[begin:end]))
		case ruleAction22:
			p.End()
		case ruleAction23:
			p.Start(FUNC)
		case ruleAction24:
			p.End()
		case ruleAction25:
			p.Start(ARGS)
		case ruleAction26:
			p.End()
		case ruleAction27:
			p.Start(FUNCAPPLY)
		case ruleAction28:
			p.End()
		case ruleAction29:
			p.Start(LIST)
		case ruleAction30:
			p.End()
		case ruleAction31:
			p.Start(TUPLE)
		case ruleAction32:
			p.End()
		case ruleAction33:
			p.Start("map")
		case ruleAction34:
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
						if !_rules[ruleGravitasse]() {
							goto l6
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
						position9 := position
						depth++
						{
							add(ruleAction0, position)
						}
						{
							position11, tokenIndex11, depth11 := position, tokenIndex, depth
							if !_rules[ruleGravitasse]() {
								goto l11
							}
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
		/* 1 Rift <- <(Action0 Gravitasse? LocalRef sp ('=' '>') sp Block Action1)> */
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
									position24, tokenIndex24, depth24 := position, tokenIndex, depth
									{
										position26 := position
										depth++
										{
											add(ruleAction2, position)
										}
										{
											position28, tokenIndex28, depth28 := position, tokenIndex, depth
											if !_rules[ruleGravitasse]() {
												goto l28
											}
											goto l29
										l28:
											position, tokenIndex, depth = position28, tokenIndex28, depth28
										}
									l29:
										if !_rules[ruleLocalRef]() {
											goto l25
										}
										if !_rules[rulesp]() {
											goto l25
										}
										if buffer[position] != rune('=') {
											goto l25
										}
										position++
										if !_rules[rulesp]() {
											goto l25
										}
										if !_rules[ruleExpr]() {
											goto l25
										}
										{
											add(ruleAction3, position)
										}
										depth--
										add(ruleAssignment, position26)
									}
									goto l24
								l25:
									position, tokenIndex, depth = position24, tokenIndex24, depth24
									{
										position31 := position
										depth++
										{
											add(ruleAction4, position)
										}
										if buffer[position] != rune('i') {
											goto l22
										}
										position++
										if buffer[position] != rune('f') {
											goto l22
										}
										position++
										if !_rules[rulesp]() {
											goto l22
										}
										if !_rules[ruleExpr]() {
											goto l22
										}
										if !_rules[rulesp]() {
											goto l22
										}
										if !_rules[ruleBlock]() {
											goto l22
										}
										{
											add(ruleAction5, position)
										}
										depth--
										add(ruleIf, position31)
									}
								}
							l24:
								depth--
								add(ruleStatement, position23)
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
						position34 := position
						depth++
						{
							position37, tokenIndex37, depth37 := position, tokenIndex, depth
							if !_rules[rulews]() {
								goto l38
							}
							goto l37
						l38:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if !_rules[rulecomment]() {
								goto l19
							}
						}
					l37:
					l35:
						{
							position36, tokenIndex36, depth36 := position, tokenIndex, depth
							{
								position39, tokenIndex39, depth39 := position, tokenIndex, depth
								if !_rules[rulews]() {
									goto l40
								}
								goto l39
							l40:
								position, tokenIndex, depth = position39, tokenIndex39, depth39
								if !_rules[rulecomment]() {
									goto l36
								}
							}
						l39:
							goto l35
						l36:
							position, tokenIndex, depth = position36, tokenIndex36, depth36
						}
						depth--
						add(rulemsp, position34)
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
		/* 3 Line <- <(Statement / Expr)> */
		nil,
		/* 4 Expr <- <(FuncApply / Value)> */
		func() bool {
			position42, tokenIndex42, depth42 := position, tokenIndex, depth
			{
				position43 := position
				depth++
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					{
						position46 := position
						depth++
						{
							add(ruleAction27, position)
						}
						if !_rules[ruleRef]() {
							goto l45
						}
						if !_rules[ruleTuple]() {
							goto l45
						}
						{
							add(ruleAction28, position)
						}
						depth--
						add(ruleFuncApply, position46)
					}
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					{
						position49 := position
						depth++
						{
							position50, tokenIndex50, depth50 := position, tokenIndex, depth
							if !_rules[ruleRef]() {
								goto l51
							}
							goto l50
						l51:
							position, tokenIndex, depth = position50, tokenIndex50, depth50
							{
								position52 := position
								depth++
								{
									position53, tokenIndex53, depth53 := position, tokenIndex, depth
									{
										position55 := position
										depth++
										{
											add(ruleAction23, position)
										}
										{
											position57 := position
											depth++
											{
												add(ruleAction25, position)
											}
											if buffer[position] != rune('(') {
												goto l54
											}
											position++
											if !_rules[rulesp]() {
												goto l54
											}
											{
												position59, tokenIndex59, depth59 := position, tokenIndex, depth
												if !_rules[ruleLocalRef]() {
													goto l59
												}
											l61:
												{
													position62, tokenIndex62, depth62 := position, tokenIndex, depth
													if !_rules[rulesp]() {
														goto l62
													}
													if buffer[position] != rune(',') {
														goto l62
													}
													position++
													if !_rules[rulesp]() {
														goto l62
													}
													if !_rules[ruleLocalRef]() {
														goto l62
													}
													goto l61
												l62:
													position, tokenIndex, depth = position62, tokenIndex62, depth62
												}
												if !_rules[rulesp]() {
													goto l59
												}
												goto l60
											l59:
												position, tokenIndex, depth = position59, tokenIndex59, depth59
											}
										l60:
											if buffer[position] != rune(')') {
												goto l54
											}
											position++
											{
												add(ruleAction26, position)
											}
											depth--
											add(ruleFuncArgs, position57)
										}
										if !_rules[rulesp]() {
											goto l54
										}
										if buffer[position] != rune('-') {
											goto l54
										}
										position++
										if buffer[position] != rune('>') {
											goto l54
										}
										position++
										if !_rules[rulesp]() {
											goto l54
										}
										{
											position64, tokenIndex64, depth64 := position, tokenIndex, depth
											if !_rules[ruleBlock]() {
												goto l65
											}
											goto l64
										l65:
											position, tokenIndex, depth = position64, tokenIndex64, depth64
											if !_rules[ruleExpr]() {
												goto l54
											}
										}
									l64:
										{
											add(ruleAction24, position)
										}
										depth--
										add(ruleFunc, position55)
									}
									goto l53
								l54:
									position, tokenIndex, depth = position53, tokenIndex53, depth53
									{
										position68 := position
										depth++
										{
											switch buffer[position] {
											case 'f', 't':
												{
													position70 := position
													depth++
													{
														add(ruleAction20, position)
													}
													{
														position72 := position
														depth++
														{
															position73, tokenIndex73, depth73 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l74
															}
															position++
															if buffer[position] != rune('r') {
																goto l74
															}
															position++
															if buffer[position] != rune('u') {
																goto l74
															}
															position++
															if buffer[position] != rune('e') {
																goto l74
															}
															position++
															goto l73
														l74:
															position, tokenIndex, depth = position73, tokenIndex73, depth73
															if buffer[position] != rune('f') {
																goto l67
															}
															position++
															if buffer[position] != rune('a') {
																goto l67
															}
															position++
															if buffer[position] != rune('l') {
																goto l67
															}
															position++
															if buffer[position] != rune('s') {
																goto l67
															}
															position++
															if buffer[position] != rune('e') {
																goto l67
															}
															position++
														}
													l73:
														depth--
														add(rulePegText, position72)
													}
													{
														add(ruleAction21, position)
													}
													{
														add(ruleAction22, position)
													}
													depth--
													add(ruleBoolean, position70)
												}
												break
											case '"':
												{
													position77 := position
													depth++
													{
														add(ruleAction13, position)
													}
													if buffer[position] != rune('"') {
														goto l67
													}
													position++
													{
														position79 := position
														depth++
													l80:
														{
															position81, tokenIndex81, depth81 := position, tokenIndex, depth
															{
																position82 := position
																depth++
																{
																	position83, tokenIndex83, depth83 := position, tokenIndex, depth
																	{
																		position85 := position
																		depth++
																		{
																			position86 := position
																			depth++
																			if buffer[position] != rune('\\') {
																				goto l84
																			}
																			position++
																			{
																				switch buffer[position] {
																				case 'v':
																					if buffer[position] != rune('v') {
																						goto l84
																					}
																					position++
																					break
																				case 't':
																					if buffer[position] != rune('t') {
																						goto l84
																					}
																					position++
																					break
																				case 'r':
																					if buffer[position] != rune('r') {
																						goto l84
																					}
																					position++
																					break
																				case 'n':
																					if buffer[position] != rune('n') {
																						goto l84
																					}
																					position++
																					break
																				case 'f':
																					if buffer[position] != rune('f') {
																						goto l84
																					}
																					position++
																					break
																				case 'b':
																					if buffer[position] != rune('b') {
																						goto l84
																					}
																					position++
																					break
																				case 'a':
																					if buffer[position] != rune('a') {
																						goto l84
																					}
																					position++
																					break
																				case '\\':
																					if buffer[position] != rune('\\') {
																						goto l84
																					}
																					position++
																					break
																				case '?':
																					if buffer[position] != rune('?') {
																						goto l84
																					}
																					position++
																					break
																				case '"':
																					if buffer[position] != rune('"') {
																						goto l84
																					}
																					position++
																					break
																				default:
																					if buffer[position] != rune('\'') {
																						goto l84
																					}
																					position++
																					break
																				}
																			}

																			depth--
																			add(ruleSimpleEsc, position86)
																		}
																		depth--
																		add(ruleStringEsc, position85)
																	}
																	goto l83
																l84:
																	position, tokenIndex, depth = position83, tokenIndex83, depth83
																	{
																		position88, tokenIndex88, depth88 := position, tokenIndex, depth
																		{
																			switch buffer[position] {
																			case '\\':
																				if buffer[position] != rune('\\') {
																					goto l88
																				}
																				position++
																				break
																			case '\n':
																				if buffer[position] != rune('\n') {
																					goto l88
																				}
																				position++
																				break
																			default:
																				if buffer[position] != rune('"') {
																					goto l88
																				}
																				position++
																				break
																			}
																		}

																		goto l81
																	l88:
																		position, tokenIndex, depth = position88, tokenIndex88, depth88
																	}
																	if !matchDot() {
																		goto l81
																	}
																}
															l83:
																depth--
																add(ruleStringChar, position82)
															}
															goto l80
														l81:
															position, tokenIndex, depth = position81, tokenIndex81, depth81
														}
														depth--
														add(rulePegText, position79)
													}
													if buffer[position] != rune('"') {
														goto l67
													}
													position++
													{
														add(ruleAction14, position)
													}
													{
														add(ruleAction15, position)
													}
													depth--
													add(ruleString, position77)
												}
												break
											default:
												{
													position92 := position
													depth++
													{
														add(ruleAction16, position)
													}
													{
														position94, tokenIndex94, depth94 := position, tokenIndex, depth
														{
															position96 := position
															depth++
															if !_rules[ruleDecimal]() {
																goto l95
															}
															{
																position97, tokenIndex97, depth97 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l98
																}
																position++
																goto l97
															l98:
																position, tokenIndex, depth = position97, tokenIndex97, depth97
																if buffer[position] != rune('E') {
																	goto l95
																}
																position++
															}
														l97:
															if !_rules[ruleInteger]() {
																goto l95
															}
															depth--
															add(ruleSciNum, position96)
														}
														goto l94
													l95:
														position, tokenIndex, depth = position94, tokenIndex94, depth94
														if !_rules[ruleDecimal]() {
															goto l99
														}
														goto l94
													l99:
														position, tokenIndex, depth = position94, tokenIndex94, depth94
														if !_rules[ruleInteger]() {
															goto l67
														}
													}
												l94:
													{
														add(ruleAction17, position)
													}
													depth--
													add(ruleNumeric, position92)
												}
												break
											}
										}

										depth--
										add(ruleScalar, position68)
									}
									goto l53
								l67:
									position, tokenIndex, depth = position53, tokenIndex53, depth53
									{
										position101 := position
										depth++
										{
											switch buffer[position] {
											case '{':
												{
													position103 := position
													depth++
													{
														add(ruleAction33, position)
													}
													if buffer[position] != rune('{') {
														goto l42
													}
													position++
													if !_rules[rulesp]() {
														goto l42
													}
													{
														position105, tokenIndex105, depth105 := position, tokenIndex, depth
														if !_rules[ruleExpr]() {
															goto l105
														}
														if !_rules[rulesp]() {
															goto l105
														}
														if buffer[position] != rune(':') {
															goto l105
														}
														position++
														if !_rules[rulesp]() {
															goto l105
														}
														if !_rules[ruleExpr]() {
															goto l105
														}
													l107:
														{
															position108, tokenIndex108, depth108 := position, tokenIndex, depth
															if !_rules[rulesp]() {
																goto l108
															}
															if buffer[position] != rune(',') {
																goto l108
															}
															position++
															if !_rules[rulesp]() {
																goto l108
															}
															if !_rules[ruleExpr]() {
																goto l108
															}
															if !_rules[rulesp]() {
																goto l108
															}
															if buffer[position] != rune(':') {
																goto l108
															}
															position++
															if !_rules[rulesp]() {
																goto l108
															}
															if !_rules[ruleExpr]() {
																goto l108
															}
															goto l107
														l108:
															position, tokenIndex, depth = position108, tokenIndex108, depth108
														}
														if !_rules[rulesp]() {
															goto l105
														}
														goto l106
													l105:
														position, tokenIndex, depth = position105, tokenIndex105, depth105
													}
												l106:
													if buffer[position] != rune('}') {
														goto l42
													}
													position++
													{
														add(ruleAction34, position)
													}
													depth--
													add(ruleMap, position103)
												}
												break
											case '(':
												if !_rules[ruleTuple]() {
													goto l42
												}
												break
											default:
												{
													position110 := position
													depth++
													{
														add(ruleAction29, position)
													}
													if buffer[position] != rune('[') {
														goto l42
													}
													position++
													if !_rules[rulesp]() {
														goto l42
													}
													{
														position112, tokenIndex112, depth112 := position, tokenIndex, depth
														if !_rules[ruleExpr]() {
															goto l112
														}
													l114:
														{
															position115, tokenIndex115, depth115 := position, tokenIndex, depth
															if !_rules[rulesp]() {
																goto l115
															}
															if buffer[position] != rune(',') {
																goto l115
															}
															position++
															if !_rules[rulesp]() {
																goto l115
															}
															if !_rules[ruleExpr]() {
																goto l115
															}
															goto l114
														l115:
															position, tokenIndex, depth = position115, tokenIndex115, depth115
														}
														if !_rules[rulesp]() {
															goto l112
														}
														goto l113
													l112:
														position, tokenIndex, depth = position112, tokenIndex112, depth112
													}
												l113:
													if buffer[position] != rune(']') {
														goto l42
													}
													position++
													{
														add(ruleAction30, position)
													}
													depth--
													add(ruleList, position110)
												}
												break
											}
										}

										depth--
										add(ruleVector, position101)
									}
								}
							l53:
								depth--
								add(ruleLiteral, position52)
							}
						}
					l50:
						depth--
						add(ruleValue, position49)
					}
				}
			l44:
				depth--
				add(ruleExpr, position43)
			}
			return true
		l42:
			position, tokenIndex, depth = position42, tokenIndex42, depth42
			return false
		},
		/* 5 Statement <- <(Assignment / If)> */
		nil,
		/* 6 Assignment <- <(Action2 Gravitasse? LocalRef sp '=' sp Expr Action3)> */
		nil,
		/* 7 If <- <(Action4 ('i' 'f') sp Expr sp Block Action5)> */
		nil,
		/* 8 Ref <- <(FullRef / LocalRef)> */
		func() bool {
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
							add(ruleAction6, position)
						}
						{
							position126 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l123
							}
						l127:
							{
								position128, tokenIndex128, depth128 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l128
								}
								goto l127
							l128:
								position, tokenIndex, depth = position128, tokenIndex128, depth128
							}
							depth--
							add(rulePegText, position126)
						}
						{
							add(ruleAction7, position)
						}
						if buffer[position] != rune(':') {
							goto l123
						}
						position++
						{
							position130 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l123
							}
						l131:
							{
								position132, tokenIndex132, depth132 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l132
								}
								goto l131
							l132:
								position, tokenIndex, depth = position132, tokenIndex132, depth132
							}
							depth--
							add(rulePegText, position130)
						}
						{
							add(ruleAction8, position)
						}
						{
							add(ruleAction9, position)
						}
						depth--
						add(ruleFullRef, position124)
					}
					goto l122
				l123:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
					if !_rules[ruleLocalRef]() {
						goto l120
					}
				}
			l122:
				depth--
				add(ruleRef, position121)
			}
			return true
		l120:
			position, tokenIndex, depth = position120, tokenIndex120, depth120
			return false
		},
		/* 9 FullRef <- <(Action6 <RefChar+> Action7 ':' <RefChar+> Action8 Action9)> */
		nil,
		/* 10 LocalRef <- <(Action10 <RefChar+> Action11 Action12)> */
		func() bool {
			position136, tokenIndex136, depth136 := position, tokenIndex, depth
			{
				position137 := position
				depth++
				{
					add(ruleAction10, position)
				}
				{
					position139 := position
					depth++
					if !_rules[ruleRefChar]() {
						goto l136
					}
				l140:
					{
						position141, tokenIndex141, depth141 := position, tokenIndex, depth
						if !_rules[ruleRefChar]() {
							goto l141
						}
						goto l140
					l141:
						position, tokenIndex, depth = position141, tokenIndex141, depth141
					}
					depth--
					add(rulePegText, position139)
				}
				{
					add(ruleAction11, position)
				}
				{
					add(ruleAction12, position)
				}
				depth--
				add(ruleLocalRef, position137)
			}
			return true
		l136:
			position, tokenIndex, depth = position136, tokenIndex136, depth136
			return false
		},
		/* 11 RefChar <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			{
				position145 := position
				depth++
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l144
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l144
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l144
						}
						position++
						break
					}
				}

				depth--
				add(ruleRefChar, position145)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 12 Value <- <(Ref / Literal)> */
		nil,
		/* 13 Literal <- <(Func / Scalar / Vector)> */
		nil,
		/* 14 Scalar <- <((&('f' | 't') Boolean) | (&('"') String) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric))> */
		nil,
		/* 15 Vector <- <((&('{') Map) | (&('(') Tuple) | (&('[') List))> */
		nil,
		/* 16 String <- <(Action13 '"' <StringChar*> '"' Action14 Action15)> */
		nil,
		/* 17 StringChar <- <(StringEsc / (!((&('\\') '\\') | (&('\n') '\n') | (&('"') '"')) .))> */
		nil,
		/* 18 StringEsc <- <SimpleEsc> */
		nil,
		/* 19 SimpleEsc <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 20 Numeric <- <(Action16 (SciNum / Decimal / Integer) Action17)> */
		nil,
		/* 21 SciNum <- <(Decimal ('e' / 'E') Integer)> */
		nil,
		/* 22 Decimal <- <(Integer '.' <Digit*> Action18)> */
		func() bool {
			position157, tokenIndex157, depth157 := position, tokenIndex, depth
			{
				position158 := position
				depth++
				if !_rules[ruleInteger]() {
					goto l157
				}
				if buffer[position] != rune('.') {
					goto l157
				}
				position++
				{
					position159 := position
					depth++
				l160:
					{
						position161, tokenIndex161, depth161 := position, tokenIndex, depth
						if !_rules[ruleDigit]() {
							goto l161
						}
						goto l160
					l161:
						position, tokenIndex, depth = position161, tokenIndex161, depth161
					}
					depth--
					add(rulePegText, position159)
				}
				{
					add(ruleAction18, position)
				}
				depth--
				add(ruleDecimal, position158)
			}
			return true
		l157:
			position, tokenIndex, depth = position157, tokenIndex157, depth157
			return false
		},
		/* 23 Integer <- <(<WholeNum> Action19)> */
		func() bool {
			position163, tokenIndex163, depth163 := position, tokenIndex, depth
			{
				position164 := position
				depth++
				{
					position165 := position
					depth++
					{
						position166 := position
						depth++
						{
							position167, tokenIndex167, depth167 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l168
							}
							position++
							goto l167
						l168:
							position, tokenIndex, depth = position167, tokenIndex167, depth167
							{
								position169, tokenIndex169, depth169 := position, tokenIndex, depth
								if buffer[position] != rune('-') {
									goto l169
								}
								position++
								goto l170
							l169:
								position, tokenIndex, depth = position169, tokenIndex169, depth169
							}
						l170:
							if c := buffer[position]; c < rune('1') || c > rune('9') {
								goto l163
							}
							position++
						l171:
							{
								position172, tokenIndex172, depth172 := position, tokenIndex, depth
								if !_rules[ruleDigit]() {
									goto l172
								}
								goto l171
							l172:
								position, tokenIndex, depth = position172, tokenIndex172, depth172
							}
						}
					l167:
						depth--
						add(ruleWholeNum, position166)
					}
					depth--
					add(rulePegText, position165)
				}
				{
					add(ruleAction19, position)
				}
				depth--
				add(ruleInteger, position164)
			}
			return true
		l163:
			position, tokenIndex, depth = position163, tokenIndex163, depth163
			return false
		},
		/* 24 WholeNum <- <('0' / ('-'? [1-9] Digit*))> */
		nil,
		/* 25 Digit <- <[0-9]> */
		func() bool {
			position175, tokenIndex175, depth175 := position, tokenIndex, depth
			{
				position176 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l175
				}
				position++
				depth--
				add(ruleDigit, position176)
			}
			return true
		l175:
			position, tokenIndex, depth = position175, tokenIndex175, depth175
			return false
		},
		/* 26 Boolean <- <(Action20 <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> Action21 Action22)> */
		nil,
		/* 27 Func <- <(Action23 FuncArgs sp ('-' '>') sp (Block / Expr) Action24)> */
		nil,
		/* 28 FuncArgs <- <(Action25 '(' sp (LocalRef (sp ',' sp LocalRef)* sp)? ')' Action26)> */
		nil,
		/* 29 FuncApply <- <(Action27 Ref Tuple Action28)> */
		nil,
		/* 30 List <- <(Action29 '[' sp (Expr (sp ',' sp Expr)* sp)? ']' Action30)> */
		nil,
		/* 31 Tuple <- <(Action31 '(' sp (Expr (sp ',' sp Expr)* sp)? ')' Action32)> */
		func() bool {
			position182, tokenIndex182, depth182 := position, tokenIndex, depth
			{
				position183 := position
				depth++
				{
					add(ruleAction31, position)
				}
				if buffer[position] != rune('(') {
					goto l182
				}
				position++
				if !_rules[rulesp]() {
					goto l182
				}
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					if !_rules[ruleExpr]() {
						goto l185
					}
				l187:
					{
						position188, tokenIndex188, depth188 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l188
						}
						if buffer[position] != rune(',') {
							goto l188
						}
						position++
						if !_rules[rulesp]() {
							goto l188
						}
						if !_rules[ruleExpr]() {
							goto l188
						}
						goto l187
					l188:
						position, tokenIndex, depth = position188, tokenIndex188, depth188
					}
					if !_rules[rulesp]() {
						goto l185
					}
					goto l186
				l185:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
				}
			l186:
				if buffer[position] != rune(')') {
					goto l182
				}
				position++
				{
					add(ruleAction32, position)
				}
				depth--
				add(ruleTuple, position183)
			}
			return true
		l182:
			position, tokenIndex, depth = position182, tokenIndex182, depth182
			return false
		},
		/* 32 Map <- <(Action33 '{' sp (Expr sp ':' sp Expr (sp ',' sp Expr sp ':' sp Expr)* sp)? '}' Action34)> */
		nil,
		/* 33 Gravitasse <- <'@'> */
		func() bool {
			position191, tokenIndex191, depth191 := position, tokenIndex, depth
			{
				position192 := position
				depth++
				if buffer[position] != rune('@') {
					goto l191
				}
				position++
				depth--
				add(ruleGravitasse, position192)
			}
			return true
		l191:
			position, tokenIndex, depth = position191, tokenIndex191, depth191
			return false
		},
		/* 34 msp <- <(ws / comment)+> */
		nil,
		/* 35 sp <- <(ws / comment)*> */
		func() bool {
			{
				position195 := position
				depth++
			l196:
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					{
						position198, tokenIndex198, depth198 := position, tokenIndex, depth
						if !_rules[rulews]() {
							goto l199
						}
						goto l198
					l199:
						position, tokenIndex, depth = position198, tokenIndex198, depth198
						if !_rules[rulecomment]() {
							goto l197
						}
					}
				l198:
					goto l196
				l197:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
				}
				depth--
				add(rulesp, position195)
			}
			return true
		},
		/* 36 comment <- <('#' (!'\n' .)*)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				if buffer[position] != rune('#') {
					goto l200
				}
				position++
			l202:
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					{
						position204, tokenIndex204, depth204 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l204
						}
						position++
						goto l203
					l204:
						position, tokenIndex, depth = position204, tokenIndex204, depth204
					}
					if !matchDot() {
						goto l203
					}
					goto l202
				l203:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
				}
				depth--
				add(rulecomment, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 37 ws <- <((&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))> */
		func() bool {
			position205, tokenIndex205, depth205 := position, tokenIndex, depth
			{
				position206 := position
				depth++
				{
					switch buffer[position] {
					case '\r':
						if buffer[position] != rune('\r') {
							goto l205
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l205
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l205
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l205
						}
						position++
						break
					}
				}

				depth--
				add(rulews, position206)
			}
			return true
		l205:
			position, tokenIndex, depth = position205, tokenIndex205, depth205
			return false
		},
		/* 39 Action0 <- <{ p.Start(RIFT) }> */
		nil,
		/* 40 Action1 <- <{ p.End() }> */
		nil,
		/* 41 Action2 <- <{ p.Start(ASSIGNMENT) }> */
		nil,
		/* 42 Action3 <- <{ p.End() }> */
		nil,
		/* 43 Action4 <- <{ p.Start(IF) }> */
		nil,
		/* 44 Action5 <- <{ p.End() }> */
		nil,
		/* 45 Action6 <- <{ p.Start(REF) }> */
		nil,
		nil,
		/* 47 Action7 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 48 Action8 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 49 Action9 <- <{ p.End() }> */
		nil,
		/* 50 Action10 <- <{ p.Start(REF) }> */
		nil,
		/* 51 Action11 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 52 Action12 <- <{ p.End() }> */
		nil,
		/* 53 Action13 <- <{ p.Start(STRING) }> */
		nil,
		/* 54 Action14 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 55 Action15 <- <{ p.End() }> */
		nil,
		/* 56 Action16 <- <{ p.Start(NUM) }> */
		nil,
		/* 57 Action17 <- <{ p.End() }> */
		nil,
		/* 58 Action18 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 59 Action19 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 60 Action20 <- <{ p.Start(BOOL) }> */
		nil,
		/* 61 Action21 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 62 Action22 <- <{ p.End() }> */
		nil,
		/* 63 Action23 <- <{ p.Start(FUNC) }> */
		nil,
		/* 64 Action24 <- <{ p.End() }> */
		nil,
		/* 65 Action25 <- <{ p.Start(ARGS) }> */
		nil,
		/* 66 Action26 <- <{ p.End() }> */
		nil,
		/* 67 Action27 <- <{ p.Start(FUNCAPPLY) }> */
		nil,
		/* 68 Action28 <- <{ p.End() }> */
		nil,
		/* 69 Action29 <- <{ p.Start(LIST) }> */
		nil,
		/* 70 Action30 <- <{ p.End() }> */
		nil,
		/* 71 Action31 <- <{ p.Start(TUPLE) }> */
		nil,
		/* 72 Action32 <- <{ p.End() }> */
		nil,
		/* 73 Action33 <- <{ p.Start("map") }> */
		nil,
		/* 74 Action34 <- <{ p.End() }> */
		nil,
	}
	p.rules = _rules
}
