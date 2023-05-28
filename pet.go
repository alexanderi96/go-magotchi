package main

import (
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
	Dead         bool
}

const (
	Normal    = "-_-"
	Happy     = "^_^"
	Sad       = "T_T"
	Angry     = ">_<"
	Surprised = "0_0"
	Scared    = "!_!"
	Hungry    = "O_O"
	Dead      = "X_X"
)

var stages = [][]string{
	{
		"  __\n (  )\n (__)",
	},
	{
		"(\\_/)\n(#)",
		"(/_\\)\n(#)",
	},
	{
		"(\\(\\\n(#)\n(\")(\")",
		"(/(/ \n(#)\n(\")(\")",
	},
	{
		"(\\(\\\n(#)\no(\")(\")",
		"(/(/ \n(#)\no(\")(\")",
	},
	{
		"(\\(\\\n(#)\no(\")(\")",
		"(/(/ \n(#)\no(\")(\")",
	},
}

func (p *Pet) Move(foods []Food, dirts []Dirt, width int, height int) {
	var foodNeeded = p.Energy < 60 || p.Hunger > 30
	closestFoodDistance := width * height

	var closestFood Food

	oldX := p.X
	oldY := p.Y

	if foodNeeded {
		for _, food := range foods {
			distance := min(abs(p.X-food.X), abs(p.X-food.X-width), abs(p.X-food.X+width)) +
				min(abs(p.Y-food.Y), abs(p.Y-food.Y-height), abs(p.Y-food.Y+height))
			if distance < closestFoodDistance {
				closestFoodDistance = distance
				closestFood = food
			}
		}
	}

	if !foodNeeded || closestFoodDistance >= width*height/2 {

		maxDistance := 0
		bestX, bestY := p.X, p.Y
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				currDistance := 0
				for _, dirt := range dirts {
					currDistance += min(abs(x-dirt.X), abs(x-dirt.X-width), abs(x-dirt.X+width)) +
						min(abs(y-dirt.Y), abs(y-dirt.Y-height), abs(y-dirt.Y+height))
				}
				if currDistance > maxDistance {
					maxDistance = currDistance
					bestX, bestY = x, y
				}
			}
		}
		if p.X != bestX {
			if p.X < bestX {
				p.X++
			} else {
				p.X--
			}
		}
		if p.Y != bestY {
			if p.Y < bestY {
				p.Y++
			} else {
				p.Y--
			}
		}
	} else {

		if p.X < closestFood.X {
			p.X++
		} else if p.X > closestFood.X {
			p.X--
		}
		if p.Y < closestFood.Y {
			p.Y++
		} else if p.Y > closestFood.Y {
			p.Y--
		}
	}

	if p.X < 0 {
		p.X = width - 1
	} else if p.X >= width {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = height - 1
	} else if p.Y >= height {
		p.Y = 0
	}

	if oldX != p.X || oldY != p.Y {
		distance := min(abs(p.X-oldX), abs(p.X-oldX-width), abs(p.X-oldX+width)) +
			min(abs(p.Y-oldY), abs(p.Y-oldY-height), abs(p.Y-oldY+height))
		p.Energy -= distance
		p.Hunger += distance
		p.Distance += distance
	}
}

func min(args ...int) int {
	minValue := args[0]
	for _, v := range args {
		if v < minValue {
			minValue = v
		}
	}
	return minValue
}

func (p *Pet) AgePet() {
	p.Energy--
	p.Hunger++
	p.Age += 1 * time.Second
	if p.Age >= time.Hour || p.Energy <= 0 {
		p.Dead = true
	} else if p.Age >= 45*time.Minute {
		p.Stage = 4
	} else if p.Age >= 30*time.Minute {
		p.Stage = 3
	} else if p.Age >= 15*time.Minute {
		p.Stage = 2
	} else if p.Age >= 3*time.Second {
		p.Stage = 1
	}
}

func (p *Pet) Eat(food int) {
	if p.Hunger > 40 {
		foods = append(foods[:food], foods[food+1:]...)
		p.TotalFood++
		p.DigestionEnd = append(p.DigestionEnd, time.Now().Add(DigestionDuration))
		p.Hunger -= 50
		p.Energy += 30
	}
}

func (p *Pet) Digest() {
	if len(p.DigestionEnd) > 0 {

		if time.Now().After(p.DigestionEnd[0]) {
			p.DigestionEnd = p.DigestionEnd[1:]

			p.TotalDirt++
			p.Energy -= 10
			if p.X > 0 {
				dirts = append(dirts, Dirt{
					X: p.X - 1,
					Y: p.Y,
				})
			} else if p.X < width-1 {
				dirts = append(dirts, Dirt{
					X: p.X + 1,
					Y: p.Y,
				})
			}
		}
	}
}

func (p *Pet) getPetFace() string {
	switch {
	case p.Dead:
		return Dead
	case p.Energy < 50:
		return Sad
	case p.Energy <= 10:
		return Scared
	case p.Hunger > 50:
		return Hungry
	case p.Hunger < 70 && p.Energy > 50:
		return Happy
	default:
		return Normal
	}
}

func (p *Pet) Display() {
	pFrame := strings.Split(stages[p.Stage][p.Frame], "\n")
	petFace := p.getPetFace()
	for y, line := range pFrame {
		xOffset := 0
		for x, char := range line {
			if char == '#' {
				tbprint(p.X+x, p.Y+y, termbox.ColorDefault, termbox.ColorDefault, petFace)
				xOffset = len(petFace) - 1
			} else {
				tbprint(p.X+x+xOffset, p.Y+y, termbox.ColorDefault, termbox.ColorDefault, string(char))
			}
		}
	}
	if !p.Dead {
		p.Frame = (p.Frame + 1) % len(stages[p.Stage])
	}
}
