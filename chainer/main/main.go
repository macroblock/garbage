package main

import (
	"fmt"

	"github.com/macroblock/imed/pkg/misc"
	"github.com/macroblock/imed/pkg/zlog/loglevel"

	"github.com/macroblock/imed/pkg/zlog/zlog"

	"github.com/macroblock/garbage/chainer"
)

var (
	log = zlog.Instance("main")
)

func handler(keychain string) bool {
	fmt.Printf("handler: %q\n", keychain)
	return true
}

func main() {
	log.Add(misc.NewAnsiLogger(loglevel.All, ""))

	ch := chainer.NewChainer()

	ch.Add(
		chainer.NewAction("first", "asdf", "", handler),
		chainer.NewAction("second", "asdg", "", handler),
		chainer.NewAction("third", "asdh", "", handler),
		chainer.NewAction("third", "asdk", "", handler),
		chainer.NewAction("fourth", "asdk", "", handler),
	)

	ch.Apply()

	keys, err := chainer.NewAction("test", "asdgasdfasdkxxxasdhasdasdaasdkx", "", nil).BinaryKey()
	if err != nil {
		log.Error(err)
	}

	for _, key := range keys {
		fmt.Println(key)
		ch.Handle(key)
	}

}
