package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

var (
	width, height, statsAreaWidth, gameAreaWidth, gameAreaHeight, gameAreaStartX, gameAreaStartY int
	foods                                                                                        = []Food{}
	dirts                                                                                        = []Dirt{}
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	width, height = termbox.Size()

	statsAreaWidth = (width / 5)
	gameAreaWidth = width - statsAreaWidth
	if gameAreaWidth < 0 {
		gameAreaWidth = 0
	}
	gameAreaHeight = height
	gameAreaStartX = statsAreaWidth + 1
	gameAreaStartY = 0

	pet := &Pet{
		X:            gameAreaStartX + gameAreaWidth/2,
		Y:            gameAreaStartY + gameAreaHeight/2,
		Stage:        0,
		Frame:        0,
		Age:          0,
		Energy:       100,
		Hunger:       50,
		DigestionEnd: []time.Time{},
		TotalFood:    0,
		TotalDirt:    0,
		Distance:     0,
		Dead:         false,
	}

	stats := &Stats{Pet: pet}

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	gameTicker := time.NewTicker(1 * time.Second)
	foodTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:

				switch ev.Key {
				case termbox.KeyEsc:

					handlePause()
				}
			case termbox.EventResize:

			}
		case <-gameTicker.C:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

			if !pet.Dead {
				pet.AgePet()
			}

			if pet.Stage != 0 && !pet.Dead {
				pet.Move(foods, dirts, width, height)
			}
			pet.Digest()

			stats.Display()

			pet.Display()

			for i, food := range foods {
				if food.X >= gameAreaStartX && food.X < gameAreaStartX+gameAreaWidth &&
					food.Y >= gameAreaStartY && food.Y < gameAreaStartY+gameAreaHeight {
					tbprint(food.X, food.Y, termbox.ColorGreen, termbox.ColorDefault, "*")
				}

				if food.X == pet.X && food.Y == pet.Y {
					pet.Eat(i)
				}
			}

			for _, dirt := range dirts {
				tbprint(dirt.X, dirt.Y, termbox.ColorRed, termbox.ColorDefault, "~")
			}

			termbox.Flush()

		case <-foodTicker.C:
			foodX := rand.Intn(gameAreaWidth) + gameAreaStartX
			foodY := rand.Intn(gameAreaHeight) + gameAreaStartY
			foods = append(foods, Food{
				X: foodX,
				Y: foodY,
			})
		}
	}
}

func tbprint(x int, y int, fg termbox.Attribute, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func drawMenu(title string, options []string) {
	menuWidth := len(title)
	for _, option := range options {
		if len(option) > menuWidth {
			menuWidth = len(option)
		}
	}
	menuWidth += 6

	menuHeight := len(options) + 5

	menuStartX := (width - menuWidth) / 2
	menuStartY := (height - menuHeight) / 2

	for i := 0; i < menuHeight; i++ {
		for j := 0; j < menuWidth; j++ {
			var ch rune
			switch {
			case i == 0 && j == 0:
				ch = '╔'
			case i == menuHeight-1 && j == 0:
				ch = '╚'
			case i == 0 && j == menuWidth-1:
				ch = '╗'
			case i == menuHeight-1 && j == menuWidth-1:
				ch = '╝'
			case i == 0 || i == menuHeight-1:
				ch = '═'
			case j == 0 || j == menuWidth-1:
				ch = '║'
			default:
				ch = ' '
			}
			termbox.SetCell(menuStartX+j, menuStartY+i, ch, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	tbprint(menuStartX+1, menuStartY, termbox.ColorWhite, termbox.ColorDefault, title)

	for i, option := range options {
		tbprint(menuStartX+(menuWidth-len(option))/2, menuStartY+i+3, termbox.ColorWhite, termbox.ColorDefault, option)
	}
}

func handlePause() {
	menuOptions := []string{
		"1. Resume",
		"2. Save",
		"3. Exit",
	}

	drawMenu("Pause", menuOptions)

	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case '1':
				return
			case '2':

			case '3':
				termbox.Close()
				os.Exit(0)
			}
		}
	}
}
