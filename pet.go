package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextureSet struct {
	IdleTextures   []rl.Texture2D
	MovingTextures []rl.Texture2D
}

type Pet struct {
	X, Y                                        float32
	Health, Hunger, Happiness, Energy, FrameIdx int
	Textures                                    TextureSet
	FlipSprite, Moving                          bool
	Age                                         int
}

var selectedTextures []rl.Texture2D

func (p *Pet) MoveUserInput() {
	moveSpeed := float32(10)

	oldX := p.X
	oldY := p.Y
	oldState := p.Moving

	if rl.IsKeyDown(rl.KeyRight) {
		p.X += moveSpeed
		if p.X > float32(world.WorldWidth) {
			p.X = float32(0)
		}
		p.FlipSprite = false
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		p.X -= moveSpeed
		if p.X < float32(0) {
			p.X = float32(world.WorldWidth)
		}
		p.FlipSprite = true
	}

	if rl.IsKeyDown(rl.KeyUp) {
		p.Y -= moveSpeed
		if p.Y < 0 {
			p.Y = float32(world.WorldHeight)
		}
	}

	if rl.IsKeyDown(rl.KeyDown) {
		p.Y += moveSpeed
		if p.Y > float32(world.WorldHeight) {
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

	if selectedTextures == nil {
		selectedTextures = p.Textures.IdleTextures
	}

	// Draw the pet in the center of the game area
	scale := float32(1)    // Scala iniziale
	maxScale := float32(3) // Dimensione massima
	maxAge := float32(60)  // Età massima per raggiungere la dimensione massima

	if p.Age <= int(maxAge) {
		// Calcola la scala in base all'età
		scale = 1 + ((maxScale - 1) * (float32(p.Age) / maxAge))
	} else {
		scale = maxScale // Mantieni la scala al massimo dopo maxAge
	}

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
	closestDistance := float32(world.WorldWidth + world.WorldHeight) // A value greater than the maximum possible distance

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
		p.Energy--
		p.FlipSprite = false
	} else if p.X > food.X {
		p.X--
		p.Energy--
		p.FlipSprite = true
	}

	if p.Y < food.Y {
		p.Y++
		p.Energy--
	} else if p.Y > food.Y {
		p.Y--
		p.Energy--
	}

	p.Moving = oldX != p.X || oldY != p.Y
	if p.Moving != oldState {
		p.FrameIdx = 0
	}

	// Check if the pet has reached the food
	if rl.Vector2Distance(rl.NewVector2(p.X, p.Y), rl.NewVector2(food.X, food.Y)) < foodSize {
		food.Eaten = true
		world.Foods = append(world.Foods[:closestFoodIdx], world.Foods[closestFoodIdx+1:]...)
		p.Energy += food.Energy
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

func (p *Pet) GetOlder() {
	p.Age++
}
