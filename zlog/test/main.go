package main

import (
	"fmt"

	"github.com/macroblock/garbage/zlog"
)

var (
	log = zlog.Instance()
)

var (
	text  = "a bit longer single-line message to print"
	text2 = "unexpected multi-line message to print for\ndemonstration purpose and nothing else"

	err  = fmt.Errorf("another unexpected error for demonstration purpose only")
	err2 = fmt.Errorf("something wrong heppend causing half-fatal consequences\nso you have to do something smart and brave")
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
	for level := zlog.LevelCritical; level <= zlog.LevelDebug; level++ {
		log.Log(level, true, text)
	}

	for level := zlog.LevelCritical; level <= zlog.LevelDebug; level++ {
		log.Log(level, true)
	}

	for level := zlog.LevelCritical; level <= zlog.LevelDebug; level++ {
		log.Log(level, true, "")
	}

	log.Error(1, "mesage")
}
