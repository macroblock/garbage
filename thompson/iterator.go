package main

import (
	"fmt"

	"github.com/macroblock/imed/pkg/ptool"
)

// TNodeIterator -
type TNodeIterator struct {
	root   *ptool.TNode
	node   *ptool.TNode
	idx    int
	parser *ptool.TParser
}

// NewNodeIterator -
func NewNodeIterator(root *ptool.TNode, parser *ptool.TParser) *TNodeIterator {
	if root == nil || 0 > len(root.Links)-1 {
		return nil
	}
	o := &TNodeIterator{root: root, parser: parser}
	o.node = o.root.Links[o.idx]
	return o
}

// Clone -
func (o *TNodeIterator) Clone() *TNodeIterator {
	if o == nil || o.root == nil {
		return nil
	}
	ret := &TNodeIterator{}
	*ret = *o
	return ret
}

// Enter -
func (o *TNodeIterator) Enter() *TNodeIterator {
	if o == nil || o.root == nil {
		return nil
	}
	o.root = o.node
	o.idx = 0
	if o.idx > len(o.root.Links)-1 {
		o.root = nil
		return nil
	}
	o.node = o.root.Links[o.idx]
	return o
}

// ForEach -
func (o *TNodeIterator) ForEach(fn func(*TNodeIterator)) *TNodeIterator {
	if o == nil || o.root == nil {
		return nil
	}
	temp := &TNodeIterator{}
	for it := o.Clone().Enter(); it.Node() != nil; it.Next() {
		*temp = *it
		fn(temp)
		// fmt.Printf("\niterator: %v\nindex: %v\n%v\n", it, it.idx, ptool.TreeToString(it.node, it.parser.ByID))
	}
	return o
}

// Find -
func (o *TNodeIterator) Find(items ...interface{}) *TNodeIterator {
	it := o.Clone()
	it.idx = 0
	return it.Next(items...)
}

// Next -
func (o *TNodeIterator) Next(items ...interface{}) *TNodeIterator {
	if o == nil || o.root == nil {
		return nil
	}

	// when no arguments just step to the next node
	if len(items) == 0 {
		o.idx++
		if o.root == nil || o.idx > len(o.root.Links)-1 {
			o.root = nil
			return nil
		}
		o.node = o.root.Links[o.idx]
		return o
	}

	// attempt to find the next node that satisfies an items path
	root := o.root
	idx := o.idx
	if idx > 0 {
		idx++
	}
	node := root
	i := idx
	for _, item := range items {
		tid := -1
		switch t := item.(type) {
		default:
			panic("illegal type")
		case int:
		case string:
			tid = o.parser.ByName(t)
		}
		root = node
		node = nil
		for i < len(root.Links) {
			v := root.Links[i]
			fmt.Printf(":: %v ~ %v -- %v\n", o.parser.ByID(tid), o.parser.ByID(v.Type), v.Value)
			if v.Type == tid {
				node = v
				idx = i
				break
			}
			i++
		}
		i = 0
		if node == nil {
			o.root = nil
			fmt.Println("-- end nil")
			return nil
		}
	}
	o.root = root
	o.node = node
	o.idx = idx
	fmt.Println("++ end ok")
	return o
}

// Node -
func (o *TNodeIterator) Node() *ptool.TNode {
	if o == nil || o.root == nil {
		return nil
	}
	return o.node
}

// Value -
func (o *TNodeIterator) Value() string {
	if o == nil || o.root == nil {
		return ""
	}
	return o.node.Value
}

// NodeID -
func (o *TNodeIterator) NodeID() int {
	if o == nil || o.root == nil {
		return -1
	}
	return o.node.Type
}

// NodeName -
func (o *TNodeIterator) NodeName() string {
	if o == nil || o.root == nil {
		return ""
	}
	return o.parser.ByID(o.node.Type)
}
