package lang

type riftParser Peg {
	parseStack
}

Source     <- sp (Rift sp)+ !.

Rift       <- { p.Start(RIFT) } Name sp '=>' sp Block { p.End() }

Block      <- '{' sp (Line msp)* '}'

Line       <- Assignment / Expr

Expr       <- Op / FuncApply / Value

Op         <- { p.Start(OP) } Value (sp BinaryOp sp Value)+ { p.End() }

BinaryOp   <- { p.Start(BINOP) } <'+' / '-' / '*' / '/' / '**' / '%'> { p.Emit(string(buffer[begin:end])) } { p.End() }

Assignment <- { p.Start(ASSIGNMENT) } Name sp '=' sp Expr { p.End() }

Name       <- { p.Start(REF) } <[[a-z_]]+> { p.Emit(string(buffer[begin:end])) } { p.End() }

Value      <- Name / Literal

Literal    <- String / Numeric / Boolean / List / Func / Tuple

String     <- { p.Start(STRING) } '"' <StringChar*> '"' { p.Emit(string(buffer[begin:end])) } { p.End() }

StringChar <- StringEsc / ![\"\n\\] .

StringEsc  <- SimpleEsc

SimpleEsc  <- '\\' ['\"?\\abfnrtv]

Numeric    <- SciNum / Decimal / Integer

SciNum     <- { p.Start(SCI) } Decimal [[e]] Integer { p.End() }

Decimal    <- { p.Start(DEC) } Integer '.' <Digit*> { p.Emit(string(buffer[begin:end])) } { p.End() }

Integer    <- { p.Start(INT) } <WholeNum> { p.Emit(string(buffer[begin:end])) } { p.End() }

WholeNum   <- '0' / '-'? [1-9] Digit*

Digit      <- [0-9]

Boolean    <- { p.Start(BOOL) } <'true' / 'false'> { p.Emit(string(buffer[begin:end])) } { p.End() }

Func       <- { p.Start(FUNC) } { p.Start(ARGS) } '(' sp (Name (sp ',' sp Name)* sp)? ')' { p.End() } sp '->' sp (Block / Expr)  { p.End() }

FuncApply  <- { p.Start(FUNCAPPLY) } Name Tuple { p.End() }

List       <- { p.Start(LIST) } '[' sp (Expr (sp ',' sp Expr)* sp)? ']' { p.End() }

Tuple      <- { p.Start(TUPLE) } '(' sp (Expr (sp ',' sp Expr)* sp)? ')' { p.End() }

msp        <- (ws / comment)+

sp         <- (ws / comment)*

comment    <- '#'  (!'\n' .)*

ws         <- [ \t\n\r]
