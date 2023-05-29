package main

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

type Stats struct {
	Pet *Pet
}

func (s *Stats) Display() {
	maxWidth := statsAreaWidth
	stats := []string{
		"Type: " + s.Pet.SpriteName,
		"Energy: " + strconv.Itoa(s.Pet.Energy),
		"Hunger: " + strconv.Itoa(s.Pet.Hunger),
		"Age: " + strconv.Itoa(int(s.Pet.Age.Seconds())),
		"Distance: " + strconv.Itoa(s.Pet.Distance),
		"Total Food: " + strconv.Itoa(s.Pet.TotalFood),
		"Total Dirt: " + strconv.Itoa(s.Pet.TotalDirt),
	}
	for i, stat := range stats {
		if len(stat) > maxWidth {
			stat = "<" + stat[maxWidth-1:]
		}
		tbprint(0, i, termbox.ColorDefault, termbox.ColorDefault, stat)
	}

	for i := 0; i < height; i++ {
		tbprint(statsAreaWidth, i, termbox.ColorDefault, termbox.ColorDefault, "║")
	}
}
