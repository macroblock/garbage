package iterator

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

// New -
func New(root *ptool.TNode, parser *ptool.TParser) *TNodeIterator {
	ret := &TNodeIterator{root: root, idx: -1, parser: parser}
	ret.first()
	// fmt.Println("idx: ", ret.idx)
	return ret
}

// First -
func (o *TNodeIterator) first() *TNodeIterator {
	if o == nil {
		return nil
	}
	o.idx = 0
	// fmt.Println("a1")
	if o.root == nil || o.idx > len(o.root.Links)-1 {
		// fmt.Println("a2")
		o.idx, o.node = -1, nil
		return o
	}
	// fmt.Println("a3")
	o.node = o.root.Links[o.idx]
	// fmt.Println("idx: ", o.idx)
	// fmt.Println(o.node)
	return o
}

// Next -
func (o *TNodeIterator) Next() *TNodeIterator {
	if o == nil || o.idx < 0 {
		return o
	}
	o.idx++
	if o.root == nil || o.idx > len(o.root.Links)-1 {
		o.idx, o.node = -1, nil
		return o
	}
	o.node = o.root.Links[o.idx]
	return o
}

// Enter -
func (o *TNodeIterator) Enter() *TNodeIterator {
	if o == nil || o.node == nil {
		return nil
	}
	o.root = o.node
	o.node = nil
	o.idx = -1
	o.first()
	return o
}

// FindFirst -
func (o *TNodeIterator) findFirst(item interface{}) *TNodeIterator {
	if o == nil {
		return nil
	}
	o.idx, o.node = find(o.root, 0, o.parser, item)
	return o
}

// FindNext -
func (o *TNodeIterator) FindNext(item interface{}) *TNodeIterator {
	if o == nil || o.idx < 0 ||
		o.root == nil || o.idx > len(o.root.Links)-1 {
		return o
	}
	o.idx, o.node = find(o.root, o.idx, o.parser, item)
	return o
}

// ForEach -
func (o *TNodeIterator) ForEach(fn func(*TNodeIterator)) *TNodeIterator {
	if o == nil || o.node == nil {
		return o
	}
	temp := &TNodeIterator{}
	for it := New(o.node, o.parser); it.Node() != nil; it.Next() {
		// fmt.Printf("\niterator: %v\nindex: %v\n%v\n", it, it.idx, ptool.TreeToString(it.node, it.parser.ByID))
		*temp = *it
		fn(temp)
	}
	return o
}

// Node -
func (o *TNodeIterator) Node() *ptool.TNode {
	if o == nil {
		return nil
	}
	return o.node
}

// Value -
func (o *TNodeIterator) Value() string {
	if o == nil || o.node == nil {
		return ""
	}
	return o.node.Value
}

// Name -
func (o *TNodeIterator) Name() string {
	if o == nil || o.node == nil {
		return ""
	}
	return o.parser.ByID(o.node.Type)
}

// func getValue(node *ptool.TNode) string {
// 	if node != nil {
// 		return node.Value
// 	}
// 	return ""
// }

func find(where *ptool.TNode, from int, p *ptool.TParser, items ...interface{}) (int, *ptool.TNode) {
	i := from
	for _, item := range items {
		tid := -1
		switch t := item.(type) {
		default:
			panic("illegal type")
		case int:
		case string:
			tid = p.ByName(t)
		}
		node := where
		where = nil
		for i < len(node.Links) {
			v := node.Links[i]
			fmt.Printf(":: %v ~ %v\n", p.ByID(tid), p.ByID(v.Type))
			if v.Type == tid {
				where = node
				break
			}
			i++
		}
		i = 0
		if where == nil {
			return -1, nil
		}
		node = where
	}
	fmt.Println("finish")
	return i, where
}

// func findNode(node *ptool.TNode, types ...string) *ptool.TNode {
// 	if node == nil {
// 		return nil
// 	}
// 	for _, typ := range types {
// 		tid := parser.ByName(typ)
// 		nextLevel := node
// 		node = nil
// 		for _, v := range nextLevel.Links {
// 			if v.Type == tid {
// 				node = v
// 				break
// 			}
// 		}
// 		if node == nil {
// 			return nil
// 		}
// 		nextLevel = node
// 	}
// 	return node
// }
