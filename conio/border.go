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
	BorderMap.val = map[string]TBorder{
		"Default":                  TBorder{'+', '-', '+', '|', '|', '+', '-', '+'},
		"Single (ASCII)":           TBorder{'+', '~', '+', '|', '|', '+', '~', '+'},
		"Double (ASCII)":           TBorder{'#', '=', '#', 'N', 'N', '#', '=', '#'},
		"Single":                   TBorder{'┌', '─', '┐', '│', '│', '└', '─', '┘'},
		"Single (rounded)":         TBorder{'╭', '─', '╮', '│', '│', '╰', '─', '╯'},
		"Double":                   TBorder{'╔', '═', '╗', '║', '║', '╚', '═', '╝'},
		"Shadowed (mix)":           TBorder{'┌', '─', '╖', '│', '║', '╘', '═', '╝'},
		"Solid (full block)":       TBorder{'█', '█', '█', '█', '█', '█', '█', '█'},
		"Solid (inner half block)": TBorder{'▄', '▄', '▄', '█', '█', '▀', '▀', '▀'},
		"Solid (outer half block)": TBorder{'█', '▀', '█', '█', '█', '█', '▄', '█'},
	}
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

// AddBorder -
func (bm *tBorderMap) Add(name string, border TBorder) {
	BorderMap.val[name] = border
}

// DeleteBorder -
func (bm *tBorderMap) Delete(name string) {
	delete(BorderMap.val, name)
}
