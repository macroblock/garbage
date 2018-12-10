package main

import (
	"fmt"

	"github.com/macroblock/imed/pkg/ptool"
)

//
//
var parserRules = `
entry           = '' decl {';' decl} $;

decl            = nodeDecl|blockDecl;

nodeDecl        = @lval @options'=' @expr
                | @ERR_incorrect_declaration;
options         = [@optKeepMode] @optSpaceMode [@optRuneSize];
optKeepMode     = '@';
optSpaceMode    = '+'|'-';
optRuneSize     = number;

expr            = exprBody {exprBody};
exprBody        = [@keep]#[@counter]#( @ident | @select | @group | @term );
term            = @range | r | str | @eof;
select          = '[' {@expr} ']';
group           = '(' {@expr} ')';
keepValue       = '<' {@expr} '>';

lval            =  @ident [str] { ',' @ident [str] };
name            = ident;
keep            = '@';

counter         = '+' | '*' | '?'| @rangeCounter;
rangeCounter    =

range           = r '..' r;
r               = \x27 # !\x27 #@rune # \x27;
rune            = anyRune;
str             = \x27 # @string # \x27;
string          = {# !\x27 # anyRune };
ident           = letter#{#letter|digit};
number          = digit#{#digit};

eof             = '$';
digit           = '0'..'9';
letter          = 'a'..'z'|'A'..'Z'|'_';
anyRune         = \x00..\xff;

                = {# ' '| \x09 | \x0D | \x0A };

ERR_incorrect_declaration = '' {# !';' # !$ # anyRune}

`

var parser *ptool.TParser

func init() {
	var err error
	builder := ptool.NewBuilder().FromString(parserRules).Entries("entry")
	parser, err = builder.Build()
	fmt.Println("=================")
	fmt.Println(builder.TreeToString())
	fmt.Println("=================")
	if err != nil {
		fmt.Println("\nparser error: ", err)
		return
	}
}
