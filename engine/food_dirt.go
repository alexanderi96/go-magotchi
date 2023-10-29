package engine

import (
	"math/rand"
	"time"
)

const (
	MaxFoodEnergy     = 20
	MinFoodEnergy     = 10
	foodSpawnInterval = 10 * time.Second
)

type Food struct {
	X, Y, Energy float32
	Eaten        bool
	SpawnTime    time.Time
}

type Dirt struct {
	X, Y   float32
	Energy float32
}

func (f *Food) Size() float32 {
	return f.Energy
}
func (d *Dirt) Size() float32 {
	return d.Energy
}

func NewFood(x, y float32) *Food {
	return &Food{
		X:         x,
		Y:         y,
		Eaten:     false,
		SpawnTime: time.Now(),
		Energy:    float32(rand.Intn(MaxFoodEnergy-MinFoodEnergy+1) + MinFoodEnergy),
	}
}
