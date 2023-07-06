package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextureSet struct {
	IdleTextures     []rl.Texture2D
	MovingTextures   []rl.Texture2D
	SleepingTextures []rl.Texture2D
}

type Pet struct {
	X          float32
	Y          float32
	Health     float32
	Hunger     float32
	Happiness  float32
	Energy     float32
	Textures   TextureSet
	FlipSprite bool
	Moving     bool
	Sleeping   bool
	Dead       bool
	FrameIdx   int
	Age        int
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
	if selectedTextures == nil {
		selectedTextures = p.Textures.IdleTextures
	}

	textureWidth := float32(selectedTextures[p.FrameIdx].Width)
	textureHeight := float32(selectedTextures[p.FrameIdx].Height)

	sourceRec := rl.NewRectangle(0, 0, textureWidth, textureHeight)
	if p.FlipSprite {
		sourceRec.Width *= -1
	}

	destRec := rl.NewRectangle(
		float32(p.X),
		float32(p.Y),
		(22*float32(world.WorldWidth))/400,
		(15*float32(world.WorldHeight))/240,
	)

	destRec.X += (float32(world.cellSize) - destRec.Width) / 2
	destRec.Y += (float32(world.cellSize) - destRec.Height) / 2

	rl.DrawTexturePro(selectedTextures[p.FrameIdx], sourceRec, destRec, rl.NewVector2(0, 0), 0, rl.White)
}

func (p *Pet) WantToMove() bool {
	return !p.Dead && p.Hunger > 50 || p.Energy > 50
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

	closestFoodIdx := 0
	closestDistance := float32(world.WorldWidth + world.WorldHeight)

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

	food := world.Foods[closestFoodIdx]
	speed := (float32(world.cellSize) * (float32(world.WorldWidth) / 200)) * p.Energy
	if p.X < food.X {
		p.X += speed
		p.Energy -= 1
		p.FlipSprite = false
	} else if p.X > food.X {
		p.X -= speed
		p.Energy -= 1
		p.FlipSprite = true
	}

	if p.Y < food.Y {
		p.Y += speed
		p.Energy -= 1
	} else if p.Y > food.Y {
		p.Y -= speed
		p.Energy -= 1
	}

	p.Moving = oldX != p.X || oldY != p.Y
	if p.Moving != oldState {
		p.FrameIdx = 0
	}

	if rl.CheckCollisionRecs(rl.NewRectangle(p.X, p.Y, float32(world.cellSize), float32(world.cellSize)), rl.NewRectangle(food.X, food.Y, float32(world.cellSize), float32(world.cellSize))) {
		food.Eaten = true
		world.Foods = append(world.Foods[:closestFoodIdx], world.Foods[closestFoodIdx+1:]...)
		p.Energy += food.Energy
	}
}

func (p *Pet) Animate() {

	if p.Moving {

		selectedTextures = p.Textures.MovingTextures
		// } else if p.Sleeping {
		// 	selectedTextures = p.Textures.SleepingTextures
	} else {
		selectedTextures = p.Textures.IdleTextures
	}

	p.FrameIdx = (p.FrameIdx + 1) % len(selectedTextures)

}

func (p *Pet) GetOlder() {
	p.Age++
}

func (p *Pet) Update() {
	if p.Energy < 0 {
		p.Energy = 0
	} else if p.Energy > 100 {
		p.Energy = 100
	}

	if p.Hunger < 0 {
		p.Hunger = 0
	} else if p.Hunger > 100 {
		p.Hunger = 100
	}

	if p.Energy == 0 || p.Health == 0 {
		p.Dead = true
	}

}
