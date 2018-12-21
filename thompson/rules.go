package main

import (
	"fmt"

	"github.com/macroblock/imed/pkg/ptool"
)

var parser *ptool.TParser

var parserRules = `
entry           = '' {';'} (decl|useBelowMode) {decl|useBelowMode} $;

decl            = (@nodeDecl|@blockDecl) ';' {';'};

nodeDecl        = [useMode] @lval '=' @sequence
                | @ERR_incorrect_declaration;
options         = [@optKeepMode] [@optSpaceMode] [@optRuneSize];
optKeepMode     = '@';
optSpaceMode    = '+'|'-';
optRuneSize     = number;

blockDecl       = [useMode] @lval '=>' @options '{' decl '}';

sequence        = expr {expr};
expr            = repeat(@ident | @keepNode | @split | seq | @keepValue | term );
term            = @range | r | str | @eof;
seq             = '(' {@sequence} ')';
split           = '[' {sequence} ']';
keepValue       = '<' {sequence} '>';
keepNode        = '@' # @ident;

lval            = [useMode] nodeVar|sysVar;
nodeVar         = @ident [str];
sysVar          = @sysIdent [str];

useBelowMode    = (@useBelowOn|@useBelowOff) {';'};
useMode         = @useOn|@useOff|@useExclude;
useBelowOn      = '++:';
useBelowOff     = '--:';
useExclude      = '***';
useOn           = '+++';
useOff          = '---';

repeat          = [ @repeat_01 | @repeat_0f | @repeat_1f | @repeat_xy | @repeat_xf | @repeat_x ];
repeat_01       = '?';
repeat_0f       = '*' # [@lazy];
repeat_1f       = '+' # [@lazy];
repeat_xy       = @number # '-' # @number # [@lazy];
repeat_xf       = @number # ('-'|'+') # [@lazy];
repeat_x        = @number # [@lazy];
lazy            = '?'; 

comment         = '//' # {# !\x0a # !$ # anyRune };

range           = r # '-' # r;
r               = (\x27 # !\x27 #@rune # \x27) | ('\x'#@runeCode);
rune            = anyRune;
runeCode        = hexDigit hexDigit;
str             = '"' # @string1 # '"' | \x27 # @string2 # \x27;
string1         = {# !'"' # anyRune };
string2         = {# !\x27 # anyRune };
ident           = letter#{#letter|digit};
sysIdent        = '$'#letter#{#letter|digit};
number          = digit#{#digit};

eof             = '$eof' | '$EOF';
hexDigit        = '0'..'9' | 'a'..'f' | 'A'..'F';
digit           = '0'..'9';
letter          = 'a'..'z'|'A'..'Z'|'_';
anyRune         = \x00..\xff;

                = {# ' ' | \x0a | comment | \x09 | \x0d  };

ERR_incorrect_declaration = '' '### OFF ###' {# !';' # !$ # anyRune}

`

func init() {
	var err error
	builder := ptool.NewBuilder().FromString(parserRules).Entries("entry")
	parser, err = builder.Build()
	fmt.Println("=================")
	fmt.Println(builder.TreeToString())
	fmt.Println("=================")
	if err != nil {
		fmt.Println("\nparser initialization error: ", err)
		return
	}
}
