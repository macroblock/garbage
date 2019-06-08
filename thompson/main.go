package main

import (
	"fmt"

	"github.com/macroblock/garbage/zlog"
)

var (
	log = zlog.Instance()
)

var testProg = `
+++ entry01 = x ['1234' ('4321' 'abcd' [ 'x' 'y' 'z' ] 'aabcc' '0000' ) ( ['left' 'center' 'right' ] 'aaaaaaa' 'bbb')] zz ;
// [ '9' '8' '7'] [ '3' '2']a b c;
// +++ entry02 = sysIdent;
seq = '1' '2' '3';
a = 'a';
b = 'b';
c = 'c';
d = 'd';
x = 'x';
z = 'z';
zz = 'zz';

// comment here
// y "test" = unk id;

// string          = [('"'<*?anyRune>'"') ("'"<*? anyRune>"'")];

// dent            = letter *[letter digit];
// sysIdent        = '$' letter *[letter digit];
// number          = +digit;

// eof             = ['$eof' '$EOF'];
// hexDigit        = ['0'-'9' 'a'-'f' 'A'-'F'];
// digit           = '0'-'9';
// letter          = ['a'-'z' 'A'-'Z' '_'];
// anyRune         = \x00-\xfe;
`

func main() {
	// log.Add(misc.NewAnsiLogger(loglevel.All, ""))

	if parser == nil {
		return
	}
	// tree, err := parser.Parse(testProg)
	// if err != nil {
	// 	fmt.Println("\n*TParser.Parse error: ", err)
	// 	return
	// }
	// fmt.Println(ptool.TreeToString(tree, parser.ByID))

	parser := NewParser()
	errors := parser.Parse(testProg)
	log.Notice(errors, "parse error(s):")
	elem, errors := parser.Build()
	log.Notice(errors, "build error(s):")
	log.Notice(elem, "result:")
	parser.Run("abcdefg")
	// frag, errors := Thompson(elem)
	// log.Notice(errors, "thompson error(s):")
	// if frag != nil {
	// 	log.Notice(frag.first, "result2:")
	// 	viz := newViz(frag.first)
	// 	_ = viz
	// 	fmt.Println(viz.String())
	// }
}

func print(str string, errors ...error) {
	fmt.Println(str)
	for _, v := range errors {
		fmt.Println(v)
	}
}
