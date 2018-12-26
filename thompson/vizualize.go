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
	for _, v := range o.items {
		ret = append(ret, v.Render()...)
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
	ret := []string{}
	h := o.h
	for _, v := range o.branches {
		w, _ := v.Size()
		ret = append(ret, v.Render()...)
		for i := 0; i < h; i++ {
			if i > len(ret)-1 {
				ret = append(ret, strings.Repeat("+", w))
				continue
			}
			len := len(ret[i])
			if len < w {
				ret[i] += strings.Repeat("+", w-len)
			}
		}
	}
	return ret
}

func newTerm(text string) *tTerm {
	return &tTerm{name: text}
}

func (o *tTerm) Size() (int, int) {
	return len(o.name), 1
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
	for state != nil && len(state.out) > 0 {
		if len(state.out) == 0 {
			fmt.Printf("!!!!unroll0 len!!!!, %v\n", state)
		}
		ret = append(ret, state)
		state = state.out[0]
	}
	return ret
}

func (o *tViz) buildBranch(state *TState, stopList []*TState) (*tBranch, *TState) {
	if state == nil {
		fmt.Printf("!!!!branch!!!!")
		return nil, nil
	}
	maxW, maxH := 1, 1
	branch := &tBranch{}
	for notInStopList(state, stopList) {
		term := newTerm(state.Name())
		branch.Add(term)

		w, h := term.Size()
		maxW = misc.MaxInt(maxW, w)
		maxH = misc.MaxInt(maxH, h)

		if len(state.out) > 1 {
			split, st := o.buildSplit(state)
			state = st
			branch.Add(split)

			w, h = split.Size()
			maxW = misc.MaxInt(maxW, w)
			maxH = misc.MaxInt(maxH, h)
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
	split.h = maxH
	return split, st
}

func (o *tViz) String() string {
	list := o.head.Render()
	return strings.Join(list, "\n")
}

// func (o *tViz) render(items ...*tItem) string {
// 	str := ""
// 	for _, item := range items {
// 		if item == nil {
// 			str += strings.Repeat(" ", 12)
// 			continue
// 		}
// 		if len(item.branches) == 0 {
// 			str += fmt.Sprintf("%12v", item.name)
// 			continue
// 		}
// 		str += "!!!ERROR!!!"
// 	}
// 	return str
// }
