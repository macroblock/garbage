package main

import (
	"fmt"
	"strings"

	"github.com/macroblock/imed/pkg/misc"
)

// TBorder -
type TBorder struct {
	LUp   rune
	H1    rune
	RUp   rune
	V1    rune
	Space rune
	V2    rune
	LDn   rune
	H2    rune
	RDn   rune
	HUp   rune
	HDn   rune
}

var borderMap = map[string]TBorder{
	"Default":                  TBorder{'+', '-', '+', '|', ' ', '|', '+', '-', '+', '+', '+'},
	"Single (ASCII)":           TBorder{'+', '~', '+', '|', ' ', '|', '+', '~', '+', '+', '+'},
	"Double (ASCII)":           TBorder{'#', '=', '#', 'N', ' ', 'N', '#', '=', '#', '#', '#'},
	"Single":                   TBorder{'┌', '─', '┐', '│', ' ', '│', '└', '─', '┘', '┰', '┸'},
	"Single (rounded)":         TBorder{'╭', '─', '╮', '│', ' ', '│', '╰', '─', '╯', '┰', '┸'},
	"Double":                   TBorder{'╔', '═', '╗', '║', ' ', '║', '╚', '═', '╝', '╩', '╩'},
	"Shadowed (mix)":           TBorder{'┌', '─', '╖', '│', ' ', '║', '╘', '═', '╝', '+', '+'},
	"Solid (full block)":       TBorder{'█', '█', '█', '█', ' ', '█', '█', '█', '█', '+', '+'},
	"Solid (inner half block)": TBorder{'▄', '▄', '▄', '█', ' ', '█', '▀', '▀', '▀', '+', '+'},
	"Solid (outer half block)": TBorder{'█', '▀', '█', '█', '·', '█', '█', '▄', '█', '+', '+'},
}

var border = borderMap["Single"]

const (
	cPlaceLUp = iota
	cPlaceMUp
	cPlaceRUp
	cPlaceLMid
	cPlaceMMid
	cPlaceRMid
	cPlaceLDn
	cPlaceMDn
	cPlaceRDn
)

func placeText(width int, centered bool, str string) string {
	total := width - len(str)
	half := (width-len(str))/2 - 1
	preGap := misc.MaxInt(0, half)
	postGap := misc.MaxInt(0, total-half)
	return fmt.Sprint(
		strings.Repeat(string(border.Space), preGap),
		str,
		strings.Repeat(string(border.Space), postGap))
}

func whereToIndex(whereX, whereY int) int {
	x := 0
	switch {
	case whereX < 0:
		x = -1
	case whereX > 0:
		x = 1
	}
	y := 0
	switch {
	case whereY < 0:
		y = -1
	case whereY > 0:
		y = 1
	}
	return (x + 1) + (y+1)*3
}

func placeLine(whereX, whereY int, width int, centered bool) string {
	preGap := misc.MaxInt(0, width/2-1)
	postGap := misc.MaxInt(0, width-width/2)
	idx := whereToIndex(whereX, whereY)
	switch idx {
	case cPlaceLUp:
		return fmt.Sprint(
			strings.Repeat(string(border.Space), preGap),
			string(border.LUp),
			strings.Repeat(string(border.H1), postGap))
	case cPlaceMUp:
		return fmt.Sprint(
			strings.Repeat(string(border.H1), preGap),
			string(border.HUp),
			strings.Repeat(string(border.H1), postGap))
	case cPlaceRUp:
		return fmt.Sprint(
			strings.Repeat(string(border.H1), preGap),
			string(border.RUp),
			strings.Repeat(string(border.Space), postGap))
	case cPlaceLDn:
		return fmt.Sprint(
			strings.Repeat(string(border.Space), preGap),
			string(border.LDn),
			strings.Repeat(string(border.H2), postGap))
	case cPlaceMDn:
		return fmt.Sprint(
			strings.Repeat(string(border.H2), preGap),
			string(border.HDn),
			strings.Repeat(string(border.H2), postGap))
	case cPlaceRDn:
		return fmt.Sprint(
			strings.Repeat(string(border.H2), preGap),
			string(border.RDn),
			strings.Repeat(string(border.Space), postGap))
	case cPlaceMMid:
		return fmt.Sprint(
			strings.Repeat(string(border.Space), preGap),
			string(border.V1),
			strings.Repeat(string(border.Space), postGap))
	}
	return fmt.Sprintf("an unsupported place index %v", idx)
}
