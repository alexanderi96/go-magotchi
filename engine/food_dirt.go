package engine

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

type Food struct {
	Texture      *Texture
	X, Y, Energy float32
	Eaten        bool
	SpawnTime    time.Time
}

type Dirt struct {
	Texture *Texture
	X, Y    float32
}

func (f *Food) Draw(x, y int32) {
	if f.X >= float32(x) || f.X < float32(0) || f.Y > float32(y) || f.Y < float32(0) {
		log.Printf("food out of bounds: %f %f ", f.X, f.Y)
		log.Printf("viewport bounds: %f %f ", float32(x), float32(y))
	} else if !f.Eaten {
		textureWidth := float32(f.Texture.Data.Width)
		textureHeight := float32(f.Texture.Data.Height)
		scale := float32(f.Energy) / foodSize
		x := f.X - (textureWidth*scale)/2
		y := f.Y - (textureHeight*scale)/2

		rl.DrawTextureEx(f.Texture.Data, rl.Vector2{X: x, Y: y}, 0, scale, rl.White)
	}
}

func (f *Dirt) Draw() {
	// TODO
}

func NewFood(x, y float32, foodTexture *Texture) *Food {
	return &Food{
		X:         x,
		Y:         y,
		Eaten:     false,
		Texture:   foodTexture,
		SpawnTime: time.Now(),
		Energy:    float32(rand.Intn(MaxFoodEnergy-MinFoodEnergy+1) + MinFoodEnergy),
	}
}
