package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth       = 400
	screenHeight      = 240
	foodSpawnInterval = 10
	maxFood           = 10
	foodSize          = 10
)

var (
	statsAreaWidth = int32(float32(screenWidth) * 1 / 4)
	gameAreaWidth  = int32(screenWidth) - statsAreaWidth
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Go-Magotchi")
	defer rl.CloseWindow()

	pet := &Pet{
		X:               float32(gameAreaWidth / 2),
		Y:               float32(screenHeight / 2),
		Health:          100,
		Hunger:          0,
		Happiness:       50,
		Energy:          50,
		Age:             3,
		FrameIdx:        0,
		Foods:           []Food{},
		LastFoodSpawned: time.Now(),
	}

	for i := 1; i <= 5; i++ {
		//https://elthen.itch.io/2d-pixel-art-fox-sprites
		texture := rl.LoadTexture(fmt.Sprintf("asset/pet/fox/still/%d.png", i))
		pet.Textures = append(pet.Textures, texture)
	}

	gameLoop(pet)

	for _, texture := range pet.Textures {
		rl.UnloadTexture(texture)
	}

	rl.CloseWindow()
}

func gameLoop(p *Pet) {
	rl.SetTargetFPS(30)
	// https://lucapixel.itch.io/free-food-pixel-art-45-icons
	foodTexture := rl.LoadTexture("asset/food/food.png")
	if foodTexture.ID == 0 {
		fmt.Println("Errore durante il caricamento della texture del cibo")
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		p.SpawnFood(foodTexture)
		p.MoveToFood()
		p.DrawFoods()

		p.Update()
		p.MoveUserInput()
		p.Draw()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
