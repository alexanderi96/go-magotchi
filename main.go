package main

import (
	"embed"
	"fmt"
	"log"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/*
var assets embed.FS

var (
	world           *World
	animationTicker *time.Ticker
	timeTicker      *time.Ticker
	movementTicker  *time.Ticker
	slidePosition   float32
)

func init() {
	if len(scrHgt) > 0 && len(scrWdt) > 0 {
		screenHeight, _ = strconv.Atoi(scrHgt)
		screenWidth, _ = strconv.Atoi(scrWdt)
	}

	rl.InitWindow(int32(screenWidth), int32(screenHeight), "Go-Magotchi")

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
			Sleeping:  false,
			Dead:      false,
		},
		WorldWidth:  int(screenWidth),
		WorldHeight: int(screenHeight),
		Paused:      false,
	}

	fox_still_frames, err := assets.ReadDir("assets/pet/fox/still")
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= len(fox_still_frames); i++ {
		//https://elthen.itch.io/2d-pixel-art-fox-sprites
		textureData, err := assets.ReadFile(fmt.Sprintf("assets/pet/fox/still/%d.png", i))
		if err != nil {
			log.Fatal(err)
		}
		texture := rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", textureData, int32(len(textureData))))
		world.Pet.Textures.IdleTextures = append(world.Pet.Textures.IdleTextures, texture)
	}

	fox_walking_frames, err := assets.ReadDir("assets/pet/fox/walking")
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= len(fox_walking_frames); i++ {
		//https://elthen.itch.io/2d-pixel-art-fox-sprites
		textureData, err := assets.ReadFile(fmt.Sprintf("assets/pet/fox/walking/%d.png", i))
		if err != nil {
			log.Fatal(err)
		}
		texture := rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", textureData, int32(len(textureData))))
		world.Pet.Textures.MovingTextures = append(world.Pet.Textures.MovingTextures, texture)
	}

	// for i := 1; i <= 5; i++ {
	// 	texture := rl.LoadTexture(fmt.Sprintf("assets/pet/fox/sleeping/%d.png", i))
	// 	world.Pet.Textures.SleepingTextures = append(world.Pet.Textures.SleepingTextures, texture)
	// }

	// https://lucapixel.itch.io/free-food-pixel-art-45-icons
	foodTextureData, err := assets.ReadFile("assets/food/food.png")
	if err != nil {
		log.Fatal(err)
	}
	foodTexture = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", foodTextureData, int32(len(foodTextureData))))

	// https://www.flaticon.com/free-icon/pharmacy_2695914
	healthIconData, err := assets.ReadFile("assets/hud/health.png")
	if err != nil {
		log.Fatal(err)
	}
	healthIcon = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", healthIconData, int32(len(healthIconData))))

	// https://www.flaticon.com/free-icon/hunger_4968451
	hungerIconData, err := assets.ReadFile("assets/hud/hunger.png")
	if err != nil {
		log.Fatal(err)
	}
	hungerIcon = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", hungerIconData, int32(len(hungerIconData))))

	// https://www.flaticon.com/free-icon/happy_1023758
	happinessIconData, err := assets.ReadFile("assets/hud/happiness.png")
	if err != nil {
		log.Fatal(err)
	}
	happinessIcon = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", happinessIconData, int32(len(happinessIconData))))

	// https://www.flaticon.com/free-icon/flash_2511629
	energyIconData, err := assets.ReadFile("assets/hud/energy.png")
	if err != nil {
		log.Fatal(err)
	}
	energyIcon = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", energyIconData, int32(len(energyIconData))))

	// https://www.flaticon.com/free-icon/time_3240587
	ageIconData, err := assets.ReadFile("assets/hud/age.png")
	if err != nil {
		log.Fatal(err)
	}
	ageIcon = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", ageIconData, int32(len(ageIconData))))
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
	movementTicker = time.NewTicker(time.Millisecond)

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
	for _, texture := range world.Pet.Textures.SleepingTextures {
		rl.UnloadTexture(texture)
	}
	rl.UnloadTexture(foodTexture)

	rl.UnloadTexture(healthIcon)
	rl.UnloadTexture(hungerIcon)
	rl.UnloadTexture(happinessIcon)
	rl.UnloadTexture(energyIcon)
	rl.UnloadTexture(ageIcon)
}
