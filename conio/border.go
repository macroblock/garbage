package conio

// TBorder -
type TBorder struct {
	LU rune
	H1 rune
	RU rune
	V1 rune
	V2 rune
	LD rune
	H2 rune
	RD rune
}

type tBorderMap struct {
	val map[string]TBorder
}

// Border -
var Border tBorderMap

func initBorder() {
	Border.val = make(map[string]TBorder)
	Border.val["Default"] = TBorder{'+', '-', '+', '|', '|', '+', '-', '+'}
	Border.val["Single (ASCII)"] = TBorder{'+', '~', '+', '|', '|', '+', '~', '+'}
	Border.val["Double (ASCII)"] = TBorder{'#', '=', '#', 'N', 'N', '#', '=', '#'}
	Border.val["Single"] = TBorder{'╭', '─', '╮', '│', '│', '╰', '─', '╯'}
	Border.val["Double"] = TBorder{'╔', '═', '╗', '║', '║', '╚', '═', '╝'}
}

func (b *tBorderMap) Get(name string) TBorder {
	return Border.val[name]
}
