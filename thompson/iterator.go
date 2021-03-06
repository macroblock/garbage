package main

import (
	"fmt"

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

const (
	errInvalidObject    = " invalid iterator"
	errInvalidOperation = " invalid operation"
)

// NewNodeIterator -
func NewNodeIterator(root *ptool.TNode, parser *ptool.TParser) *TNodeIterator {
	if root == nil || 0 > len(root.Links)-1 || parser == nil {
		log.Error(true, "invalid parameteres")
		return nil
	}
	o := &TNodeIterator{root: root, ignore: []int{}, parser: parser}
	o.inNode = o.root.Links[o.index]
	return o
}

// Clone -
func (o *TNodeIterator) Clone() *TNodeIterator {
	if o == nil || o.root == nil {
		log.Error(true, errInvalidObject)
		return nil
	}
	ret := &TNodeIterator{}
	*ret = *o
	return ret
}

// Enter -
func (o *TNodeIterator) Enter() *TNodeIterator {
	if o == nil || o.root == nil {
		log.Error(true, errInvalidObject)
		return nil
	}
	if o.retNode == nil {
		log.Error(true, "invalid operation")
		return nil
	}
	o.root = o.retNode
	o.retNode = nil
	o.index = 0
	if o.index > len(o.root.Links)-1 {
		o.root = nil
		log.Error(true, "node has no links")
		return nil
	}
	o.inNode = o.root.Links[o.index]
	return o
}

// AutoIgnore -
func (o *TNodeIterator) AutoIgnore(items ...interface{}) *TNodeIterator {
	if o == nil || o.root == nil {
		log.Error(true, errInvalidObject)
		return nil
	}
	ok := false
	o.ignore, ok = itemsToIntSlice(o.parser, items...)
	if !ok {
		o.root = nil
		return nil
	}
	return o
}

// Try -
func (o *TNodeIterator) Try(items ...interface{}) bool {
	it := &TNodeIterator{}
	*it = *o
	err := it.tryToAccept(items...)
	if err != nil {
		return false
	}
	*o = *it
	return true
}

// Accept -
func (o *TNodeIterator) Accept(items ...interface{}) *TNodeIterator {
	err := o.tryToAccept(items...)
	log.Error(err)
	return o
}

func (o *TNodeIterator) tryToAccept(items ...interface{}) error {
	// retrieveCallInfo()
	if o == nil || o.root == nil {
		return fmt.Errorf(errInvalidObject)
	}
	if o.inNode == nil {
		o.root = nil
		return fmt.Errorf("no more nodes")
	}

	if len(items) > 0 {
		where, ok := itemsToIntSlice(o.parser, items...)
		if !ok {
			o.root = nil
			return nil
		}
		// fmt.Printf("id: %v; where: %v\n", o.inNode.Type, where)
		// log.Notice(where, "type:", o.inNode.Type)
		if !in(o.inNode.Type, where) {
			o.root = nil
			// log.Info("exit 0")
			return fmt.Errorf("unexpected node %q", o.parser.ByID(o.inNode.Type))
		}
	}

	o.index++
	if o.index > len(o.root.Links)-1 {
		if o.inNode == nil {
			o.root = nil
			// log.Info("exit 1")
			return fmt.Errorf("index %v out of range 0..%v", o.index, len(o.root.Links)-1)
		}
		o.retNode = o.inNode
		o.inNode = nil
		// log.Info("exit 2")
		return nil
	}
	o.retNode = o.inNode
	o.inNode = o.root.Links[o.index]
	// log.Info("exit 3", o.inNode.Type)
	return nil
}

// Node -
func (o *TNodeIterator) Node() *ptool.TNode {
	if o == nil || o.root == nil {
		log.Error(true, errInvalidObject)
		return nil
	}
	return o.retNode
}

// Value -
func (o *TNodeIterator) Value() string {
	if o == nil || o.root == nil {
		log.Error(true, errInvalidObject)
		return ""
	}
	if o.retNode == nil {
		log.Error(true, errInvalidOperation)
		return ""
	}
	return o.retNode.Value
}

// ID -
func (o *TNodeIterator) ID() int {
	if o == nil || o.root == nil {
		log.Error(true, errInvalidObject)
		return -1
	}
	if o.retNode == nil {
		log.Error(true, errInvalidOperation)
		return -1
	}
	return o.retNode.Type
}

// Name -
func (o *TNodeIterator) Name() string {
	if o == nil || o.root == nil {
		log.Error(true, errInvalidObject)
		return ""
	}
	if o.retNode == nil {
		log.Error(true, errInvalidOperation)
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

func itemToInt(p *ptool.TParser, item interface{}) int {
	ret := -1
	switch t := item.(type) {
	default:
		log.Errorf(true, "illegal type %v", t)
	case int:
		ret = t
		log.Errorf(ret < 0, "invalid node id %v", t)
	case string:
		ret = p.ByName(t)
		log.Errorf(ret < 0, "unknown node type %q", t)
	}
	return ret
}

func itemsToIntSlice(p *ptool.TParser, items ...interface{}) ([]int, bool) {
	ret := []int{}
	ok := true
	for _, i := range items {
		id := itemToInt(p, i)
		if id < 0 {
			ok = false
			continue
		}
		ret = append(ret, id)
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
