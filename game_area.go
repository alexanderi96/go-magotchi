package main

import (
	"github.com/nsf/termbox-go"
)

type GameArea struct {
	Width  int
	Height int
	StartX int
	StartY int
}

func (ga *GameArea) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	ga.Width, ga.Height = termbox.Size()
	ga.StartX = 0
	ga.StartY = 0
	return nil
}

func (ga *GameArea) Close() {
	termbox.Close()
}
