package conio

import "io"

type tActionMode struct {
	stack []tActionModeItem
}

type tActionModeItem struct {
	mode    string
	closers []io.Closer
}

// Push -
func (am *tActionMode) Push(mode string) {
	item := tActionModeItem{mode: mode, closers: nil}
	am.stack = append(am.stack, item)
}

// Pop -
func (am *tActionMode) Pop(mode string) {
	am.stack = append(am.stack, mode)
}
