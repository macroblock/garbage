package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/macroblock/imed/pkg/misc"
	"github.com/macroblock/imed/pkg/zlog/loglevel"
	runewidth "github.com/mattn/go-runewidth"

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

func oldMain() {
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

var row = 0
var style = tcell.StyleDefault

func putln(s tcell.Screen, str string) {

	puts(s, style, 1, row, str)
	row++
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}

// KeyStruct -
type KeyStruct struct {
	Key   int16
	mod   int16
	r     rune
	alias string
}

const (
	modNone  int16 = 0
	modShift int16 = 1 << iota
	modCtrl
	modAlt
	modMeta
)

func (o KeyStruct) String() string {
	slice := []string{}
	if o.mod&modMeta != 0 {
		slice = append(slice, "Meta")
	}
	if o.mod&modCtrl != 0 {
		slice = append(slice, "Ctrl")
	}
	if o.mod&modAlt != 0 {
		slice = append(slice, "Alt")
	}
	if o.mod&modShift != 0 {
		slice = append(slice, "Shift")
	}
	if len(o.alias) == 0 {
		slice = append(slice, fmt.Sprintf("%q", o.r))
	} else {
		slice = append(slice, fmt.Sprintf("%v", o.alias))
	}
	return strings.Join(slice, "+")
}

func initKeyStruct(keyEvent tcell.EventKey) KeyStruct {
	key := KeyStruct{}

	mod := keyEvent.Modifiers()
	if mod&tcell.ModMeta != 0 {
		key.mod |= modMeta
	}
	if mod&tcell.ModCtrl != 0 {
		key.mod |= modCtrl
	}
	if mod&tcell.ModAlt != 0 {
		key.mod |= modAlt
	}
	if mod&tcell.ModShift != 0 {
		key.mod |= modShift
	}
	k := keyEvent.Key()
	if s, exists := tcell.KeyNames[k]; exists {
		if key.mod&modCtrl != 0 && strings.HasPrefix(s, "Ctrl-") {
			s = s[5:]
		}
		if key.mod&modShift == 0 && len(s) == 1 {
			s = strings.ToLower(s)
		}
		key.alias = s
		return key
	}
	if k == tcell.KeyRune {
		key.alias = string(keyEvent.Rune())
		if key.alias == " " {
			key.alias = "Space"
		}
		return key
	}
	panic("something went wrong")

}

func parseTCellKey(str string) KeyStruct {
	invMap := map[string]int16{}
	for key, val := range tcell.KeyNames {
		invMap[val] = int16(key)
	}
	key := KeyStruct{}
	for len(str) > 0 {
		switch {
		default:
			if _, exists := invMap[str]; exists {
				key.alias = str
				return key
			}
			key.r, _ = utf8.DecodeRuneInString(str)
			return key
		case strings.HasPrefix(str, "Rune["):
			str = strings.TrimPrefix(str, "Rune[")
			key.r, _ = utf8.DecodeRuneInString(str)
			return key
		case strings.HasPrefix(str, "Shift+"):
			str = strings.TrimPrefix(str, "Shift+")
			key.mod |= modShift
		case strings.HasPrefix(str, "Ctrl+"):
			str = strings.TrimPrefix(str, "Ctrl+")
			key.mod |= modCtrl
		case strings.HasPrefix(str, "Alt+"):
			str = strings.TrimPrefix(str, "Alt+")
			key.mod |= modAlt
		case strings.HasPrefix(str, "Meta+"):
			str = strings.TrimPrefix(str, "Meta+")
			key.mod |= modMeta
		}
	}
	panic("something went wrong")
}

func main() {

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	plain := tcell.StyleDefault
	bold := style.Bold(true)

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	quit := make(chan struct{})

	style = bold
	putln(s, "Press ESC to Exit")
	putln(s, "Character set: "+s.CharacterSet())
	style = plain

	putln(s, "English:   October")
	putln(s, "Icelandic: október")
	putln(s, "Arabic:    أكتوبر")
	putln(s, "Russian:   октября")
	putln(s, "Greek:     Οκτωβρίου")
	putln(s, "Chinese:   十月 (note, two double wide characters)")
	putln(s, "Combining: A\u030a (should look like Angstrom)")
	putln(s, "Emoticon:  \U0001f618 (blowing a kiss)")
	putln(s, "Airplane:  \u2708 (fly away)")
	putln(s, "Command:   \u2318 (mac clover key)")
	putln(s, "Enclose:   !\u20e3 (should be enclosed exclamation)")
	putln(s, "ZWJ:       \U0001f9db\u200d\u2640 (female vampire)")
	putln(s, "ZWJ:       \U0001f9db\u200d\u2642 (male vampire)")
	putln(s, "Family:    \U0001f469\u200d\U0001f467\u200d\U0001f467 (woman girl girl)\n")

	putln(s, "")
	putln(s, "Box:")
	putln(s, string([]rune{
		tcell.RuneULCorner,
		tcell.RuneHLine,
		tcell.RuneTTee,
		tcell.RuneHLine,
		tcell.RuneURCorner,
	}))
	putln(s, string([]rune{
		tcell.RuneVLine,
		tcell.RuneBullet,
		tcell.RuneVLine,
		tcell.RuneLantern,
		tcell.RuneVLine,
	})+"  (bullet, lantern/section)")
	putln(s, string([]rune{
		tcell.RuneLTee,
		tcell.RuneHLine,
		tcell.RunePlus,
		tcell.RuneHLine,
		tcell.RuneRTee,
	}))
	putln(s, string([]rune{
		tcell.RuneVLine,
		tcell.RuneDiamond,
		tcell.RuneVLine,
		tcell.RuneUArrow,
		tcell.RuneVLine,
	})+"  (diamond, up arrow)")
	putln(s, string([]rune{
		tcell.RuneLLCorner,
		tcell.RuneHLine,
		tcell.RuneBTee,
		tcell.RuneHLine,
		tcell.RuneLRCorner,
	}))

	s.Show()
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
				keyName := ev.Name()
				// putln(s, fmt.Sprintf("%v --> %v", keyName, parseTCellKey(keyName)))
				putln(s, fmt.Sprintf("%v --> %v", keyName, initKeyStruct(*ev)))
				s.Sync()
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	<-quit

	s.Fini()
}
