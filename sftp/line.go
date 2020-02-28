package main

import "github.com/c-bata/go-prompt"

type Line struct {
	Document prompt.Document
	Args     []string
}

func (l Line) Cmd() string {
	return l.Args[0]
}

func (l Line) OptionBeforeCursor() string {
	return ""
}

// #TODO: 28-02-20 implement me //

func (l Line) CursorIndex() int {
	return 0
}
