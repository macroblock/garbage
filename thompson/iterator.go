package main

import (
	"github.com/macroblock/imed/pkg/ptool"
)

// TNodeIterator -
type TNodeIterator struct {
	root    *ptool.TNode
	inNode  *ptool.TNode
	retNode *ptool.TNode
	index   int
	ignore  []int
	parser  *ptool.TParser
}

// NewNodeIterator -
func NewNodeIterator(root *ptool.TNode, parser *ptool.TParser) *TNodeIterator {
	if root == nil || 0 > len(root.Links)-1 || parser == nil {
		log.Error(true, "NewNodeIterator: invalid parameteres")
		return nil
	}
	o := &TNodeIterator{root: root, ignore: []int{}, parser: parser}
	o.inNode = o.root.Links[o.index]
	return o
}

// Clone -
func (o *TNodeIterator) Clone() *TNodeIterator {
	if o == nil || o.root == nil {
		log.Error(true, "Clone: invalid object")
		return nil
	}
	ret := &TNodeIterator{}
	*ret = *o
	return ret
}

// Enter -
// func (o *TNodeIterator) Enter() *TNodeIterator {
// 	if o == nil || o.root == nil {
// 		log.Error(true, "Enter: invalid object")
// 		return nil
// 	}
// 	o.root = o.node
// 	o.index = 0
// 	if o.index > len(o.root.Links)-1 {
// 		o.root = nil
// 		log.Error(true, "Enter: node has no links")
// 		return nil
// 	}
// 	o.node = o.root.Links[o.index]
// 	return o
// }

// AutoIgnore -
func (o *TNodeIterator) AutoIgnore(items ...interface{}) *TNodeIterator {
	if o == nil || o.root == nil {
		log.Error(true, "Ignore: invalid object")
		return nil
	}
	ok := false
	o.ignore, ok = itemsToIntSlice("Ignore", o.parser, items...)
	if !ok {
		o.root = nil
		return nil
	}
	return o
}

// Accept -
func (o *TNodeIterator) Accept(items ...interface{}) *TNodeIterator {
	if o == nil || o.root == nil || o.inNode == nil {
		log.Error(true, "Accept: invalid object")
		o.root = nil
		return nil
	}

	if len(items) > 0 {
		where, ok := itemsToIntSlice("Accept", o.parser, items...)
		if !ok {
			o.root = nil
			return nil
		}
		// fmt.Printf("id: %v; where: %v\n", o.inNode.Type, where)
		if !in(o.inNode.Type, where) {
			o.root = nil
			log.Errorf(true, "Accept: unexpected node %q", o.parser.ByID(o.inNode.Type))
			// return nil
		}
	}

	o.index++
	if o.index > len(o.root.Links)-1 {
		if o.inNode == nil {
			log.Errorf(true, "Accept: index %v out of range 0..%v", o.index, len(o.root.Links)-1)
			o.root = nil
			return nil
		}
		o.retNode = o.inNode
		o.inNode = nil
		return o
	}
	o.retNode = o.inNode
	o.inNode = o.root.Links[o.index]
	return o
}

// Node -
func (o *TNodeIterator) Node() *ptool.TNode {
	if o == nil || o.root == nil {
		log.Error(true, "Node: invalid object")
		return nil
	}
	if o.retNode == nil {
		log.Error(true, "Node: invalid operation")
		return nil
	}
	return o.retNode
}

// Value -
func (o *TNodeIterator) Value() string {
	if o == nil || o.root == nil {
		return ""
	}
	if o.retNode == nil {
		log.Error(true, "Value: invalid operation")
		return ""
	}
	return o.retNode.Value
}

// ID -
func (o *TNodeIterator) ID() int {
	if o == nil || o.root == nil {
		return -1
	}
	if o.retNode == nil {
		log.Error(true, "ID: invalid operation")
		return -1
	}
	return o.retNode.Type
}

// Name -
func (o *TNodeIterator) Name() string {
	if o == nil || o.root == nil {
		return ""
	}
	if o.retNode == nil {
		log.Error(true, "Name: invalid operation")
		return ""
	}
	return o.parser.ByID(o.retNode.Type)
}

func in(what int, where []int) bool {
	for _, v := range where {
		if what == v {
			return true
		}
	}
	return false
}

func itemToInt(errPrefix string, parser *ptool.TParser, item interface{}) (int, bool) {
	ret := -1
	switch t := item.(type) {
	default:
		log.Errorf(true, "%v: illegal type %v", errPrefix, t)
	case int:
		ret = t
	case string:
		ret := parser.ByName(t)
		if ret < 0 {
			log.Errorf(true, "%v: unknown node type %q", errPrefix, t)
			return ret, false
		}
	}
	return ret, true
}

func itemsToIntSlice(errPrefix string, parser *ptool.TParser, items ...interface{}) ([]int, bool) {
	ret := []int{}
	ok := true
	for _, i := range items {
		switch t := i.(type) {
		default:
			log.Errorf(true, "%v: illegal type %v", errPrefix, t)
		case int:
			ret = append(ret, t)
		case string:
			id := parser.ByName(t)
			if id < 0 {
				ok = false
				log.Errorf(true, "%v: unknown node type %q", errPrefix, t)
				continue
			}
			ret = append(ret, id)
		}
	}
	return ret, ok
}

// ForEach -
// func (o *TNodeIterator) ForEach(fn func(*TNodeIterator)) *TNodeIterator {
// 	if o == nil || o.root == nil {
// 		return nil
// 	}
// 	temp := &TNodeIterator{}
// 	for it := o.Clone().Enter(); it.Node() != nil; it.Next() {
// 		*temp = *it
// 		fn(temp)
// 		// fmt.Printf("\niterator: %v\nindex: %v\n%v\n", it, it.idx, ptool.TreeToString(it.node, it.parser.ByID))
// 	}
// 	return o
// }
