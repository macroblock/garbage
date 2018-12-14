package main

import (
	"fmt"
	"path"
	"runtime"
	"strings"

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
	errInvalidObject = "iterator is invalid"
)

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
		log.Errorf(true, "%v%v", "Clone: ", errInvalidObject)
		return nil
	}
	ret := &TNodeIterator{}
	*ret = *o
	return ret
}

// Enter -
func (o *TNodeIterator) Enter() *TNodeIterator {
	errPrefix := "Enter: "
	if o == nil || o.root == nil {
		log.Errorf(true, "%v%v", errPrefix, errInvalidObject)
		return nil
	}
	if o.retNode == nil {
		log.Errorf(true, "%vinvalid operation", errPrefix)
		return nil
	}
	o.root = o.retNode
	o.retNode = nil
	o.index = 0
	if o.index > len(o.root.Links)-1 {
		o.root = nil
		log.Errorf(true, "%vnode has no links", errPrefix)
		return nil
	}
	o.inNode = o.root.Links[o.index]
	return o
}

// AutoIgnore -
func (o *TNodeIterator) AutoIgnore(items ...interface{}) *TNodeIterator {
	errPrefix := "AutoIgnore: "
	if o == nil || o.root == nil {
		log.Errorf(true, "%v%v", errPrefix, errInvalidObject)
		return nil
	}
	ok := false
	o.ignore, ok = itemsToIntSlice(errPrefix, o.parser, items...)
	if !ok {
		o.root = nil
		return nil
	}
	return o
}
func retrieveCallInfo() {
	pc, file, line, _ := runtime.Caller(1)
	_, fileName := path.Split(file)
	name := runtime.FuncForPC(pc).Name()
	parts := strings.Split(name, ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	fmt.Printf("name: %v\nfile: %v\nline: %v\nfunc: %v\npkg: %v\n\n", name, fileName, line, funcName, packageName)
	// return &callInfo{
	// 	packageName: packageName,
	// 	fileName:    fileName,
	// 	funcName:    funcName,
	// 	line:        line,
	// }
}

// Accept -
func (o *TNodeIterator) Accept(items ...interface{}) *TNodeIterator {
	// retrieveCallInfo()
	errPrefix := "Accept: "
	if o == nil || o.root == nil || o.inNode == nil {
		log.Errorf(true, "%v%v", errPrefix, errInvalidObject)
		o.root = nil
		return nil
	}

	if len(items) > 0 {
		where, ok := itemsToIntSlice(errPrefix, o.parser, items...)
		if !ok {
			o.root = nil
			return nil
		}
		// fmt.Printf("id: %v; where: %v\n", o.inNode.Type, where)
		if !in(o.inNode.Type, where) {
			o.root = nil
			log.Errorf(true, "%vunexpected node %q", errPrefix, o.parser.ByID(o.inNode.Type))
			// return nil
		}
	}

	o.index++
	if o.index > len(o.root.Links)-1 {
		if o.inNode == nil {
			log.Errorf(true, "%vindex %v out of range 0..%v", errPrefix, o.index, len(o.root.Links)-1)
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
		log.Errorf(true, "%v%v", "Node: ", errInvalidObject)
		return nil
	}
	return o.retNode
}

// Value -
func (o *TNodeIterator) Value() string {
	if o == nil || o.root == nil {
		log.Errorf(true, "%v%v", "Value: ", errInvalidObject)
		return ""
	}
	return o.retNode.Value
}

// ID -
func (o *TNodeIterator) ID() int {
	if o == nil || o.root == nil {
		log.Errorf(true, "%v%v", "ID: ", errInvalidObject)
		return -1
	}
	return o.retNode.Type
}

// Name -
func (o *TNodeIterator) Name() string {
	if o == nil || o.root == nil {
		log.Errorf(true, "%v%v", "Name: ", errInvalidObject)
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

func itemToInt(errPrefix string, p *ptool.TParser, item interface{}) int {
	ret := -1
	switch t := item.(type) {
	default:
		log.Errorf(true, "%villegal type %v", errPrefix, t)
	case int:
		ret = t
		if ret < 0 {
			log.Errorf(true, "%vinvalid node id %v", errPrefix, t)
		}
	case string:
		ret = p.ByName(t)
		if ret < 0 {
			log.Errorf(true, "%vunknown node type %q", errPrefix, t)
		}
	}
	return ret
}

func itemsToIntSlice(errPrefix string, p *ptool.TParser, items ...interface{}) ([]int, bool) {
	ret := []int{}
	ok := true
	for _, i := range items {
		id := itemToInt(errPrefix, p, i)
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
