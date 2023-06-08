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
)

var (
	world           *World
	animationTicker *time.Ticker
	timeTicker      *time.Ticker
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Go-Magotchi")
	defer rl.CloseWindow()
	world = &World{
		Foods: []*Food{},
		Dirts: []*Dirt{},
		Pet: &Pet{
			X:         float32(screenWidth / 2),
			Y:         float32(screenHeight / 2),
			Health:    100,
			Hunger:    50,
			Happiness: 50,
			Energy:    100,
			Age:       0,
			FrameIdx:  0,
		},
		WorldWidth:  screenWidth / 10,
		WorldHeight: screenHeight / 10,
	}

	world.WorldGrid = make([][]Cell, world.WorldHeight)
	for i := range world.WorldGrid {
		world.WorldGrid[i] = make([]Cell, world.WorldWidth)
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

	gameLoop()

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

	rl.CloseWindow()
}

func gameLoop() {
	rl.SetTargetFPS(30)

	foodTicker = time.NewTicker(10 * time.Second)
	animationTicker = time.NewTicker(150 * time.Millisecond)
	timeTicker = time.NewTicker(time.Second)

	for !rl.WindowShouldClose() {

		world.Update()

		world.Draw()

	}

	rl.CloseWindow()
}
