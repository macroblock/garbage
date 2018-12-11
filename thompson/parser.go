package main

import (
	"fmt"

	"github.com/macroblock/imed/pkg/ptool"
)

//
//
var parserRules = `
entry           = '' {';'} decl { ';'{';'} decl} {';'} $;

decl            = nodeDecl|blockDecl;

nodeDecl        = @lval @options'=' @sequence
                | @ERR_incorrect_declaration;
options         = [@optKeepMode] [@optSpaceMode] [@optRuneSize];
optKeepMode     = '@';
optSpaceMode    = '+'|'-';
optRuneSize     = number;

blockDecl       = @lval @options '=' str '{' decl '}';

sequence        = expr {expr};
expr            = ( repeat| @ident | @keepNode | @select | seq | @keepValue | @term );
term            = @range | r | str | @eof;
select          = '[' {sequence} ']';
seq             = '(' {@sequence} ')';
keepValue       = '<' {sequence} '>';
keepNode        = '@' # @ident;

lval            =  @ident [str] { ',' @ident [str] };
name            = ident;

repeat          = @repeat_01 | @repeat_0f | @repeat_1f | @repeat_xy | @repeat_xf | @repeat_x;
repeat_01       = '?' # repeatList;
repeat_0f       = '*' # repeatList;
repeat_1f       = '+' # repeatList;
repeat_xy       = @number # '-' # @number # repeatList;
repeat_xf       = @number # ('-'|'+') # repeatList;
repeat_x        = @number # repeatList;
repeatList      = @ident | @keepNode | @select | seq | @keepValue | @term;

comment         = '//' # {# !\x0a # !$ # anyRune };

range           = r # '-' # r;
r               = \x27 # !\x27 #@rune # \x27;
rune            = anyRune;
str             = \x27 # @string # \x27;
string          = {# !\x27 # anyRune };
ident           = letter#{#letter|digit};
number          = digit#{#digit};

eof             = '$eof' | '$EOF';
digit           = '0'..'9';
letter          = 'a'..'z'|'A'..'Z'|'_';
anyRune         = \x00..\xff;

                = {# ' ' | \x0a | @comment | \x09 | \x0d  };

ERR_incorrect_declaration = '' '### OFF ###' {# !';' # !$ # anyRune}

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
