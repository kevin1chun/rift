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
	rules  [82]func() bool
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
			p.Start(REF)
		case ruleAction12:
			p.Emit(string(buffer[begin:end]))
		case ruleAction13:
			p.Emit(string(buffer[begin:end]))
		case ruleAction14:
			p.End()
		case ruleAction15:
			p.Start(REF)
		case ruleAction16:
			p.Emit(string(buffer[begin:end]))
		case ruleAction17:
			p.End()
		case ruleAction18:
			p.Start(STRING)
		case ruleAction19:
			p.Emit(string(buffer[begin:end]))
		case ruleAction20:
			p.End()
		case ruleAction21:
			p.Start(NUM)
		case ruleAction22:
			p.End()
		case ruleAction23:
			p.Emit(string(buffer[begin:end]))
		case ruleAction24:
			p.Emit(string(buffer[begin:end]))
		case ruleAction25:
			p.Start(BOOL)
		case ruleAction26:
			p.Emit(string(buffer[begin:end]))
		case ruleAction27:
			p.End()
		case ruleAction28:
			p.Start(FUNC)
		case ruleAction29:
			p.End()
		case ruleAction30:
			p.Start(ARGS)
		case ruleAction31:
			p.End()
		case ruleAction32:
			p.Start(FUNCAPPLY)
		case ruleAction33:
			p.End()
		case ruleAction34:
			p.Start(LIST)
		case ruleAction35:
			p.End()
		case ruleAction36:
			p.Start(TUPLE)
		case ruleAction37:
			p.End()
		case ruleAction38:
			p.Start("map")
		case ruleAction39:
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
								goto l21
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
		/* 4 Expr <- <(Op / FuncApply / Value)> */
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
							add(ruleAction2, position)
						}
						if !_rules[ruleValue]() {
							goto l45
						}
						if !_rules[rulesp]() {
							goto l45
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
												goto l45
											}
											position++
											break
										case '*':
											if buffer[position] != rune('*') {
												goto l45
											}
											position++
											if buffer[position] != rune('*') {
												goto l45
											}
											position++
											break
										case '/':
											if buffer[position] != rune('/') {
												goto l45
											}
											position++
											break
										case '-':
											if buffer[position] != rune('-') {
												goto l45
											}
											position++
											break
										default:
											if buffer[position] != rune('+') {
												goto l45
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
							goto l45
						}
						if !_rules[ruleValue]() {
							goto l45
						}
					l48:
						{
							position49, tokenIndex49, depth49 := position, tokenIndex, depth
							if !_rules[rulesp]() {
								goto l49
							}
							{
								position58 := position
								depth++
								{
									add(ruleAction4, position)
								}
								{
									position60 := position
									depth++
									{
										position61, tokenIndex61, depth61 := position, tokenIndex, depth
										if buffer[position] != rune('*') {
											goto l62
										}
										position++
										goto l61
									l62:
										position, tokenIndex, depth = position61, tokenIndex61, depth61
										{
											switch buffer[position] {
											case '%':
												if buffer[position] != rune('%') {
													goto l49
												}
												position++
												break
											case '*':
												if buffer[position] != rune('*') {
													goto l49
												}
												position++
												if buffer[position] != rune('*') {
													goto l49
												}
												position++
												break
											case '/':
												if buffer[position] != rune('/') {
													goto l49
												}
												position++
												break
											case '-':
												if buffer[position] != rune('-') {
													goto l49
												}
												position++
												break
											default:
												if buffer[position] != rune('+') {
													goto l49
												}
												position++
												break
											}
										}

									}
								l61:
									depth--
									add(rulePegText, position60)
								}
								{
									add(ruleAction5, position)
								}
								{
									add(ruleAction6, position)
								}
								depth--
								add(ruleBinaryOp, position58)
							}
							if !_rules[rulesp]() {
								goto l49
							}
							if !_rules[ruleValue]() {
								goto l49
							}
							goto l48
						l49:
							position, tokenIndex, depth = position49, tokenIndex49, depth49
						}
						{
							add(ruleAction3, position)
						}
						depth--
						add(ruleOp, position46)
					}
					goto l44
				l45:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					{
						position68 := position
						depth++
						{
							add(ruleAction32, position)
						}
						if !_rules[ruleRef]() {
							goto l67
						}
						if !_rules[ruleTuple]() {
							goto l67
						}
						{
							add(ruleAction33, position)
						}
						depth--
						add(ruleFuncApply, position68)
					}
					goto l44
				l67:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
					if !_rules[ruleValue]() {
						goto l42
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
		/* 5 Op <- <(Action2 Value (sp BinaryOp sp Value)+ Action3)> */
		nil,
		/* 6 BinaryOp <- <(Action4 <('*' / ((&('%') '%') | (&('*') ('*' '*')) | (&('/') '/') | (&('-') '-') | (&('+') '+')))> Action5 Action6)> */
		nil,
		/* 7 Statement <- <(Assignment / If)> */
		nil,
		/* 8 Assignment <- <(Action7 LocalRef sp '=' sp Expr Action8)> */
		nil,
		/* 9 If <- <(Action9 ('i' 'f') sp Expr sp Block Action10)> */
		nil,
		/* 10 Ref <- <(FullRef / LocalRef)> */
		func() bool {
			position76, tokenIndex76, depth76 := position, tokenIndex, depth
			{
				position77 := position
				depth++
				{
					position78, tokenIndex78, depth78 := position, tokenIndex, depth
					{
						position80 := position
						depth++
						{
							add(ruleAction11, position)
						}
						{
							position82 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l79
							}
						l83:
							{
								position84, tokenIndex84, depth84 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l84
								}
								goto l83
							l84:
								position, tokenIndex, depth = position84, tokenIndex84, depth84
							}
							depth--
							add(rulePegText, position82)
						}
						{
							add(ruleAction12, position)
						}
						if buffer[position] != rune(':') {
							goto l79
						}
						position++
						{
							position86 := position
							depth++
							if !_rules[ruleRefChar]() {
								goto l79
							}
						l87:
							{
								position88, tokenIndex88, depth88 := position, tokenIndex, depth
								if !_rules[ruleRefChar]() {
									goto l88
								}
								goto l87
							l88:
								position, tokenIndex, depth = position88, tokenIndex88, depth88
							}
							depth--
							add(rulePegText, position86)
						}
						{
							add(ruleAction13, position)
						}
						{
							add(ruleAction14, position)
						}
						depth--
						add(ruleFullRef, position80)
					}
					goto l78
				l79:
					position, tokenIndex, depth = position78, tokenIndex78, depth78
					if !_rules[ruleLocalRef]() {
						goto l76
					}
				}
			l78:
				depth--
				add(ruleRef, position77)
			}
			return true
		l76:
			position, tokenIndex, depth = position76, tokenIndex76, depth76
			return false
		},
		/* 11 FullRef <- <(Action11 <RefChar+> Action12 ':' <RefChar+> Action13 Action14)> */
		nil,
		/* 12 LocalRef <- <(Action15 <RefChar+> Action16 Action17)> */
		func() bool {
			position92, tokenIndex92, depth92 := position, tokenIndex, depth
			{
				position93 := position
				depth++
				{
					add(ruleAction15, position)
				}
				{
					position95 := position
					depth++
					if !_rules[ruleRefChar]() {
						goto l92
					}
				l96:
					{
						position97, tokenIndex97, depth97 := position, tokenIndex, depth
						if !_rules[ruleRefChar]() {
							goto l97
						}
						goto l96
					l97:
						position, tokenIndex, depth = position97, tokenIndex97, depth97
					}
					depth--
					add(rulePegText, position95)
				}
				{
					add(ruleAction16, position)
				}
				{
					add(ruleAction17, position)
				}
				depth--
				add(ruleLocalRef, position93)
			}
			return true
		l92:
			position, tokenIndex, depth = position92, tokenIndex92, depth92
			return false
		},
		/* 13 RefChar <- <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))> */
		func() bool {
			position100, tokenIndex100, depth100 := position, tokenIndex, depth
			{
				position101 := position
				depth++
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l100
						}
						position++
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l100
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l100
						}
						position++
						break
					}
				}

				depth--
				add(ruleRefChar, position101)
			}
			return true
		l100:
			position, tokenIndex, depth = position100, tokenIndex100, depth100
			return false
		},
		/* 14 Value <- <(Ref / Literal)> */
		func() bool {
			position103, tokenIndex103, depth103 := position, tokenIndex, depth
			{
				position104 := position
				depth++
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if !_rules[ruleRef]() {
						goto l106
					}
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					{
						position107 := position
						depth++
						{
							position108, tokenIndex108, depth108 := position, tokenIndex, depth
							{
								position110 := position
								depth++
								{
									add(ruleAction28, position)
								}
								{
									position112 := position
									depth++
									{
										add(ruleAction30, position)
									}
									if buffer[position] != rune('(') {
										goto l109
									}
									position++
									if !_rules[rulesp]() {
										goto l109
									}
									{
										position114, tokenIndex114, depth114 := position, tokenIndex, depth
										if !_rules[ruleLocalRef]() {
											goto l114
										}
									l116:
										{
											position117, tokenIndex117, depth117 := position, tokenIndex, depth
											if !_rules[rulesp]() {
												goto l117
											}
											if buffer[position] != rune(',') {
												goto l117
											}
											position++
											if !_rules[rulesp]() {
												goto l117
											}
											if !_rules[ruleLocalRef]() {
												goto l117
											}
											goto l116
										l117:
											position, tokenIndex, depth = position117, tokenIndex117, depth117
										}
										if !_rules[rulesp]() {
											goto l114
										}
										goto l115
									l114:
										position, tokenIndex, depth = position114, tokenIndex114, depth114
									}
								l115:
									if buffer[position] != rune(')') {
										goto l109
									}
									position++
									{
										add(ruleAction31, position)
									}
									depth--
									add(ruleFuncArgs, position112)
								}
								if !_rules[rulesp]() {
									goto l109
								}
								if buffer[position] != rune('-') {
									goto l109
								}
								position++
								if buffer[position] != rune('>') {
									goto l109
								}
								position++
								if !_rules[rulesp]() {
									goto l109
								}
								{
									position119, tokenIndex119, depth119 := position, tokenIndex, depth
									if !_rules[ruleBlock]() {
										goto l120
									}
									goto l119
								l120:
									position, tokenIndex, depth = position119, tokenIndex119, depth119
									if !_rules[ruleExpr]() {
										goto l109
									}
								}
							l119:
								{
									add(ruleAction29, position)
								}
								depth--
								add(ruleFunc, position110)
							}
							goto l108
						l109:
							position, tokenIndex, depth = position108, tokenIndex108, depth108
							{
								position123 := position
								depth++
								{
									switch buffer[position] {
									case 'f', 't':
										{
											position125 := position
											depth++
											{
												add(ruleAction25, position)
											}
											{
												position127 := position
												depth++
												{
													position128, tokenIndex128, depth128 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l129
													}
													position++
													if buffer[position] != rune('r') {
														goto l129
													}
													position++
													if buffer[position] != rune('u') {
														goto l129
													}
													position++
													if buffer[position] != rune('e') {
														goto l129
													}
													position++
													goto l128
												l129:
													position, tokenIndex, depth = position128, tokenIndex128, depth128
													if buffer[position] != rune('f') {
														goto l122
													}
													position++
													if buffer[position] != rune('a') {
														goto l122
													}
													position++
													if buffer[position] != rune('l') {
														goto l122
													}
													position++
													if buffer[position] != rune('s') {
														goto l122
													}
													position++
													if buffer[position] != rune('e') {
														goto l122
													}
													position++
												}
											l128:
												depth--
												add(rulePegText, position127)
											}
											{
												add(ruleAction26, position)
											}
											{
												add(ruleAction27, position)
											}
											depth--
											add(ruleBoolean, position125)
										}
										break
									case '"':
										{
											position132 := position
											depth++
											{
												add(ruleAction18, position)
											}
											if buffer[position] != rune('"') {
												goto l122
											}
											position++
											{
												position134 := position
												depth++
											l135:
												{
													position136, tokenIndex136, depth136 := position, tokenIndex, depth
													{
														position137 := position
														depth++
														{
															position138, tokenIndex138, depth138 := position, tokenIndex, depth
															{
																position140 := position
																depth++
																{
																	position141 := position
																	depth++
																	if buffer[position] != rune('\\') {
																		goto l139
																	}
																	position++
																	{
																		switch buffer[position] {
																		case 'v':
																			if buffer[position] != rune('v') {
																				goto l139
																			}
																			position++
																			break
																		case 't':
																			if buffer[position] != rune('t') {
																				goto l139
																			}
																			position++
																			break
																		case 'r':
																			if buffer[position] != rune('r') {
																				goto l139
																			}
																			position++
																			break
																		case 'n':
																			if buffer[position] != rune('n') {
																				goto l139
																			}
																			position++
																			break
																		case 'f':
																			if buffer[position] != rune('f') {
																				goto l139
																			}
																			position++
																			break
																		case 'b':
																			if buffer[position] != rune('b') {
																				goto l139
																			}
																			position++
																			break
																		case 'a':
																			if buffer[position] != rune('a') {
																				goto l139
																			}
																			position++
																			break
																		case '\\':
																			if buffer[position] != rune('\\') {
																				goto l139
																			}
																			position++
																			break
																		case '?':
																			if buffer[position] != rune('?') {
																				goto l139
																			}
																			position++
																			break
																		case '"':
																			if buffer[position] != rune('"') {
																				goto l139
																			}
																			position++
																			break
																		default:
																			if buffer[position] != rune('\'') {
																				goto l139
																			}
																			position++
																			break
																		}
																	}

																	depth--
																	add(ruleSimpleEsc, position141)
																}
																depth--
																add(ruleStringEsc, position140)
															}
															goto l138
														l139:
															position, tokenIndex, depth = position138, tokenIndex138, depth138
															{
																position143, tokenIndex143, depth143 := position, tokenIndex, depth
																{
																	switch buffer[position] {
																	case '\\':
																		if buffer[position] != rune('\\') {
																			goto l143
																		}
																		position++
																		break
																	case '\n':
																		if buffer[position] != rune('\n') {
																			goto l143
																		}
																		position++
																		break
																	default:
																		if buffer[position] != rune('"') {
																			goto l143
																		}
																		position++
																		break
																	}
																}

																goto l136
															l143:
																position, tokenIndex, depth = position143, tokenIndex143, depth143
															}
															if !matchDot() {
																goto l136
															}
														}
													l138:
														depth--
														add(ruleStringChar, position137)
													}
													goto l135
												l136:
													position, tokenIndex, depth = position136, tokenIndex136, depth136
												}
												depth--
												add(rulePegText, position134)
											}
											if buffer[position] != rune('"') {
												goto l122
											}
											position++
											{
												add(ruleAction19, position)
											}
											{
												add(ruleAction20, position)
											}
											depth--
											add(ruleString, position132)
										}
										break
									default:
										{
											position147 := position
											depth++
											{
												add(ruleAction21, position)
											}
											{
												position149, tokenIndex149, depth149 := position, tokenIndex, depth
												{
													position151 := position
													depth++
													if !_rules[ruleDecimal]() {
														goto l150
													}
													{
														position152, tokenIndex152, depth152 := position, tokenIndex, depth
														if buffer[position] != rune('e') {
															goto l153
														}
														position++
														goto l152
													l153:
														position, tokenIndex, depth = position152, tokenIndex152, depth152
														if buffer[position] != rune('E') {
															goto l150
														}
														position++
													}
												l152:
													if !_rules[ruleInteger]() {
														goto l150
													}
													depth--
													add(ruleSciNum, position151)
												}
												goto l149
											l150:
												position, tokenIndex, depth = position149, tokenIndex149, depth149
												if !_rules[ruleDecimal]() {
													goto l154
												}
												goto l149
											l154:
												position, tokenIndex, depth = position149, tokenIndex149, depth149
												if !_rules[ruleInteger]() {
													goto l122
												}
											}
										l149:
											{
												add(ruleAction22, position)
											}
											depth--
											add(ruleNumeric, position147)
										}
										break
									}
								}

								depth--
								add(ruleScalar, position123)
							}
							goto l108
						l122:
							position, tokenIndex, depth = position108, tokenIndex108, depth108
							{
								position156 := position
								depth++
								{
									switch buffer[position] {
									case '{':
										{
											position158 := position
											depth++
											{
												add(ruleAction38, position)
											}
											if buffer[position] != rune('{') {
												goto l103
											}
											position++
											if !_rules[rulesp]() {
												goto l103
											}
											{
												position160, tokenIndex160, depth160 := position, tokenIndex, depth
												if !_rules[ruleExpr]() {
													goto l160
												}
												if !_rules[rulesp]() {
													goto l160
												}
												if buffer[position] != rune(':') {
													goto l160
												}
												position++
												if !_rules[rulesp]() {
													goto l160
												}
												if !_rules[ruleExpr]() {
													goto l160
												}
											l162:
												{
													position163, tokenIndex163, depth163 := position, tokenIndex, depth
													if !_rules[rulesp]() {
														goto l163
													}
													if buffer[position] != rune(',') {
														goto l163
													}
													position++
													if !_rules[rulesp]() {
														goto l163
													}
													if !_rules[ruleExpr]() {
														goto l163
													}
													if !_rules[rulesp]() {
														goto l163
													}
													if buffer[position] != rune(':') {
														goto l163
													}
													position++
													if !_rules[rulesp]() {
														goto l163
													}
													if !_rules[ruleExpr]() {
														goto l163
													}
													goto l162
												l163:
													position, tokenIndex, depth = position163, tokenIndex163, depth163
												}
												if !_rules[rulesp]() {
													goto l160
												}
												goto l161
											l160:
												position, tokenIndex, depth = position160, tokenIndex160, depth160
											}
										l161:
											if buffer[position] != rune('}') {
												goto l103
											}
											position++
											{
												add(ruleAction39, position)
											}
											depth--
											add(ruleMap, position158)
										}
										break
									case '(':
										if !_rules[ruleTuple]() {
											goto l103
										}
										break
									default:
										{
											position165 := position
											depth++
											{
												add(ruleAction34, position)
											}
											if buffer[position] != rune('[') {
												goto l103
											}
											position++
											if !_rules[rulesp]() {
												goto l103
											}
											{
												position167, tokenIndex167, depth167 := position, tokenIndex, depth
												if !_rules[ruleExpr]() {
													goto l167
												}
											l169:
												{
													position170, tokenIndex170, depth170 := position, tokenIndex, depth
													if !_rules[rulesp]() {
														goto l170
													}
													if buffer[position] != rune(',') {
														goto l170
													}
													position++
													if !_rules[rulesp]() {
														goto l170
													}
													if !_rules[ruleExpr]() {
														goto l170
													}
													goto l169
												l170:
													position, tokenIndex, depth = position170, tokenIndex170, depth170
												}
												if !_rules[rulesp]() {
													goto l167
												}
												goto l168
											l167:
												position, tokenIndex, depth = position167, tokenIndex167, depth167
											}
										l168:
											if buffer[position] != rune(']') {
												goto l103
											}
											position++
											{
												add(ruleAction35, position)
											}
											depth--
											add(ruleList, position165)
										}
										break
									}
								}

								depth--
								add(ruleVector, position156)
							}
						}
					l108:
						depth--
						add(ruleLiteral, position107)
					}
				}
			l105:
				depth--
				add(ruleValue, position104)
			}
			return true
		l103:
			position, tokenIndex, depth = position103, tokenIndex103, depth103
			return false
		},
		/* 15 Literal <- <(Func / Scalar / Vector)> */
		nil,
		/* 16 Scalar <- <((&('f' | 't') Boolean) | (&('"') String) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric))> */
		nil,
		/* 17 Vector <- <((&('{') Map) | (&('(') Tuple) | (&('[') List))> */
		nil,
		/* 18 String <- <(Action18 '"' <StringChar*> '"' Action19 Action20)> */
		nil,
		/* 19 StringChar <- <(StringEsc / (!((&('\\') '\\') | (&('\n') '\n') | (&('"') '"')) .))> */
		nil,
		/* 20 StringEsc <- <SimpleEsc> */
		nil,
		/* 21 SimpleEsc <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 22 Numeric <- <(Action21 (SciNum / Decimal / Integer) Action22)> */
		nil,
		/* 23 SciNum <- <(Decimal ('e' / 'E') Integer)> */
		nil,
		/* 24 Decimal <- <(Integer '.' <Digit*> Action23)> */
		func() bool {
			position181, tokenIndex181, depth181 := position, tokenIndex, depth
			{
				position182 := position
				depth++
				if !_rules[ruleInteger]() {
					goto l181
				}
				if buffer[position] != rune('.') {
					goto l181
				}
				position++
				{
					position183 := position
					depth++
				l184:
					{
						position185, tokenIndex185, depth185 := position, tokenIndex, depth
						if !_rules[ruleDigit]() {
							goto l185
						}
						goto l184
					l185:
						position, tokenIndex, depth = position185, tokenIndex185, depth185
					}
					depth--
					add(rulePegText, position183)
				}
				{
					add(ruleAction23, position)
				}
				depth--
				add(ruleDecimal, position182)
			}
			return true
		l181:
			position, tokenIndex, depth = position181, tokenIndex181, depth181
			return false
		},
		/* 25 Integer <- <(<WholeNum> Action24)> */
		func() bool {
			position187, tokenIndex187, depth187 := position, tokenIndex, depth
			{
				position188 := position
				depth++
				{
					position189 := position
					depth++
					{
						position190 := position
						depth++
						{
							position191, tokenIndex191, depth191 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l192
							}
							position++
							goto l191
						l192:
							position, tokenIndex, depth = position191, tokenIndex191, depth191
							{
								position193, tokenIndex193, depth193 := position, tokenIndex, depth
								if buffer[position] != rune('-') {
									goto l193
								}
								position++
								goto l194
							l193:
								position, tokenIndex, depth = position193, tokenIndex193, depth193
							}
						l194:
							if c := buffer[position]; c < rune('1') || c > rune('9') {
								goto l187
							}
							position++
						l195:
							{
								position196, tokenIndex196, depth196 := position, tokenIndex, depth
								if !_rules[ruleDigit]() {
									goto l196
								}
								goto l195
							l196:
								position, tokenIndex, depth = position196, tokenIndex196, depth196
							}
						}
					l191:
						depth--
						add(ruleWholeNum, position190)
					}
					depth--
					add(rulePegText, position189)
				}
				{
					add(ruleAction24, position)
				}
				depth--
				add(ruleInteger, position188)
			}
			return true
		l187:
			position, tokenIndex, depth = position187, tokenIndex187, depth187
			return false
		},
		/* 26 WholeNum <- <('0' / ('-'? [1-9] Digit*))> */
		nil,
		/* 27 Digit <- <[0-9]> */
		func() bool {
			position199, tokenIndex199, depth199 := position, tokenIndex, depth
			{
				position200 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l199
				}
				position++
				depth--
				add(ruleDigit, position200)
			}
			return true
		l199:
			position, tokenIndex, depth = position199, tokenIndex199, depth199
			return false
		},
		/* 28 Boolean <- <(Action25 <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> Action26 Action27)> */
		nil,
		/* 29 Func <- <(Action28 FuncArgs sp ('-' '>') sp (Block / Expr) Action29)> */
		nil,
		/* 30 FuncArgs <- <(Action30 '(' sp (LocalRef (sp ',' sp LocalRef)* sp)? ')' Action31)> */
		nil,
		/* 31 FuncApply <- <(Action32 Ref Tuple Action33)> */
		nil,
		/* 32 List <- <(Action34 '[' sp (Expr (sp ',' sp Expr)* sp)? ']' Action35)> */
		nil,
		/* 33 Tuple <- <(Action36 '(' sp (Expr (sp ',' sp Expr)* sp)? ')' Action37)> */
		func() bool {
			position206, tokenIndex206, depth206 := position, tokenIndex, depth
			{
				position207 := position
				depth++
				{
					add(ruleAction36, position)
				}
				if buffer[position] != rune('(') {
					goto l206
				}
				position++
				if !_rules[rulesp]() {
					goto l206
				}
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if !_rules[ruleExpr]() {
						goto l209
					}
				l211:
					{
						position212, tokenIndex212, depth212 := position, tokenIndex, depth
						if !_rules[rulesp]() {
							goto l212
						}
						if buffer[position] != rune(',') {
							goto l212
						}
						position++
						if !_rules[rulesp]() {
							goto l212
						}
						if !_rules[ruleExpr]() {
							goto l212
						}
						goto l211
					l212:
						position, tokenIndex, depth = position212, tokenIndex212, depth212
					}
					if !_rules[rulesp]() {
						goto l209
					}
					goto l210
				l209:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
				}
			l210:
				if buffer[position] != rune(')') {
					goto l206
				}
				position++
				{
					add(ruleAction37, position)
				}
				depth--
				add(ruleTuple, position207)
			}
			return true
		l206:
			position, tokenIndex, depth = position206, tokenIndex206, depth206
			return false
		},
		/* 34 Map <- <(Action38 '{' sp (Expr sp ':' sp Expr (sp ',' sp Expr sp ':' sp Expr)* sp)? '}' Action39)> */
		nil,
		/* 35 Gravitasse <- <'@'> */
		nil,
		/* 36 msp <- <(ws / comment)+> */
		nil,
		/* 37 sp <- <(ws / comment)*> */
		func() bool {
			{
				position218 := position
				depth++
			l219:
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					{
						position221, tokenIndex221, depth221 := position, tokenIndex, depth
						if !_rules[rulews]() {
							goto l222
						}
						goto l221
					l222:
						position, tokenIndex, depth = position221, tokenIndex221, depth221
						if !_rules[rulecomment]() {
							goto l220
						}
					}
				l221:
					goto l219
				l220:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
				}
				depth--
				add(rulesp, position218)
			}
			return true
		},
		/* 38 comment <- <('#' (!'\n' .)*)> */
		func() bool {
			position223, tokenIndex223, depth223 := position, tokenIndex, depth
			{
				position224 := position
				depth++
				if buffer[position] != rune('#') {
					goto l223
				}
				position++
			l225:
				{
					position226, tokenIndex226, depth226 := position, tokenIndex, depth
					{
						position227, tokenIndex227, depth227 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l227
						}
						position++
						goto l226
					l227:
						position, tokenIndex, depth = position227, tokenIndex227, depth227
					}
					if !matchDot() {
						goto l226
					}
					goto l225
				l226:
					position, tokenIndex, depth = position226, tokenIndex226, depth226
				}
				depth--
				add(rulecomment, position224)
			}
			return true
		l223:
			position, tokenIndex, depth = position223, tokenIndex223, depth223
			return false
		},
		/* 39 ws <- <((&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))> */
		func() bool {
			position228, tokenIndex228, depth228 := position, tokenIndex, depth
			{
				position229 := position
				depth++
				{
					switch buffer[position] {
					case '\r':
						if buffer[position] != rune('\r') {
							goto l228
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l228
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l228
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l228
						}
						position++
						break
					}
				}

				depth--
				add(rulews, position229)
			}
			return true
		l228:
			position, tokenIndex, depth = position228, tokenIndex228, depth228
			return false
		},
		/* 41 Action0 <- <{ p.Start(RIFT) }> */
		nil,
		/* 42 Action1 <- <{ p.End() }> */
		nil,
		/* 43 Action2 <- <{ p.Start(OP) }> */
		nil,
		/* 44 Action3 <- <{ p.End() }> */
		nil,
		/* 45 Action4 <- <{ p.Start(BINOP) }> */
		nil,
		nil,
		/* 47 Action5 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 48 Action6 <- <{ p.End() }> */
		nil,
		/* 49 Action7 <- <{ p.Start(ASSIGNMENT) }> */
		nil,
		/* 50 Action8 <- <{ p.End() }> */
		nil,
		/* 51 Action9 <- <{ p.Start(IF) }> */
		nil,
		/* 52 Action10 <- <{ p.End() }> */
		nil,
		/* 53 Action11 <- <{ p.Start(REF) }> */
		nil,
		/* 54 Action12 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 55 Action13 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 56 Action14 <- <{ p.End() }> */
		nil,
		/* 57 Action15 <- <{ p.Start(REF) }> */
		nil,
		/* 58 Action16 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 59 Action17 <- <{ p.End() }> */
		nil,
		/* 60 Action18 <- <{ p.Start(STRING) }> */
		nil,
		/* 61 Action19 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 62 Action20 <- <{ p.End() }> */
		nil,
		/* 63 Action21 <- <{ p.Start(NUM) }> */
		nil,
		/* 64 Action22 <- <{ p.End() }> */
		nil,
		/* 65 Action23 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 66 Action24 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 67 Action25 <- <{ p.Start(BOOL) }> */
		nil,
		/* 68 Action26 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 69 Action27 <- <{ p.End() }> */
		nil,
		/* 70 Action28 <- <{ p.Start(FUNC) }> */
		nil,
		/* 71 Action29 <- <{ p.End() }> */
		nil,
		/* 72 Action30 <- <{ p.Start(ARGS) }> */
		nil,
		/* 73 Action31 <- <{ p.End() }> */
		nil,
		/* 74 Action32 <- <{ p.Start(FUNCAPPLY) }> */
		nil,
		/* 75 Action33 <- <{ p.End() }> */
		nil,
		/* 76 Action34 <- <{ p.Start(LIST) }> */
		nil,
		/* 77 Action35 <- <{ p.End() }> */
		nil,
		/* 78 Action36 <- <{ p.Start(TUPLE) }> */
		nil,
		/* 79 Action37 <- <{ p.End() }> */
		nil,
		/* 80 Action38 <- <{ p.Start("map") }> */
		nil,
		/* 81 Action39 <- <{ p.End() }> */
		nil,
	}
	p.rules = _rules
}
