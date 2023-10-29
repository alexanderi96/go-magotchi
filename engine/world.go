package engine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alexanderi96/go-magotchi/config"
)

type World struct {
	Foods  []*Food
	Dirts  []*Dirt
	Pet    *Pet
	Config *config.Config
	Paused bool

	foodTicker      *time.Ticker
	AnimationTicker *time.Ticker
	timeTicker      *time.Ticker
	movementTicker  *time.Ticker
}

func NewGame() (*World, error) {

	config, err := config.ReadConfig("./config.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	world := &World{
		Foods: []*Food{},
		Dirts: []*Dirt{},
		Pet: &Pet{
			X:         float32(config.ViewportX / 2),
			Y:         float32(config.ViewportY / 2),
			Health:    100,
			Hunger:    50,
			Happiness: 50,
			Energy:    100,
			Age:       0,
			Sleeping:  false,
			Dead:      false,
		},
		Config: config,
		Paused: false,

		foodTicker:      time.NewTicker(10 * time.Second),
		AnimationTicker: time.NewTicker(150 * time.Millisecond),
		timeTicker:      time.NewTicker(time.Second),
		movementTicker:  time.NewTicker(time.Millisecond),
	}

	return world, nil
}

func (w *World) spawnFood() {
	X := float32(rand.Intn(int(w.Config.ViewportX)))
	Y := float32(rand.Intn(int(w.Config.ViewportY)))

	for _, food := range w.Foods {
		if food.X == X && food.Y == Y {
			return
		}
	}

	w.Foods = append(w.Foods, NewFood(X, Y))
}

func (w *World) spawnDirt() {
	// TODO: dirt spawn
}

func (w *World) Update() {
	select {
	case <-w.foodTicker.C:
		if !w.Pet.Dead {
			w.spawnFood()
		}

	case <-w.timeTicker.C:
		if !w.Pet.Dead {
			w.Pet.GetOlder()
		}

	case <-w.movementTicker.C:
		if w.Pet.WantToMove() {
			w.Pet.MoveToFood(w.Config.ViewportX, w.Config.ViewportY, w.Foods)
		}
	default:
	}
	w.Pet.Update()
}
