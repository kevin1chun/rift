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
	ruleName
	ruleValue
	ruleLiteral
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
	ruleFuncApply
	ruleList
	ruleTuple
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
	"Name",
	"Value",
	"Literal",
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
	"FuncApply",
	"List",
	"Tuple",
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
	rules  [68]func() bool
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
			p.Emit(string(buffer[begin:end]))
		case ruleAction11:
			p.End()
		case ruleAction12:
			p.Start(STRING)
		case ruleAction13:
			p.Emit(string(buffer[begin:end]))
		case ruleAction14:
			p.End()
		case ruleAction15:
			p.Start(SCI)
		case ruleAction16:
			p.End()
		case ruleAction17:
			p.Start(DEC)
		case ruleAction18:
			p.Emit(string(buffer[begin:end]))
		case ruleAction19:
			p.End()
		case ruleAction20:
			p.Start(INT)
		case ruleAction21:
			p.Emit(string(buffer[begin:end]))
		case ruleAction22:
			p.End()
		case ruleAction23:
			p.Start(BOOL)
		case ruleAction24:
			p.Emit(string(buffer[begin:end]))
		case ruleAction25:
			p.End()
		case ruleAction26:
			p.Start(FUNC)
		case ruleAction27:
			p.Start(ARGS)
		case ruleAction28:
			p.End()
		case ruleAction29:
			p.End()
		case ruleAction30:
			p.Start(FUNCAPPLY)
		case ruleAction31:
			p.End()
		case ruleAction32:
			p.Start(LIST)
		case ruleAction33:
			p.End()
		case ruleAction34:
			p.Start(TUPLE)
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
					if !_rules[ruleName]() {
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
						if !_rules[ruleName]() {
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
		/* 1 Rift <- <(Action0 Name sp ('=' '>') sp Block Action1)> */
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
								if !_rules[ruleName]() {
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
							add(ruleAction30, position)
						}
						if !_rules[ruleName]() {
							goto l55
						}
						if !_rules[ruleTuple]() {
							goto l55
						}
						{
							add(ruleAction31, position)
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
		/* 7 Assignment <- <(Action7 Name sp '=' sp Expr Action8)> */
		nil,
		/* 8 Name <- <(Action9 <((&('_') '_') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+> Action10 Action11)> */
		func() bool {
			position62, tokenIndex62, depth62 := position, tokenIndex, depth
			{
				position63 := position
				depth++
				{
					add(ruleAction9, position)
				}
				{
					position65 := position
					depth++
					{
						switch buffer[position] {
						case '_':
							if buffer[position] != rune('_') {
								goto l62
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l62
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l62
							}
							position++
							break
						}
					}

				l66:
					{
						position67, tokenIndex67, depth67 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '_':
								if buffer[position] != rune('_') {
									goto l67
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l67
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l67
								}
								position++
								break
							}
						}

						goto l66
					l67:
						position, tokenIndex, depth = position67, tokenIndex67, depth67
					}
					depth--
					add(rulePegText, position65)
				}
				{
					add(ruleAction10, position)
				}
				{
					add(ruleAction11, position)
				}
				depth--
				add(ruleName, position63)
			}
			return true
		l62:
			position, tokenIndex, depth = position62, tokenIndex62, depth62
			return false
		},
		/* 9 Value <- <(Name / Literal)> */
		func() bool {
			position72, tokenIndex72, depth72 := position, tokenIndex, depth
			{
				position73 := position
				depth++
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if !_rules[ruleName]() {
						goto l75
					}
					goto l74
				l75:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
					{
						position76 := position
						depth++
						{
							position77, tokenIndex77, depth77 := position, tokenIndex, depth
							{
								position79 := position
								depth++
								{
									add(ruleAction26, position)
								}
								{
									add(ruleAction27, position)
								}
								if buffer[position] != rune('(') {
									goto l78
								}
								position++
								if !_rules[rulesp]() {
									goto l78
								}
								{
									position82, tokenIndex82, depth82 := position, tokenIndex, depth
									if !_rules[ruleName]() {
										goto l82
									}
								l84:
									{
										position85, tokenIndex85, depth85 := position, tokenIndex, depth
										if !_rules[rulesp]() {
											goto l85
										}
										if buffer[position] != rune(',') {
											goto l85
										}
										position++
										if !_rules[rulesp]() {
											goto l85
										}
										if !_rules[ruleName]() {
											goto l85
										}
										goto l84
									l85:
										position, tokenIndex, depth = position85, tokenIndex85, depth85
									}
									if !_rules[rulesp]() {
										goto l82
									}
									goto l83
								l82:
									position, tokenIndex, depth = position82, tokenIndex82, depth82
								}
							l83:
								if buffer[position] != rune(')') {
									goto l78
								}
								position++
								{
									add(ruleAction28, position)
								}
								if !_rules[rulesp]() {
									goto l78
								}
								if buffer[position] != rune('-') {
									goto l78
								}
								position++
								if buffer[position] != rune('>') {
									goto l78
								}
								position++
								if !_rules[rulesp]() {
									goto l78
								}
								{
									position87, tokenIndex87, depth87 := position, tokenIndex, depth
									if !_rules[ruleBlock]() {
										goto l88
									}
									goto l87
								l88:
									position, tokenIndex, depth = position87, tokenIndex87, depth87
									if !_rules[ruleExpr]() {
										goto l78
									}
								}
							l87:
								{
									add(ruleAction29, position)
								}
								depth--
								add(ruleFunc, position79)
							}
							goto l77
						l78:
							position, tokenIndex, depth = position77, tokenIndex77, depth77
							{
								switch buffer[position] {
								case '(':
									if !_rules[ruleTuple]() {
										goto l72
									}
									break
								case '[':
									{
										position91 := position
										depth++
										{
											add(ruleAction32, position)
										}
										if buffer[position] != rune('[') {
											goto l72
										}
										position++
										if !_rules[rulesp]() {
											goto l72
										}
										{
											position93, tokenIndex93, depth93 := position, tokenIndex, depth
											if !_rules[ruleExpr]() {
												goto l93
											}
										l95:
											{
												position96, tokenIndex96, depth96 := position, tokenIndex, depth
												if !_rules[rulesp]() {
													goto l96
												}
												if buffer[position] != rune(',') {
													goto l96
												}
												position++
												if !_rules[rulesp]() {
													goto l96
												}
												if !_rules[ruleExpr]() {
													goto l96
												}
												goto l95
											l96:
												position, tokenIndex, depth = position96, tokenIndex96, depth96
											}
											if !_rules[rulesp]() {
												goto l93
											}
											goto l94
										l93:
											position, tokenIndex, depth = position93, tokenIndex93, depth93
										}
									l94:
										if buffer[position] != rune(']') {
											goto l72
										}
										position++
										{
											add(ruleAction33, position)
										}
										depth--
										add(ruleList, position91)
									}
									break
								case 'f', 't':
									{
										position98 := position
										depth++
										{
											add(ruleAction23, position)
										}
										{
											position100 := position
											depth++
											{
												position101, tokenIndex101, depth101 := position, tokenIndex, depth
												if buffer[position] != rune('t') {
													goto l102
												}
												position++
												if buffer[position] != rune('r') {
													goto l102
												}
												position++
												if buffer[position] != rune('u') {
													goto l102
												}
												position++
												if buffer[position] != rune('e') {
													goto l102
												}
												position++
												goto l101
											l102:
												position, tokenIndex, depth = position101, tokenIndex101, depth101
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
										l101:
											depth--
											add(rulePegText, position100)
										}
										{
											add(ruleAction24, position)
										}
										{
											add(ruleAction25, position)
										}
										depth--
										add(ruleBoolean, position98)
									}
									break
								case '"':
									{
										position105 := position
										depth++
										{
											add(ruleAction12, position)
										}
										if buffer[position] != rune('"') {
											goto l72
										}
										position++
										{
											position107 := position
											depth++
										l108:
											{
												position109, tokenIndex109, depth109 := position, tokenIndex, depth
												{
													position110 := position
													depth++
													{
														position111, tokenIndex111, depth111 := position, tokenIndex, depth
														{
															position113 := position
															depth++
															{
																position114 := position
																depth++
																if buffer[position] != rune('\\') {
																	goto l112
																}
																position++
																{
																	switch buffer[position] {
																	case 'v':
																		if buffer[position] != rune('v') {
																			goto l112
																		}
																		position++
																		break
																	case 't':
																		if buffer[position] != rune('t') {
																			goto l112
																		}
																		position++
																		break
																	case 'r':
																		if buffer[position] != rune('r') {
																			goto l112
																		}
																		position++
																		break
																	case 'n':
																		if buffer[position] != rune('n') {
																			goto l112
																		}
																		position++
																		break
																	case 'f':
																		if buffer[position] != rune('f') {
																			goto l112
																		}
																		position++
																		break
																	case 'b':
																		if buffer[position] != rune('b') {
																			goto l112
																		}
																		position++
																		break
																	case 'a':
																		if buffer[position] != rune('a') {
																			goto l112
																		}
																		position++
																		break
																	case '\\':
																		if buffer[position] != rune('\\') {
																			goto l112
																		}
																		position++
																		break
																	case '?':
																		if buffer[position] != rune('?') {
																			goto l112
																		}
																		position++
																		break
																	case '"':
																		if buffer[position] != rune('"') {
																			goto l112
																		}
																		position++
																		break
																	default:
																		if buffer[position] != rune('\'') {
																			goto l112
																		}
																		position++
																		break
																	}
																}

																depth--
																add(ruleSimpleEsc, position114)
															}
															depth--
															add(ruleStringEsc, position113)
														}
														goto l111
													l112:
														position, tokenIndex, depth = position111, tokenIndex111, depth111
														{
															position116, tokenIndex116, depth116 := position, tokenIndex, depth
															{
																switch buffer[position] {
																case '\\':
																	if buffer[position] != rune('\\') {
																		goto l116
																	}
																	position++
																	break
																case '\n':
																	if buffer[position] != rune('\n') {
																		goto l116
																	}
																	position++
																	break
																default:
																	if buffer[position] != rune('"') {
																		goto l116
																	}
																	position++
																	break
																}
															}

															goto l109
														l116:
															position, tokenIndex, depth = position116, tokenIndex116, depth116
														}
														if !matchDot() {
															goto l109
														}
													}
												l111:
													depth--
													add(ruleStringChar, position110)
												}
												goto l108
											l109:
												position, tokenIndex, depth = position109, tokenIndex109, depth109
											}
											depth--
											add(rulePegText, position107)
										}
										if buffer[position] != rune('"') {
											goto l72
										}
										position++
										{
											add(ruleAction13, position)
										}
										{
											add(ruleAction14, position)
										}
										depth--
										add(ruleString, position105)
									}
									break
								default:
									{
										position120 := position
										depth++
										{
											position121, tokenIndex121, depth121 := position, tokenIndex, depth
											{
												position123 := position
												depth++
												{
													add(ruleAction15, position)
												}
												if !_rules[ruleDecimal]() {
													goto l122
												}
												{
													position125, tokenIndex125, depth125 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l126
													}
													position++
													goto l125
												l126:
													position, tokenIndex, depth = position125, tokenIndex125, depth125
													if buffer[position] != rune('E') {
														goto l122
													}
													position++
												}
											l125:
												if !_rules[ruleInteger]() {
													goto l122
												}
												{
													add(ruleAction16, position)
												}
												depth--
												add(ruleSciNum, position123)
											}
											goto l121
										l122:
											position, tokenIndex, depth = position121, tokenIndex121, depth121
											if !_rules[ruleDecimal]() {
												goto l128
											}
											goto l121
										l128:
											position, tokenIndex, depth = position121, tokenIndex121, depth121
											if !_rules[ruleInteger]() {
												goto l72
											}
										}
									l121:
										depth--
										add(ruleNumeric, position120)
									}
									break
								}
							}

						}
					l77:
						depth--
						add(ruleLiteral, position76)
					}
				}
			l74:
				depth--
				add(ruleValue, position73)
			}
			return true
		l72:
			position, tokenIndex, depth = position72, tokenIndex72, depth72
			return false
		},
		/* 10 Literal <- <(Func / ((&('(') Tuple) | (&('[') List) | (&('f' | 't') Boolean) | (&('"') String) | (&('-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') Numeric)))> */
		nil,
		/* 11 String <- <(Action12 '"' <StringChar*> '"' Action13 Action14)> */
		nil,
		/* 12 StringChar <- <(StringEsc / (!((&('\\') '\\') | (&('\n') '\n') | (&('"') '"')) .))> */
		nil,
		/* 13 StringEsc <- <SimpleEsc> */
		nil,
		/* 14 SimpleEsc <- <('\\' ((&('v') 'v') | (&('t') 't') | (&('r') 'r') | (&('n') 'n') | (&('f') 'f') | (&('b') 'b') | (&('a') 'a') | (&('\\') '\\') | (&('?') '?') | (&('"') '"') | (&('\'') '\'')))> */
		nil,
		/* 15 Numeric <- <(SciNum / Decimal / Integer)> */
		nil,
		/* 16 SciNum <- <(Action15 Decimal ('e' / 'E') Integer Action16)> */
		nil,
		/* 17 Decimal <- <(Action17 Integer '.' <Digit*> Action18 Action19)> */
		func() bool {
			position136, tokenIndex136, depth136 := position, tokenIndex, depth
			{
				position137 := position
				depth++
				{
					add(ruleAction17, position)
				}
				if !_rules[ruleInteger]() {
					goto l136
				}
				if buffer[position] != rune('.') {
					goto l136
				}
				position++
				{
					position139 := position
					depth++
				l140:
					{
						position141, tokenIndex141, depth141 := position, tokenIndex, depth
						if !_rules[ruleDigit]() {
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
					add(ruleAction18, position)
				}
				{
					add(ruleAction19, position)
				}
				depth--
				add(ruleDecimal, position137)
			}
			return true
		l136:
			position, tokenIndex, depth = position136, tokenIndex136, depth136
			return false
		},
		/* 18 Integer <- <(Action20 <WholeNum> Action21 Action22)> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			{
				position145 := position
				depth++
				{
					add(ruleAction20, position)
				}
				{
					position147 := position
					depth++
					{
						position148 := position
						depth++
						{
							position149, tokenIndex149, depth149 := position, tokenIndex, depth
							if buffer[position] != rune('0') {
								goto l150
							}
							position++
							goto l149
						l150:
							position, tokenIndex, depth = position149, tokenIndex149, depth149
							{
								position151, tokenIndex151, depth151 := position, tokenIndex, depth
								if buffer[position] != rune('-') {
									goto l151
								}
								position++
								goto l152
							l151:
								position, tokenIndex, depth = position151, tokenIndex151, depth151
							}
						l152:
							if c := buffer[position]; c < rune('1') || c > rune('9') {
								goto l144
							}
							position++
						l153:
							{
								position154, tokenIndex154, depth154 := position, tokenIndex, depth
								if !_rules[ruleDigit]() {
									goto l154
								}
								goto l153
							l154:
								position, tokenIndex, depth = position154, tokenIndex154, depth154
							}
						}
					l149:
						depth--
						add(ruleWholeNum, position148)
					}
					depth--
					add(rulePegText, position147)
				}
				{
					add(ruleAction21, position)
				}
				{
					add(ruleAction22, position)
				}
				depth--
				add(ruleInteger, position145)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 19 WholeNum <- <('0' / ('-'? [1-9] Digit*))> */
		nil,
		/* 20 Digit <- <[0-9]> */
		func() bool {
			position158, tokenIndex158, depth158 := position, tokenIndex, depth
			{
				position159 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l158
				}
				position++
				depth--
				add(ruleDigit, position159)
			}
			return true
		l158:
			position, tokenIndex, depth = position158, tokenIndex158, depth158
			return false
		},
		/* 21 Boolean <- <(Action23 <(('t' 'r' 'u' 'e') / ('f' 'a' 'l' 's' 'e'))> Action24 Action25)> */
		nil,
		/* 22 Func <- <(Action26 Action27 '(' sp (Name (sp ',' sp Name)* sp)? ')' Action28 sp ('-' '>') sp (Block / Expr) Action29)> */
		nil,
		/* 23 FuncApply <- <(Action30 Name Tuple Action31)> */
		nil,
		/* 24 List <- <(Action32 '[' sp (Expr (sp ',' sp Expr)* sp)? ']' Action33)> */
		nil,
		/* 25 Tuple <- <(Action34 '(' sp (Expr (sp ',' sp Expr)* sp)? ')' Action35)> */
		func() bool {
			position164, tokenIndex164, depth164 := position, tokenIndex, depth
			{
				position165 := position
				depth++
				{
					add(ruleAction34, position)
				}
				if buffer[position] != rune('(') {
					goto l164
				}
				position++
				if !_rules[rulesp]() {
					goto l164
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
				if buffer[position] != rune(')') {
					goto l164
				}
				position++
				{
					add(ruleAction35, position)
				}
				depth--
				add(ruleTuple, position165)
			}
			return true
		l164:
			position, tokenIndex, depth = position164, tokenIndex164, depth164
			return false
		},
		/* 26 msp <- <(ws / comment)+> */
		nil,
		/* 27 sp <- <(ws / comment)*> */
		func() bool {
			{
				position174 := position
				depth++
			l175:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					{
						position177, tokenIndex177, depth177 := position, tokenIndex, depth
						if !_rules[rulews]() {
							goto l178
						}
						goto l177
					l178:
						position, tokenIndex, depth = position177, tokenIndex177, depth177
						if !_rules[rulecomment]() {
							goto l176
						}
					}
				l177:
					goto l175
				l176:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
				}
				depth--
				add(rulesp, position174)
			}
			return true
		},
		/* 28 comment <- <('#' (!'\n' .)*)> */
		func() bool {
			position179, tokenIndex179, depth179 := position, tokenIndex, depth
			{
				position180 := position
				depth++
				if buffer[position] != rune('#') {
					goto l179
				}
				position++
			l181:
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					{
						position183, tokenIndex183, depth183 := position, tokenIndex, depth
						if buffer[position] != rune('\n') {
							goto l183
						}
						position++
						goto l182
					l183:
						position, tokenIndex, depth = position183, tokenIndex183, depth183
					}
					if !matchDot() {
						goto l182
					}
					goto l181
				l182:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
				}
				depth--
				add(rulecomment, position180)
			}
			return true
		l179:
			position, tokenIndex, depth = position179, tokenIndex179, depth179
			return false
		},
		/* 29 ws <- <((&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))> */
		func() bool {
			position184, tokenIndex184, depth184 := position, tokenIndex, depth
			{
				position185 := position
				depth++
				{
					switch buffer[position] {
					case '\r':
						if buffer[position] != rune('\r') {
							goto l184
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l184
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l184
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l184
						}
						position++
						break
					}
				}

				depth--
				add(rulews, position185)
			}
			return true
		l184:
			position, tokenIndex, depth = position184, tokenIndex184, depth184
			return false
		},
		/* 31 Action0 <- <{ p.Start(RIFT) }> */
		nil,
		/* 32 Action1 <- <{ p.End() }> */
		nil,
		/* 33 Action2 <- <{ p.Start(OP) }> */
		nil,
		/* 34 Action3 <- <{ p.End() }> */
		nil,
		/* 35 Action4 <- <{ p.Start(BINOP) }> */
		nil,
		nil,
		/* 37 Action5 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 38 Action6 <- <{ p.End() }> */
		nil,
		/* 39 Action7 <- <{ p.Start(ASSIGNMENT) }> */
		nil,
		/* 40 Action8 <- <{ p.End() }> */
		nil,
		/* 41 Action9 <- <{ p.Start(REF) }> */
		nil,
		/* 42 Action10 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 43 Action11 <- <{ p.End() }> */
		nil,
		/* 44 Action12 <- <{ p.Start(STRING) }> */
		nil,
		/* 45 Action13 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 46 Action14 <- <{ p.End() }> */
		nil,
		/* 47 Action15 <- <{ p.Start(SCI) }> */
		nil,
		/* 48 Action16 <- <{ p.End() }> */
		nil,
		/* 49 Action17 <- <{ p.Start(DEC) }> */
		nil,
		/* 50 Action18 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 51 Action19 <- <{ p.End() }> */
		nil,
		/* 52 Action20 <- <{ p.Start(INT) }> */
		nil,
		/* 53 Action21 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 54 Action22 <- <{ p.End() }> */
		nil,
		/* 55 Action23 <- <{ p.Start(BOOL) }> */
		nil,
		/* 56 Action24 <- <{ p.Emit(string(buffer[begin:end])) }> */
		nil,
		/* 57 Action25 <- <{ p.End() }> */
		nil,
		/* 58 Action26 <- <{ p.Start(FUNC) }> */
		nil,
		/* 59 Action27 <- <{ p.Start(ARGS) }> */
		nil,
		/* 60 Action28 <- <{ p.End() }> */
		nil,
		/* 61 Action29 <- <{ p.End() }> */
		nil,
		/* 62 Action30 <- <{ p.Start(FUNCAPPLY) }> */
		nil,
		/* 63 Action31 <- <{ p.End() }> */
		nil,
		/* 64 Action32 <- <{ p.Start(LIST) }> */
		nil,
		/* 65 Action33 <- <{ p.End() }> */
		nil,
		/* 66 Action34 <- <{ p.Start(TUPLE) }> */
		nil,
		/* 67 Action35 <- <{ p.End() }> */
		nil,
	}
	p.rules = _rules
}
