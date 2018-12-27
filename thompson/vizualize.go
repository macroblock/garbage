package main

import (
	"fmt"
	"strings"

	"github.com/macroblock/imed/pkg/misc"
)

const (
	typeError = iota
	typeSplit
	typeSeq
)

type (
	tViz struct {
		head *tBranch
	}

	// IItem -
	IItem interface {
		Render() []string
	}

	tTerm struct {
		name string
	}

	tBranch struct {
		items []IItem
		w, h  int
	}

	tSplit struct {
		branches []*tBranch
		w, h     int
	}
)

func (o *tBranch) Add(v IItem) {
	o.items = append(o.items, v)
}

func (o *tBranch) Size() (int, int) {
	return o.w, o.h
}

func (o *tBranch) Render() []string {
	ret := []string{}
	w := o.w
	for _, v := range o.items {
		for _, s := range v.Render() {
			// total := w - len(s)
			// half := (w - len(s)) / 2
			// preGap := misc.MaxInt(0, half)
			// postGap := misc.MaxInt(0, total-half)
			ret = append(ret, placeText(w, true, s))
			// strings.Repeat(".", preGap)+s+ // fmt.Sprintf("-%vx%v", o.w, o.h)+
			// strings.Repeat(".", postGap),
			// )
		}
	}
	return ret
}

func (o *tSplit) Add(v *tBranch) {
	o.branches = append(o.branches, v)
}

func (o *tSplit) Size() (int, int) {
	return o.w, o.h
}
func (o *tSplit) Render() []string {
	h := o.h
	ret := make([]string, o.h)
	n := o.h - 1
	for i, v := range o.branches {
		w, _ := v.Size()
		// preGap := misc.MaxInt(0, w/2-1)
		// postGap := misc.MaxInt(0, w-w/2)
		switch i {
		default:
			// ret[0] += strings.Repeat("-", preGap) + "V" + strings.Repeat("-", postGap)
			// ret[n] += strings.Repeat("-", preGap) + "A" + strings.Repeat("-", postGap)
			ret[0] += placeLine(0, -1, w, true)
			ret[n] += placeLine(0, 1, w, true)
		case 0:
			// ret[0] += strings.Repeat(".", preGap) + "/" + strings.Repeat("-", postGap)
			// ret[n] += strings.Repeat(".", preGap) + "\\" + strings.Repeat("-", postGap)
			ret[0] += placeLine(-1, -1, w, true)
			ret[n] += placeLine(-1, 1, w, true)
		case len(o.branches) - 1:
			// ret[0] += strings.Repeat("-", preGap) + "\\" + strings.Repeat(".", postGap)
			// ret[n] += strings.Repeat("-", preGap) + "/" + strings.Repeat(".", postGap)
			ret[0] += placeLine(1, -1, w, true)
			ret[n] += placeLine(1, 1, w, true)
		}
		br := v.Render()
		for i := 0; i < h-2; i++ {
			if i < len(br) {
				ret[i+1] += br[i]
				continue
			}
			// ret[i+1] += strings.Repeat(".", preGap) + "|" + strings.Repeat(".", postGap)
			ret[i+1] += placeLine(0, 0, w, true)
		}
	}
	return ret
}

func newTerm(text string) *tTerm {
	return &tTerm{name: text}
}

func (o *tTerm) Size() (int, int) {
	return len(o.name) + 2, 1
}

func (o *tTerm) Render() []string {
	ret := []string{o.name}
	return ret
}

func newViz(head *TState) *tViz {
	viz := &tViz{}
	viz.head = &tBranch{}
	viz.head, _ = viz.buildBranch(head, nil)
	return viz
}

func notInStopList(state *TState, stopList []*TState) bool {
	if state == nil {
		return false
	}
	for _, v := range stopList {
		if v == state {
			return false
		}
	}
	return true
}

func unroll0(state *TState) []*TState {
	if state == nil {
		fmt.Printf("!!!!unroll0!!!!")
		return nil
	}
	ret := []*TState{}
	for state != nil {
		ret = append(ret, state)
		if len(state.out) == 0 {
			// fmt.Printf("!!!!unroll0 len!!!!, %v\n", state)
			state = nil
			continue
		}
		state = state.out[0]
	}
	return ret
}

func (o *tViz) buildBranch(state *TState, stopList []*TState) (*tBranch, *TState) {
	if state == nil {
		fmt.Printf("!!!!branch!!!!")
		return nil, nil
	}
	maxW, maxH := 0, 0
	branch := &tBranch{}
	fmt.Println("stop list", stopList)
	for notInStopList(state, stopList) {
		fmt.Println("    state ", state.Name())
		term := newTerm(state.Name())
		branch.Add(term)

		w, h := term.Size()
		maxW = misc.MaxInt(maxW, w)
		maxH += h

		if len(state.out) > 1 {
			split, st := o.buildSplit(state)
			state = st
			branch.Add(split)

			w, h = split.Size()
			maxW = misc.MaxInt(maxW, w)
			maxH += h
			continue
		}
		if len(state.out) == 0 {
			state = nil
			continue
		}
		state = state.out[0]
	}
	branch.w = maxW
	branch.h = maxH
	fmt.Println("<<<<<< out")
	return branch, state
}

func (o *tViz) buildSplit(state *TState) (*tSplit, *TState) {
	if state == nil {
		fmt.Printf("!!!!split!!!!")
		return nil, nil
	}
	if len(state.out) < 2 {
		fmt.Printf("!!!!split len!!!!")
		return nil, nil
	}
	unrolled := unroll0(state.out[0])
	fmt.Println("     ----unrolled", unrolled)
	split := &tSplit{}
	br1, st := o.buildBranch(state.out[1], unrolled)
	stopNode := []*TState{st}
	br0, _ := o.buildBranch(state.out[0], stopNode)
	split.Add(br0)
	split.Add(br1)

	maxW, maxH := br0.Size()
	w, h := br1.Size()
	maxW += w
	maxH = misc.MaxInt(maxH, h)

	for _, v := range state.out[2:] {
		br, _ := o.buildBranch(v, stopNode)
		split.Add(br)

		w, h := br.Size()
		maxW += w
		maxH = misc.MaxInt(maxH, h)
	}
	split.w = maxW
	split.h = maxH + 2 // +2 for drawing header and footer
	return split, st
}

func (o *tViz) String() string {
	list := o.head.Render()
	return strings.Join(list, "\n")
}
