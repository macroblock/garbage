package main

import (
	"fmt"

	"github.com/macroblock/garbage/zlog"
)

var (
	log = zlog.Instance()
)

var testProg = `
+++ entry01 = x [a b] z;
// +++ entry02 = sysIdent;

a = 'f';
b = 'g';
x = 'x';
z = 'z';

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
	frag, errors := Thompson(elem)
	log.Notice(errors, "thompson error(s):")
	if frag != nil {
		log.Notice(frag.first, "result2:")
	}
}

func print(str string, errors ...error) {
	fmt.Println(str)
	for _, v := range errors {
		fmt.Println(v)
	}
}
