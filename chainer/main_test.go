package chainer

import (
	"testing"

	"github.com/macroblock/imed/pkg/ptool"
)

var (
	keyChainParseCorrect = []string{
		"",
		"asdf",
		"a-d",
		"a-d 0-9",
		"<enter>",
		"<ctrl>a",
		"<alt+ctrl+enter>",
	}
	keyChainParseIncorrect = []string{
		"<>",
		"<shift>",
		"asdf-",
		"<enter+control>",
	}
)

func TestParseCorrect(t *testing.T) {
	for _, v := range keyChainParseCorrect {
		tree, err := parser.Parse(v)
		if err != nil {
			t.Errorf("\n%q\nParse() error: %v\n%v", v, err, ptool.TreeToString(tree, parser.ByID))
			continue
		}
	}
}

func TestParseIncorrect(t *testing.T) {
	for _, v := range keyChainParseIncorrect {
		tree, err := parser.Parse(v)
		if err == nil {
			t.Errorf("\n%q\nParse() has no error: %v\n%v", v, err, ptool.TreeToString(tree, parser.ByID))
			continue
		}
	}
}
