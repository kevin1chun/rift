package lang

type riftParser Peg {
	parseStack
}

Source     <- sp (Rift sp)+ !.

Rift       <- { p.Start(RIFT) } '@'? LocalRef sp '=>' sp Block { p.End() }

# TODO: Do you have to use an msp here? I wonder if there is another way to delimit lines
Block      <- '{' sp (Line msp)* '}'

Line       <- Assignment / Expr

Expr       <- Op / FuncApply / Value

Op         <- { p.Start(OP) } Value (sp BinaryOp sp Value)+ { p.End() }

BinaryOp   <- { p.Start(BINOP) } <'+' / '-' / '*' / '/' / '**' / '%'> { p.Emit(string(buffer[begin:end])) } { p.End() }

Assignment <- { p.Start(ASSIGNMENT) } LocalRef sp '=' sp Expr { p.End() }

Ref        <- { p.Start(REF) } (FullRef / LocalRef) { p.End() }

FullRef    <- <RefChar+> { p.Emit(string(buffer[begin:end])) } ':' <RefChar+> { p.Emit(string(buffer[begin:end])) }

LocalRef   <- <RefChar+> { p.Emit(string(buffer[begin:end])) }

RefChar    <- [[a-z_]]

Value      <- Ref / Literal

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

msp        <- (ws / comment)+

sp         <- (ws / comment)*

comment    <- '#'  (!'\n' .)*

ws         <- [ \t\n\r]
