package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Food struct {
	Texture   rl.Texture2D
	X, Y      float32
	Eaten     bool
	SpawnTime time.Time
	Energy    int
}

type Dirt struct {
	Texture rl.Texture2D
	X, Y    float32
	Eaten   bool
}

func (f *Food) DrawFoods() {
	if !f.Eaten {

		textureWidth := float32(foodTexture.Width)
		textureHeight := float32(foodTexture.Height)
		scale := float32(f.Energy) / foodSize
		x := f.X - (textureWidth*scale)/2
		y := f.Y - (textureHeight*scale)/2

		rl.DrawTextureEx(foodTexture, rl.Vector2{X: x, Y: y}, 0, scale, rl.White)
	}
}
