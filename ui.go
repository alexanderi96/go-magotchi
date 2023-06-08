package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// func DrawStats() {
// 	// Draw the pet's stats in the stats area
// 	statsY := int32(3) // Start drawing stats 10 pixels from the top
// 	statsX := int32(3)
// 	// Create the stats text
// 	stats := fmt.Sprintf("H:%d H:%d, H:%d, E:%d, A:%d", world.Pet.Health, world.Pet.Hunger, world.Pet.Happiness, world.Pet.Energy, world.Pet.Age)

// 	rl.DrawText(stats, statsX, statsY, 15, rl.Black)
// }

var (
	healthIcon    rl.Texture2D
	hungerIcon    rl.Texture2D
	happinessIcon rl.Texture2D
	energyIcon    rl.Texture2D
	ageIcon       rl.Texture2D
)

func DrawStats() {
	statsY := int32(3) // Start drawing stats 10 pixels from the top
	statsX := int32(3)

	rl.DrawTextureV(healthIcon, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += healthIcon.Width + 5

	stats := fmt.Sprintf(":%d", world.Pet.Health)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(hungerIcon, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += hungerIcon.Width + 5

	stats = fmt.Sprintf(":%d", world.Pet.Hunger)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(happinessIcon, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += happinessIcon.Width + 5

	stats = fmt.Sprintf(":%d", world.Pet.Happiness)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(energyIcon, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += energyIcon.Width + 5

	stats = fmt.Sprintf(":%d", world.Pet.Energy)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(ageIcon, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += ageIcon.Width + 5

	stats = fmt.Sprintf(":%d", world.Pet.Age)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)

}

func DrawFloor() {
	rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Green) // Game area
}
