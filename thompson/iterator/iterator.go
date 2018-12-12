package iterator

import (
	"github.com/macroblock/imed/pkg/ptool"
)

// Type -
type Type struct {
	root   *ptool.TNode
	node   *ptool.TNode
	idx    int
	parser *ptool.TParser
}

// New -
func New(root *ptool.TNode, parser *ptool.TParser) *Type {
	return &Type{root: root, idx: -1, parser: parser}
}

// First -
func (o *Type) First() *Type {
	if o == nil {
		return nil
	}
	o.idx = 0
	if o.root == nil || o.idx > len(o.root.Links)-1 {
		o.idx, o.node = -1, nil
		return o
	}
	o.node = o.root.Links[o.idx]
	return o
}

// Next -
func (o *Type) Next() *Type {
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

// FindFirst -
func (o *Type) FindFirst(item interface{}) *Type {
	if o == nil {
		return nil
	}
	o.idx, o.node = find(o.root, 0, item, o.parser)
	return o
}

// FindNext -
func (o *Type) FindNext(item interface{}) *Type {
	if o == nil || o.idx < 0 ||
		o.root == nil || o.idx > len(o.root.Links)-2 {
		return o
	}
	o.idx, o.node = find(o.root, o.idx+1, item, o.parser)
	return o
}

// ForEach -
func (o *Type) ForEach(fn func(*Type)) *Type {
	if o == nil || o.node == nil {
		return o
	}
	for it := New(o.node, o.parser).First(); it.Node() != nil; it.Next() {
		fn(it)
	}
	return o
}

// Node -
func (o *Type) Node() *ptool.TNode {
	if o == nil {
		return nil
	}
	return o.node
}

// Value -
func (o *Type) Value() string {
	if o == nil || o.idx < 0 ||
		o.root == nil || o.idx > len(o.root.Links)-1 {
		return ""
	}
	return o.root.Links[o.idx].Value
}

func getValue(node *ptool.TNode) string {
	if node != nil {
		return node.Value
	}
	return ""
}

func find(where *ptool.TNode, from int, item interface{}, p *ptool.TParser) (int, *ptool.TNode) {
	tid := -1
	switch t := item.(type) {
	default:
		return -1, nil
	case int:
	case string:
		tid = p.ByName(t)
	}
	for i, v := range where.Links[from:] {
		if v.Type == tid {
			return from + i, v
		}
	}
	return -1, nil
}
