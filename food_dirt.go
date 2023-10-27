package main

import (
	"log"
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
	if f.X >= float32(world.WorldWidth) || f.X < float32(0) || f.Y > float32(world.WorldHeight) || f.Y < float32(0) {
		log.Printf("food out of bounds: %f %f ", f.X, f.Y)
		log.Printf("world bounds: %f %f ", float32(world.WorldWidth), float32(world.WorldHeight))
	} else if !f.Eaten {
		textureWidth := float32(foodTexture.Width)
		textureHeight := float32(foodTexture.Height)
		scale := float32(f.Energy) / foodSize
		x := f.X - (textureWidth*scale)/2
		y := f.Y - (textureHeight*scale)/2

		rl.DrawTextureEx(foodTexture, rl.Vector2{X: x, Y: y}, 0, scale, rl.White)
	}
}

func (f *Dirt) Draw() {
	// TODO
}

func NewFood(x, y float32) *Food {
	return &Food{
		X:         x,
		Y:         y,
		Eaten:     false,
		Texture:   foodTexture,
		SpawnTime: time.Now(),
		Energy:    float32(rand.Intn(MaxFoodEnergy-MinFoodEnergy+1) + MinFoodEnergy),
	}
}
