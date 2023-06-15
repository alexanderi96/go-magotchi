package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MaxFoodEnergy     = 20
	MinFoodEnergy     = 10
	foodSize          = 10
	foodSpawnInterval = 10 * time.Second
)

var (
	foodTicker  *time.Ticker
	foodTexture rl.Texture2D
)

type Food struct {
	Texture      rl.Texture2D
	X, Y, Energy float32
	Eaten        bool
	SpawnTime    time.Time
}

type Dirt struct {
	Texture rl.Texture2D
	X, Y    float32
}

func (f *Food) Draw() {
	if !f.Eaten {
		textureWidth := float32(foodTexture.Width)
		textureHeight := float32(foodTexture.Height)

		destRec := rl.NewRectangle(float32(f.X), float32(f.Y), float32(world.cellSize), float32(world.cellSize))

		rl.DrawTexturePro(foodTexture, rl.NewRectangle(0, 0, textureWidth, textureHeight), destRec, rl.NewVector2(0, 0), 0, rl.White)
	}
}

func (f *Dirt) Draw() {
	// TODO
}

func NewFood() *Food {
	return &Food{
		X:         float32((rand.Intn(world.WorldWidth/world.cellSize-2) + 1) * world.cellSize),
		Y:         float32(world.cellSize) + float32((rand.Intn(world.WorldHeight/world.cellSize-2)+1)*world.cellSize),
		Eaten:     false,
		Texture:   foodTexture,
		SpawnTime: time.Now(),
		Energy:    float32(rand.Intn(MaxFoodEnergy-MinFoodEnergy+1) + MinFoodEnergy),
	}
}
