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
	X                float32
	Y                float32
	Health           float32
	Hunger           float32
	Happiness        float32
	Energy           float32
	Textures         TextureSet
	SelectedTextures []rl.Texture2D
	FlipSprite       bool
	Moving           bool
	Sleeping         bool
	Dead             bool
	FrameIdx         int
	Age              int
}

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
	if p.SelectedTextures == nil {
		p.SelectedTextures = p.Textures.IdleTextures
	}

	scale := float32(1)
	maxScale := float32(3)
	maxAge := float32(60)

	if p.Age <= int(maxAge) {
		scale = 1 + ((maxScale - 1) * (float32(p.Age) / maxAge))
	} else {
		scale = maxScale
	}

	textureWidth := float32(p.SelectedTextures[p.FrameIdx].Width)
	textureHeight := float32(p.SelectedTextures[p.FrameIdx].Height)

	sourceRec := rl.NewRectangle(0, 0, textureWidth, textureHeight)
	if p.FlipSprite {
		sourceRec.Width *= -1
	}

	destRec := rl.NewRectangle(p.X-textureWidth*scale/2, p.Y-textureHeight*scale/2, textureWidth*scale, textureHeight*scale)

	rl.DrawTexturePro(p.SelectedTextures[p.FrameIdx], sourceRec, destRec, rl.NewVector2(0, 0), 0, rl.White)
}

func (p *Pet) WantToMove() bool {
	return !p.Dead && p.Hunger > 50 || p.Energy > 0
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
	if p.X < food.X {
		p.X++
		p.Energy -= 0.1
		p.FlipSprite = false
	} else if p.X > food.X {
		p.X--
		p.Energy -= 0.1
		p.FlipSprite = true
	}

	if p.Y < food.Y {
		p.Y++
		p.Energy -= 0.1
	} else if p.Y > food.Y {
		p.Y--
		p.Energy -= 0.1
	}

	p.Moving = oldX != p.X || oldY != p.Y
	if p.Moving != oldState {
		p.FrameIdx = 0
	}

	if rl.CheckCollisionRecs(rl.NewRectangle(p.X, p.Y, float32(p.SelectedTextures[p.FrameIdx].Width), float32(p.SelectedTextures[p.FrameIdx].Height)), rl.NewRectangle(food.X, food.Y, float32(food.Texture.Width), float32(food.Texture.Height))) {
		food.Eaten = true
		world.Foods = append(world.Foods[:closestFoodIdx], world.Foods[closestFoodIdx+1:]...)
		p.Energy += food.Energy
	}
}

func (p *Pet) Animate() {

	if p.Moving {

		p.SelectedTextures = p.Textures.MovingTextures
		// } else if p.Sleeping {
		// 	p.SelectedTextures = p.Textures.SleepingTextures
	} else {
		p.SelectedTextures = p.Textures.IdleTextures
	}

	p.FrameIdx = (p.FrameIdx + 1) % len(p.SelectedTextures)

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
