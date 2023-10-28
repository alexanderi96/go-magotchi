package engine

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Pet struct {
	X                float32
	Y                float32
	Health           float32
	Hunger           float32
	Happiness        float32
	Energy           float32
	SelectedTextures []*Texture
	FlipSprite       bool
	Moving           bool
	Sleeping         bool
	Dead             bool
	FrameIdx         int
	Age              int
}

func (p *Pet) MoveUserInput(x, y float32) {
	moveSpeed := float32(10)

	oldX := p.X
	oldY := p.Y
	oldState := p.Moving

	if rl.IsKeyDown(rl.KeyRight) {
		p.X += moveSpeed
		if p.X > float32(x) {
			p.X = float32(0)
		}
		p.FlipSprite = false
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		p.X -= moveSpeed
		if p.X < float32(0) {
			p.X = float32(x)
		}
		p.FlipSprite = true
	}

	if rl.IsKeyDown(rl.KeyUp) {
		p.Y -= moveSpeed
		if p.Y < 0 {
			p.Y = float32(y)
		}
	}

	if rl.IsKeyDown(rl.KeyDown) {
		p.Y += moveSpeed
		if p.Y > float32(y) {
			p.Y = 0
		}
	}

	p.Moving = oldX != p.X || oldY != p.Y
	if p.Moving != oldState {
		p.FrameIdx = 0
	}
}

func (p *Pet) WantToMove() bool {
	return !p.Dead && p.Hunger > 50 || p.Energy > 0
}

func (p *Pet) MoveToFood(x, y float32, f []*Food) {
	oldX := p.X
	oldY := p.Y
	oldState := p.Moving

	if len(f) == 0 {
		p.Moving = false
		if p.Moving != oldState {
			p.FrameIdx = 0
		}
		return
	}

	closestFoodIdx := 0
	closestDistance := float32(x + y)

	for idx, food := range f {
		if food.Eaten {
			continue
		}

		distance := rl.Vector2Distance(rl.NewVector2(p.X, p.Y), rl.NewVector2(food.X, food.Y))
		if distance < closestDistance {
			closestFoodIdx = idx
			closestDistance = distance
		}
	}

	food := f[closestFoodIdx]
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

	if rl.CheckCollisionRecs(rl.NewRectangle(p.X, p.Y, float32(p.SelectedTextures[p.FrameIdx].Data.Width), float32(p.SelectedTextures[p.FrameIdx].Data.Height)), rl.NewRectangle(food.X, food.Y, float32(food.Texture.Data.Width), float32(food.Texture.Data.Height))) {
		food.Eaten = true
		f = append(f[:closestFoodIdx], f[closestFoodIdx+1:]...)
		p.Energy += food.Energy
	}
}

func (p *Pet) Animate(t *TextureSet) {

	if p.Moving {

		p.SelectedTextures = append(t.MovingTextures)
		// } else if p.Sleeping {
		// 	p.SelectedTextures = t.SleepingTextures
	} else {
		p.SelectedTextures = t.IdleTextures
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
