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

// BorderMap -
var BorderMap tBorderMap

func initBorder() {
	BorderMap.val = make(map[string]TBorder)
	BorderMap.val["Default"] = TBorder{'+', '-', '+', '|', '|', '+', '-', '+'}
	BorderMap.val["Single (ASCII)"] = TBorder{'+', '~', '+', '|', '|', '+', '~', '+'}
	BorderMap.val["Double (ASCII)"] = TBorder{'#', '=', '#', 'N', 'N', '#', '=', '#'}
	BorderMap.val["Single"] = TBorder{'┌', '─', '┐', '│', '│', '└', '─', '┘'}
	BorderMap.val["Single (rounded)"] = TBorder{'╭', '─', '╮', '│', '│', '╰', '─', '╯'}
	BorderMap.val["Double"] = TBorder{'╔', '═', '╗', '║', '║', '╚', '═', '╝'}
	BorderMap.val["Shadowed (mix)"] = TBorder{'┌', '─', '╖', '│', '║', '╘', '═', '╝'}
	BorderMap.val["Solid (full block)"] = TBorder{'█', '█', '█', '█', '█', '█', '█', '█'}
	BorderMap.val["Solid (inner half block)"] = TBorder{'▄', '▄', '▄', '█', '█', '▀', '▀', '▀'}
	BorderMap.val["Solid (outer half block)"] = TBorder{'█', '▀', '█', '█', '█', '█', '▄', '█'}
}

// Get -
func (bm *tBorderMap) Get(name string) TBorder {
	b, ok := bm.val[name]
	if !ok {
		b = bm.val["Default"]
	}
	return b
}

// Names -
func (bm *tBorderMap) Names() []string {
	keys := make([]string, 0, len(bm.val))
	for k := range BorderMap.val {
		keys = append(keys, k)
	}
	return keys
}
