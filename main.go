package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type World struct {
	Foods []*Food
	Dirts []*Dirt
	Pet   *Pet
}

func (w *World) SpawnFood() {
	newFood := Food{
		X:         float32(int(statsAreaWidth) + rand.Intn(int(gameAreaWidth-foodSize))),
		Y:         float32(rand.Intn(screenHeight - foodSize)),
		Eaten:     false,
		Texture:   foodTexture,
		SpawnTime: time.Now(),
		Energy:    rand.Intn(20),
	}
	w.Foods = append(w.Foods, &newFood)
}

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
	foodTexture    rl.Texture2D
	world          *World
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Go-Magotchi")
	defer rl.CloseWindow()
	world = &World{
		Foods: []*Food{},
		Dirts: []*Dirt{},
		Pet: &Pet{
			X:         float32(gameAreaWidth / 2),
			Y:         float32(screenHeight / 2),
			Health:    100,
			Hunger:    0,
			Happiness: 50,
			Energy:    50,
			Age:       3,
			FrameIdx:  0,
		}}

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
	if foodTexture.ID == 0 {
		fmt.Println("Errore durante il caricamento della texture del cibo")
	}

	gameLoop()

	for _, texture := range world.Pet.Textures.IdleTextures {
		rl.UnloadTexture(texture)
	}
	for _, texture := range world.Pet.Textures.MovingTextures {
		rl.UnloadTexture(texture)
	}
	rl.UnloadTexture(foodTexture)

	rl.CloseWindow()
}

func gameLoop() {
	rl.SetTargetFPS(30)

	foodTicker := time.NewTicker(10 * time.Second)
	animationTicker := time.NewTicker(150 * time.Millisecond)

	for !rl.WindowShouldClose() {

		select {
		case <-foodTicker.C:
			world.SpawnFood()
		case <-animationTicker.C:
			world.Pet.Animate()
		default:
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		// world.SpawnFood()

		world.Pet.Draw()

		for i, food := range world.Foods {
			if food.Eaten {
				world.Foods = append(world.Foods[:i], world.Foods[i+1:]...)
			}

			food.DrawFoods()
		}

		world.Pet.MoveToFood()
		//world.Pet.MoveUserInput()

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
