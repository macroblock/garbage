package main

import (
	"fmt"

	"github.com/macroblock/garbage/zlog"
)

var (
	log = zlog.Instance()
)

var (
	text  = "longer one-lined message to print"
	text2 = "unexpected multi-lined message to print for\ndemonstration purposes and nothing else"

	err  = fmt.Errorf("unexpected error is for demonstration porposes only")
	err2 = fmt.Errorf("something wrong heppend with half-fatal consequence\nyou have to do something")
)

func main() {
	log.Log(zlog.LevelInfo, nil, "test message")
	log.Log(zlog.LevelError, fmt.Errorf("asdf"), "because of")
	log.Log(zlog.LevelError, fmt.Errorf("line\nenother line"), "because of adfaaf\nasdfasfaf")

	for level := zlog.LevelCritical; level <= zlog.LevelDebug; level++ {
		log.Log(level, err, text)
	}
	for level := zlog.LevelCritical; level <= zlog.LevelDebug; level++ {
		log.Log(level, err2, text2)
	}

	log.Error(1, "mesage")
}
