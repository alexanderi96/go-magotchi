package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextureSet struct {
	IdleTextures   []rl.Texture2D
	MovingTextures []rl.Texture2D
}

type Pet struct {
	X, Y                                             float32
	Health, Hunger, Happiness, Age, Energy, FrameIdx int
	Textures                                         TextureSet
	FlipSprite, Moving                               bool
}

var selectedTextures []rl.Texture2D

func (p *Pet) MoveUserInput() {
	moveSpeed := float32(10) // feel free to adjust this value

	oldX := p.X
	oldY := p.Y
	oldState := p.Moving

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

	p.Moving = oldX != p.X || oldY != p.Y
	if p.Moving != oldState {
		p.FrameIdx = 0
	}
}

func (p *Pet) Draw() {
	// Draw areas
	rl.DrawRectangle(0, 0, statsAreaWidth, screenHeight, rl.LightGray)         // Stats area
	rl.DrawRectangle(statsAreaWidth, 0, gameAreaWidth, screenHeight, rl.White) // Game area

	if selectedTextures == nil {
		selectedTextures = p.Textures.IdleTextures
	}

	// Draw the pet in the center of the game area
	scale := float32(p.Age) // Scale the image by the age of the pet
	textureWidth := float32(selectedTextures[p.FrameIdx].Width)
	textureHeight := float32(selectedTextures[p.FrameIdx].Height)

	// Create the source rectangle for the texture
	sourceRec := rl.NewRectangle(0, 0, textureWidth, textureHeight)
	if p.FlipSprite {
		sourceRec.Width *= -1
	}

	// Create the destination rectangle for the texture
	destRec := rl.NewRectangle(p.X-textureWidth*scale/2, p.Y-textureHeight*scale/2, textureWidth*scale, textureHeight*scale)

	// Draw the texture
	rl.DrawTexturePro(selectedTextures[p.FrameIdx], sourceRec, destRec, rl.NewVector2(0, 0), 0, rl.White)

	// Draw the pet's stats in the stats area
	statsY := int32(3)      // Start drawing stats 10 pixels from the top
	lineHeight := int32(15) // Each line of text is 30 pixels high
	rl.DrawText(fmt.Sprintf("Health: %d", p.Health), 3, statsY, 15, rl.Black)
	rl.DrawText(fmt.Sprintf("Hunger: %d", p.Hunger), 3, statsY+lineHeight, 15, rl.Black)
	rl.DrawText(fmt.Sprintf("Happiness: %d", p.Happiness), 3, statsY+2*lineHeight, 15, rl.Black)
	rl.DrawText(fmt.Sprintf("Energy: %d", p.Energy), 3, statsY+3*lineHeight, 15, rl.Black)
	rl.DrawText(fmt.Sprintf("Age: %d", p.Age), 3, statsY+4*lineHeight, 15, rl.Black)

}

func (p *Pet) MoveToFood() {
	oldX := p.X
	oldY := p.Y
	oldState := p.Moving

	if len(world.Foods) == 0 {
		p.Moving = false
		if p.Moving != oldState {
			p.FrameIdx = 0
		}
		return
	}

	// Get the closest food
	closestFoodIdx := 0
	closestDistance := float32(gameAreaWidth + screenHeight) // A value greater than the maximum possible distance

	for idx, food := range world.Foods {
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
	food := world.Foods[closestFoodIdx]
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

	p.Moving = oldX != p.X || oldY != p.Y
	if p.Moving != oldState {
		p.FrameIdx = 0
	}

	// Check if the pet has reached the food
	if rl.Vector2Distance(rl.NewVector2(p.X, p.Y), rl.NewVector2(food.X, food.Y)) < foodSize {
		food.Eaten = true
		p.Health += food.Energy
	}
}

func (p *Pet) Animate() {

	if p.Moving {
		// log.Println("Moving")
		selectedTextures = p.Textures.MovingTextures
	} else {
		// log.Println("Idle")
		selectedTextures = p.Textures.IdleTextures
	}

	p.FrameIdx = (p.FrameIdx + 1) % len(selectedTextures) // Advance to the next frame

}
