package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
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

	name, sprite := GetRandomPet()
	pet := &Pet{
		Sprite:       sprite,
		SpriteName:   name,
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
	pet.AdjustMoveSpeed()

	stats := &Stats{Pet: pet}

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	gameTicker := time.NewTicker(33 * time.Millisecond)
	foodTicker := time.NewTicker(10 * time.Second)
	secondTicker := time.NewTicker(time.Second)

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

		case <-pet.MoveTicker.C:
			if pet.Stage != 0 && !pet.Dead {
				pet.Move(foods, dirts, width, height)
			}

		case <-secondTicker.C:
			if !pet.Dead {
				pet.AgePet()
				pet.Digest()
				pet.UpdateFrame()
			} else {
				handleEnd(pet)
			}
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

func drawMenu(title string, topText string, options []string) {
	maxWidth := (2 * width) / 3
	wrappedTopText, lines := wrapText(topText, maxWidth-6)

	menuWidth := max(max(len(title), longestLine(wrappedTopText)), longestOption(options)) + 6
	menuHeight := len(options) + lines + 5

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

	for i, line := range strings.Split(wrappedTopText, "\n") {
		tbprint(menuStartX+1, menuStartY+i+2, termbox.ColorWhite, termbox.ColorDefault, line)
	}

	for i, option := range options {
		tbprint(menuStartX+1, menuStartY+i+lines+3, termbox.ColorWhite, termbox.ColorDefault, option)
	}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func longestOption(options []string) int {
	longest := 0
	for _, option := range options {
		if len(option) > longest {
			longest = len(option)
		}
	}
	return longest
}

func longestLine(text string) int {
	longest := 0
	for _, line := range strings.Split(text, "\n") {
		if len(line) > longest {
			longest = len(line)
		}
	}
	return longest
}

func wrapText(text string, maxWidth int) (string, int) {
	words := strings.Fields(text)
	if len(words) == 0 {
		return "", 0
	}

	wrapped := words[0]
	currentLength := len(wrapped)
	lines := 1

	for _, word := range words[1:] {
		if currentLength+len(word)+1 > maxWidth {
			wrapped += "\n" + word
			currentLength = len(word)
			lines++
		} else {
			wrapped += " " + word
			currentLength += len(word) + 1
		}
	}
	return wrapped, lines
}

func handlePause() {
	menuOptions := []string{
		"1. Resume",
		"2. Restart",
		"3. Save",
		"4. Load",
		"5. Exit",
	}

	drawMenu("Pause", "", menuOptions)

	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case '1':
				return
			case '2':

			case '3':

			case '4':

			case '5':
				termbox.Close()
				os.Exit(0)
			}
		}
	}
}

func handleEnd(p *Pet) {
	menuOptions := []string{
		"1. Restart",
		"2. Load",
		"3. Exit",
	}

	gameSummary := fmt.Sprintf(`Our pet, %s, lived for %s. With an energy peak at %d and a hunger level reaching %d, %s was a real adventurer, covering %d units of distance. It feasted on %d food items and, naturally, produced %d piles of dirt. Sadly, %s is now %s. We'll always cherish the memories we made together.`,
		p.SpriteName,
		p.Age,
		p.Energy,
		p.Hunger,
		p.SpriteName,
		p.Distance,
		p.TotalFood,
		p.TotalDirt,
		p.SpriteName,
		func() string {
			if p.Dead {
				return "no longer with us"
			} else {
				return "still with us"
			}
		}(),
	)

	drawMenu("Your pet is Dead", gameSummary, menuOptions)

	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case '1':

			case '2':

			case '3':
				termbox.Close()
				os.Exit(0)
			}
		}
	}
}
