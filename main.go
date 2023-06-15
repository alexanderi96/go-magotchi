package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	world           *World
	animationTicker *time.Ticker
	timeTicker      *time.Ticker
	movementTicker  *time.Ticker
	slidePosition   float32
)

func init() {
	rl.InitWindow(screenWidth, screenHeight, "Go-Magotchi")

	world = &World{
		Foods:    []*Food{},
		Dirts:    []*Dirt{},
		cellSize: 40,
		Paused:   false,
	}

	world.WorldHeight = screenHeight - world.cellSize
	world.WorldWidth = screenWidth

	world.Pet = &Pet{
		X:         float32(world.WorldWidth / world.cellSize / 2 * world.cellSize),
		Y:         float32(world.cellSize) + float32(world.WorldHeight/world.cellSize/2*world.cellSize),
		Health:    100,
		Hunger:    50,
		Happiness: 50,
		Energy:    100,
		Age:       0,
		FrameIdx:  0,
	}

	for i := 1; i <= 5; i++ {
		//https://elthen.itch.io/2d-pixel-art-fox-sprites
		texture := rl.LoadTexture(fmt.Sprintf("asset/pet/fox/still/%d.png", i))
		world.Pet.Textures.IdleTextures = append(world.Pet.Textures.IdleTextures, texture)
	}

	for i := 1; i <= 8; i++ {
		//https://elthen.itch.io/2d-pixel-art-fox-sprites
		texture := rl.LoadTexture(fmt.Sprintf("asset/pet/fox/walking/%d.png", i))
		world.Pet.Textures.MovingTextures = append(world.Pet.Textures.MovingTextures, texture)
	}

	// https://lucapixel.itch.io/free-food-pixel-art-45-icons
	foodTexture = rl.LoadTexture("asset/food/food.png")

	// https://www.flaticon.com/free-icon/pharmacy_2695914
	healthIcon = rl.LoadTexture("asset/hud/health.png")
	// https://www.flaticon.com/free-icon/hunger_4968451
	hungerIcon = rl.LoadTexture("asset/hud/hunger.png")
	// https://www.flaticon.com/free-icon/happy_1023758
	happinessIcon = rl.LoadTexture("asset/hud/happiness.png")
	// https://www.flaticon.com/free-icon/flash_2511629
	energyIcon = rl.LoadTexture("asset/hud/energy.png")
	// https://www.flaticon.com/free-icon/time_3240587
	ageIcon = rl.LoadTexture("asset/hud/age.png")
}
func main() {
	defer rl.CloseWindow()

	gameLoop()

	performCloseTasks()

	rl.CloseWindow()
}

func gameLoop() {
	rl.SetTargetFPS(30)

	foodTicker = time.NewTicker(10 * time.Second)
	animationTicker = time.NewTicker(150 * time.Millisecond)
	timeTicker = time.NewTicker(time.Second)
	movementTicker = time.NewTicker(time.Second)

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeyP) {
			slidePosition = float32(rl.GetScreenHeight())
			world.Paused = !world.Paused
		}

		if !world.Paused {
			world.Update()
		}

		world.Draw()
	}
}

func performCloseTasks() {
	for _, texture := range world.Pet.Textures.IdleTextures {
		rl.UnloadTexture(texture)
	}
	for _, texture := range world.Pet.Textures.MovingTextures {
		rl.UnloadTexture(texture)
	}
	rl.UnloadTexture(foodTexture)

	rl.UnloadTexture(healthIcon)
	rl.UnloadTexture(hungerIcon)
	rl.UnloadTexture(happinessIcon)
	rl.UnloadTexture(energyIcon)
	rl.UnloadTexture(ageIcon)
}
