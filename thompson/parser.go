package main

import (
	"fmt"

	"github.com/macroblock/imed/pkg/ptool"
)

//
//
var parserRules = `
entry           = '' decl {';' decl} $;

decl            = @lval '=' [@options] @expr [@exec]
                | @ERR_incorrect_declaration;
options         = '!none!';
exec            = '{' expr '}';

expr            = @keepValue | [[@name]#@keep]#[@counter]#( @ident | @or | @and | @term );
term            = @range | r | str | @eof;
or              = '[' {@expr} ']';
and             = '(' {@expr} ')';
keepValue       = '<' {@expr} '>';

lval            = ident;
name            = ident;
keep            = '@';
counter         = '+' | '*' | '?';

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
