package chainer

import (
	"fmt"

	"github.com/macroblock/imed/pkg/ptool"
)

var rule = `
entry       = {@runeSet|@keySet|@range} $;

keySet      =  '<' keyExpr '>';
runeSet     =  '[' runeExpr {runeExpr} ']';
range       = @rune '-' @rune;

keyExpr    = ( @mod [keySeparator keyExpr] ) | @key;
keySeparator= '+'|'-';

runeExpr    = @range | @rune;

rune        = !'<'!'>'!'['!']'!'-' \x21..\xfe;

key         = 'enter'|'esc'|'f1'|'space';
mod         = 'shift'|'alt'|'ctrl';

    = { ' ' | \x09 | \x0a | \x0d };
`

var parser *ptool.TParser

func init() {
	p, err := ptool.NewBuilder().FromString(rule).Entries("entry").Build()
	if err != nil {
		fmt.Println("\nparser error: ", err)
		panic("")
	}
	parser = p
}
