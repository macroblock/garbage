package main

type (
	tViz struct {
		levels []tLevel
	}

	tLevel struct {
		nodes []*tNode
	}

	tNode struct {
		name string
		out  []tLink
		in   []tLink
	}

	tLink struct {
	}
)

func newViz() *tViz {
	viz := &tViz{}
	return viz
}

func (o *tViz) Init(state *TState) {
	vizited := map[*TState]*tNode{}
	curr := []*TState{}
	o.levels = []tLevel{}
	curr = append(curr, state)
	for {
		temp := []*TState{}
	}
}
