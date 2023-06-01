package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Pet struct {
	X, Y            float32
	Health, Hunger  int
	Happiness, Age  int
	FrameIdx        int
	Textures        []rl.Texture2D
	FoodNearby      *Food
	Energy          int
	CurrentFoodIdx  int
	Foods           []Food
	LastFoodSpawned time.Time
	FlipSprite      bool
}

func (p *Pet) MoveUserInput() {
	moveSpeed := float32(10) // feel free to adjust this value
	if rl.IsKeyDown(rl.KeyRight) {
		p.X += moveSpeed
		if p.X > float32(screenWidth) {
			p.X = float32(statsAreaWidth)
		}
		p.FlipSprite = false
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		p.X -= moveSpeed
		if p.X < float32(statsAreaWidth) {
			p.X = float32(screenWidth)
		}
		p.FlipSprite = true
	}

	if rl.IsKeyDown(rl.KeyUp) {
		p.Y -= moveSpeed
		if p.Y < 0 {
			p.Y = screenHeight
		}
	}

	if rl.IsKeyDown(rl.KeyDown) {
		p.Y += moveSpeed
		if p.Y > screenHeight {
			p.Y = 0
		}
	}
}

func (p *Pet) Update() {
	// La logica di aggiornamento va qui
}

func (p *Pet) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	// Draw areas
	rl.DrawRectangle(0, 0, statsAreaWidth, screenHeight, rl.LightGray)         // Stats area
	rl.DrawRectangle(statsAreaWidth, 0, gameAreaWidth, screenHeight, rl.White) // Game area

	// Draw the pet in the center of the game area
	scale := float32(p.Age) // Scale the image by the age of the pet
	textureWidth := float32(p.Textures[p.FrameIdx].Width)
	textureHeight := float32(p.Textures[p.FrameIdx].Height)

	// Create the source rectangle for the texture
	sourceRec := rl.NewRectangle(0, 0, textureWidth, textureHeight)
	if p.FlipSprite {
		sourceRec.Width *= -1
	}

	// Create the destination rectangle for the texture
	destRec := rl.NewRectangle(p.X-textureWidth*scale/2, p.Y-textureHeight*scale/2, textureWidth*scale, textureHeight*scale)

	// Create the origin vector for the texture
	origin := rl.NewVector2(0, 0)

	// Draw the texture
	rl.DrawTexturePro(p.Textures[p.FrameIdx], sourceRec, destRec, origin, 0, rl.White)

	// Draw the pet's stats in the stats area
	statsY := int32(3)      // Start drawing stats 10 pixels from the top
	lineHeight := int32(15) // Each line of text is 30 pixels high
	rl.DrawText(fmt.Sprintf("Health: %d", p.Health), 3, statsY, 15, rl.Black)
	rl.DrawText(fmt.Sprintf("Hunger: %d", p.Hunger), 3, statsY+lineHeight, 15, rl.Black)
	rl.DrawText(fmt.Sprintf("Happiness: %d", p.Happiness), 3, statsY+2*lineHeight, 15, rl.Black)
	rl.DrawText(fmt.Sprintf("Energy: %d", p.Energy), 3, statsY+3*lineHeight, 15, rl.Black)
	rl.DrawText(fmt.Sprintf("Age: %d", p.Age), 3, statsY+4*lineHeight, 15, rl.Black)

	p.FrameIdx = (p.FrameIdx + 1) % len(p.Textures) // Advance to the next frame

	rl.EndDrawing()
}

func (p *Pet) SpawnFood(foodTexture rl.Texture2D) {
	if len(p.Foods) < maxFood && time.Since(p.LastFoodSpawned).Seconds() > foodSpawnInterval {
		newFood := Food{
			X:       float32(int(statsAreaWidth) + rand.Intn(int(gameAreaWidth-foodSize))),
			Y:       float32(rand.Intn(screenHeight - foodSize)),
			Eaten:   false,
			Texture: foodTexture,
		}
		p.Foods = append(p.Foods, newFood)
		p.LastFoodSpawned = time.Now()
	}
}

func (p *Pet) MoveToFood() {
	if len(p.Foods) == 0 {
		return
	}

	// Get the closest food
	closestFoodIdx := 0
	closestDistance := float32(gameAreaWidth + screenHeight) // A value greater than the maximum possible distance

	for idx, food := range p.Foods {
		if food.Eaten {
			continue
		}

		distance := rl.Vector2Distance(rl.NewVector2(p.X, p.Y), rl.NewVector2(food.X, food.Y))
		if distance < closestDistance {
			closestFoodIdx = idx
			closestDistance = distance
		}
	}

	// Move to closest food
	food := &p.Foods[closestFoodIdx]
	if p.X < food.X {
		p.X++
		p.FlipSprite = false
	} else if p.X > food.X {
		p.X--
		p.FlipSprite = true
	}

	if p.Y < food.Y {
		p.Y++
	} else if p.Y > food.Y {
		p.Y--
	}

	// Check if the pet has reached the food
	if rl.Vector2Distance(rl.NewVector2(p.X, p.Y), rl.NewVector2(food.X, food.Y)) < foodSize {
		food.Eaten = true
		p.Health += 10
	}
}

func (p *Pet) DrawFoods() {
	for _, food := range p.Foods {
		if !food.Eaten {
			textureWidth := float32(food.Texture.Width)
			textureHeight := float32(food.Texture.Height)
			scale := float32(foodSize) / textureWidth
			x := food.X - (textureWidth*scale)/2
			y := food.Y - (textureHeight*scale)/2

			rl.DrawTextureEx(food.Texture, rl.NewVector2(x, y), 0, scale, rl.White)
		}
	}
}
