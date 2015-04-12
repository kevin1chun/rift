package lang

type riftParser Peg {
	parseStack
}

Source     <- sp (Rift sp)+ !.

# TODO: Should gravitasse be allowable for any ref?
Rift       <- { p.Start(RIFT) } Gravitasse? LocalRef sp '=>' sp Block { p.End() }

# TODO: Do you have to use an msp here? I wonder if there is another way to delimit lines
Block      <- { p.Start(BLOCK) } '{' sp (Line msp)* '}' { p.End() }

Line       <- Statement / Expr

Expr       <- (!Op Single) / Op

Single     <- If / FuncApply / Value

Op         <- { p.Start(OP) } Single (sp BinaryOp sp Expr)+ { p.End() }

# TODO: Break down by operator type? 
# TODO: Should we even treat operators specially?
BinaryOp   <- { p.Start(BINOP) } <'**' / '>=' / '<=' / '==' / '+' / '-' / '*' / '/' / '%' / '>' / '<'> { p.Emit(string(buffer[begin:end])) } { p.End() }

Statement  <- Assignment / If

Assignment <- { p.Start(ASSIGNMENT) } LocalRef sp '=' sp Expr { p.End() }

If         <- { p.Start(IF) } 'if' sp Expr sp Block (sp 'else' sp Block)? { p.End() }

Ref        <- FullRef / LocalRef

FullRef    <- { p.Start(REF) } <RefChar+> { p.Emit(string(buffer[begin:end])) } ':' <RefChar+> { p.Emit(string(buffer[begin:end])) } { p.End() }

LocalRef   <- { p.Start(REF) } <RefChar+> { p.Emit(string(buffer[begin:end])) } { p.End() }

RefChar    <- [[a-z_]]

Value      <- Literal / Ref

Literal    <- Func / Scalar / Vector

Scalar     <- String / Numeric / Boolean

Vector     <- List / Tuple / Map

String     <- { p.Start(STRING) } '"' <StringChar*> '"' { p.Emit(string(buffer[begin:end])) } { p.End() }

StringChar <- StringEsc / ![\"\n\\] .

StringEsc  <- SimpleEsc

SimpleEsc  <- '\\' ['\"?\\abfnrtv]

Numeric    <- { p.Start(NUM) } (SciNum / Decimal / Integer) { p.End() }

SciNum     <- Decimal [[e]] Integer

Decimal    <- Integer '.' <Digit*> { p.Emit(string(buffer[begin:end])) }

Integer    <- <WholeNum> { p.Emit(string(buffer[begin:end])) }

WholeNum   <- '0' / '-'? [1-9] Digit*

Digit      <- [0-9]

Boolean    <- { p.Start(BOOL) } <'true' / 'false'> { p.Emit(string(buffer[begin:end])) } { p.End() }

Func       <- { p.Start(FUNC) } FuncArgs sp '->' sp (Block / Expr)  { p.End() }

FuncArgs   <- { p.Start(ARGS) } '(' sp (LocalRef (sp ',' sp LocalRef)* sp)? ')' { p.End() }

FuncApply  <- { p.Start(FUNCAPPLY) } Ref Tuple { p.End() }

List       <- { p.Start(LIST) } '[' sp (Expr (sp ',' sp Expr)* sp)? ']' { p.End() }

Tuple      <- { p.Start(TUPLE) } '(' sp (Expr (sp ',' sp Expr)* sp)? ')' { p.End() }

Map        <- { p.Start("map") } '{' sp (Expr sp ':' sp Expr (sp ',' sp Expr sp ':' sp Expr)* sp)? '}' { p.End() }

Gravitasse <- '@'

msp        <- (ws / comment)+

sp         <- (ws / comment)*

comment    <- '#'  (!'\n' .)*

ws         <- [ \t\n\r]
