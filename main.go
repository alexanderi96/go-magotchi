package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const DigestionDuration = 60 * time.Second

type Pet struct {
	Stage        int
	Energy       int
	Hunger       int
	Age          time.Duration
	X            int
	Y            int
	Frame        int
	DigestionEnd []time.Time
	Distance     int
	TotalFood    int
	TotalDirt    int
}

type Food struct {
	X int
	Y int
}

type Dirt struct {
	X int
	Y int
}

var stages = [][]string{
	{
		"  __\n (  )\n (__)",
		"  __\n (  )\n (__)",
		"  __\n (  )\n (__)",
	},
	{
		"(\\_/)\n(-.-)",
		"(/_\\)\n(-.-)",
		"(\\_/)\n(-.-)",
	},
	{
		"(\\(\\\n(-.-)\n(\")(\")",
		"(/(/ \n(-.-)\n(\")(\")",
		"(\\(\\\n(-.-)\n(\")(\")",
	},
	{
		"(\\(\\\n(-.-)\no(\")(\")",
		"(/(/ \n(-.-)\no(\")(\")",
		"(\\(\\\n(-.-)\no(\")(\")",
	},
	{
		"(\\(\\\n(-_-)\no(\")(\")",
		"(/(/ \n(-_-)\no(\")(\")",
		"(\\(\\\n(-_-)\no(\")(\")",
	},
	{
		"(x_x)",
		"(X_X)",
		"(x_x)",
	},
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	width, height := termbox.Size()

	// Define game area size
	gameAreaWidth := width - 20 // 20 is the width of the stats area, adjust if needed
	if gameAreaWidth < 0 {
		gameAreaWidth = 0
	}
	gameAreaHeight := height
	gameAreaStartX := 0
	gameAreaStartY := 0

	// Game objects initialization
	pet := Pet{
		X:            gameAreaStartX + gameAreaWidth/2,
		Y:            gameAreaStartY + gameAreaHeight/2,
		Stage:        0,
		Frame:        0,
		Age:          0,
		Energy:       100,
		Hunger:       0,
		DigestionEnd: []time.Time{},
		TotalFood:    0,
		TotalDirt:    0,
		Distance:     0,
	}

	foods := []Food{}
	dirts := []Dirt{}

	// Tickers
	gameTicker := time.NewTicker(1 * time.Second)
	foodTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-gameTicker.C:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

			// Pet movement
			oldX := pet.X
			oldY := pet.Y
			// Move the pet if it is not in stage 0 or 5
			if pet.Stage != 0 && pet.Stage != 5 {
				pet = movePet(pet, foods, dirts, width, height)
				pet.Energy--
			}
			pet.Distance += abs(pet.X-oldX) + abs(pet.Y-oldY)

			// Pet digestion
			now := time.Now()
			i := 0
			for _, end := range pet.DigestionEnd {
				if now.Before(end) {
					pet.DigestionEnd[i] = end
					i++
				} else {
					pet.TotalDirt++
					dirts = append(dirts, Dirt{
						X: pet.X + 2,
						Y: pet.Y + 2,
					})
				}
			}
			pet.DigestionEnd = pet.DigestionEnd[:i]

			// Pet eating
			i = 0
			for _, food := range foods {
				if pet.X <= food.X && food.X <= pet.X+3 &&
					pet.Y <= food.Y && food.Y <= pet.Y+3 {
					pet.DigestionEnd = append(pet.DigestionEnd, now.Add(DigestionDuration))
					pet.TotalFood++
				} else {
					foods[i] = food
					i++
				}
			}
			foods = foods[:i]

			// Pet aging
			pet.Age += 1 * time.Second
			if pet.Age >= time.Hour {
				pet.Stage = 5
			} else if pet.Age >= 45*time.Minute {
				pet.Stage = 4
			} else if pet.Age >= 30*time.Minute {
				pet.Stage = 3
			} else if pet.Age >= 15*time.Minute {
				pet.Stage = 2
			} else if pet.Age >= 3*time.Second {
				pet.Stage = 1
			}

			// Pet animation
			pet.Frame = (pet.Frame + 1) % 3

			// Print stats
			stats := []string{
				"Energy: " + strconv.Itoa(pet.Energy),
				"Hunger: " + strconv.Itoa(pet.Hunger),
				"Age: " + strconv.Itoa(int(pet.Age.Seconds())),
				"Distance: " + strconv.Itoa(pet.Distance),
				"Total Food: " + strconv.Itoa(pet.TotalFood),
				"Total Dirt: " + strconv.Itoa(pet.TotalDirt),
			}
			for i, stat := range stats {
				tbprint(0, i, termbox.ColorDefault, termbox.ColorDefault, stat) // Print stats on the left
			}
			// Draw separator line between stats and game
			for i := 0; i < height; i++ {
				termbox.SetCell(21, i, '|', termbox.ColorWhite, termbox.ColorBlack)
			}

			// Print the pet
			petFrame := strings.Split(stages[pet.Stage][pet.Frame], "\n")
			for y, line := range petFrame {
				for x, char := range line {
					tbprint(pet.X+x, pet.Y+y, termbox.ColorDefault, termbox.ColorDefault, string(char))
				}
			}

			// Print the food
			for _, food := range foods {
				if food.X >= gameAreaStartX && food.X < gameAreaStartX+gameAreaWidth &&
					food.Y >= gameAreaStartY && food.Y < gameAreaStartY+gameAreaHeight {
					tbprint(food.X, food.Y, termbox.ColorGreen, termbox.ColorDefault, "*")
				}
			}

			// Print the dirt
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

func movePet(pet Pet, foods []Food, dirts []Dirt, width int, height int) Pet {
	closestFoodDistance := width * height
	closestDirtDistance := width * height
	var closestFood Food
	var closestDirt Dirt

	for _, food := range foods {
		distance := abs(pet.X-food.X) + abs(pet.Y-food.Y)
		if distance < closestFoodDistance {
			closestFoodDistance = distance
			closestFood = food
		}
	}

	for _, dirt := range dirts {
		distance := abs(pet.X-dirt.X) + abs(pet.Y-dirt.Y)
		if distance < closestDirtDistance {
			closestDirtDistance = distance
			closestDirt = dirt
		}
	}

	if closestFoodDistance < closestDirtDistance {
		// Move towards the closest food
		if pet.X < closestFood.X {
			pet.X++
		} else if pet.X > closestFood.X {
			pet.X--
		}
		if pet.Y < closestFood.Y {
			pet.Y++
		} else if pet.Y > closestFood.Y {
			pet.Y--
		}
	} else {
		// Move away from the closest dirt
		if pet.X < closestDirt.X {
			pet.X--
		} else if pet.X > closestDirt.X {
			pet.X++
		}
		if pet.Y < closestDirt.Y {
			pet.Y--
		} else if pet.Y > closestDirt.Y {
			pet.Y++
		}
	}

	// If the pet can't move away from the closest dirt, move it in a random direction
	if pet.X == closestDirt.X && pet.Y == closestDirt.Y {
		switch rand.Intn(4) {
		case 0:
			pet.X++
		case 1:
			pet.X--
		case 2:
			pet.Y++
		case 3:
			pet.Y--
		}
	}

	// Wrap the pet around the screen if it goes out of bounds
	if pet.X < 0 {
		pet.X = width - 1
	} else if pet.X >= width {
		pet.X = 0
	}
	if pet.Y < 0 {
		pet.Y = height - 1
	} else if pet.Y >= height {
		pet.Y = 0
	}

	return pet
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
