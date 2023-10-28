package gui

import (
	"fmt"

	"github.com/alexanderi96/go-magotchi/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Beepberry screen dimensions: 400x240
var (
	slidePosition float32
)

func init() {
	slidePosition = float32(rl.GetScreenHeight())
}

func UnloadTextures(w *engine.World) {
	for _, texture := range w.Textures.IdleTextures {
		rl.UnloadTexture(texture.Data)
	}
	for _, texture := range w.Textures.MovingTextures {
		rl.UnloadTexture(texture.Data)
	}
	for _, texture := range w.Textures.SleepingTextures {
		rl.UnloadTexture(texture.Data)
	}
	rl.UnloadTexture(w.Textures.FoodTexture.Data)

	rl.UnloadTexture(w.Textures.HealthIcon.Data)
	rl.UnloadTexture(w.Textures.HungerIcon.Data)
	rl.UnloadTexture(w.Textures.HappinessIcon.Data)
	rl.UnloadTexture(w.Textures.EnergyIcon.Data)
	rl.UnloadTexture(w.Textures.AgeIcon.Data)
}

func drawStats(w *engine.World) {
	statsY := int32(10)
	statsX := int32(10)

	rl.DrawTextureV(w.Textures.HealthIcon.Data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += w.Textures.HealthIcon.Data.Width + 5

	stats := fmt.Sprintf(":%.0f", w.Pet.Health)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(w.Textures.HungerIcon.Data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += w.Textures.HungerIcon.Data.Width + 5

	stats = fmt.Sprintf(":%.0f", w.Pet.Hunger)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(w.Textures.HappinessIcon.Data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += w.Textures.HappinessIcon.Data.Width + 5

	stats = fmt.Sprintf(":%.0f", w.Pet.Happiness)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(w.Textures.EnergyIcon.Data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += w.Textures.EnergyIcon.Data.Width + 5

	stats = fmt.Sprintf(":%.0f", w.Pet.Energy)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(w.Textures.AgeIcon.Data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += w.Textures.AgeIcon.Data.Width + 5

	stats = fmt.Sprintf(":%d", w.Pet.Age)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)

}

func drawFloor(x, y int32) {
	rl.DrawRectangle(0, 0, x, y, rl.DarkGreen)
}

func drawFoods(w *engine.World) {
	for _, food := range w.Foods {
		food.Draw(w.Config.ViewportX, w.Config.ViewportY)
	}
}

func drawDirts(w *engine.World) {
	for _, dirt := range w.Dirts {
		dirt.Draw()
	}
}

func drawPet(w *engine.World) {
	if w.Pet.SelectedTextures == nil {
		w.Pet.SelectedTextures = w.Textures.IdleTextures
	}

	scale := float32(1)
	maxScale := float32(3)
	maxAge := float32(60)

	if w.Pet.Age <= int(maxAge) {
		scale = 1 + ((maxScale - 1) * (float32(w.Pet.Age) / maxAge))
	} else {
		scale = maxScale
	}

	textureWidth := float32(w.Pet.SelectedTextures[w.Pet.FrameIdx].Data.Width)
	textureHeight := float32(w.Pet.SelectedTextures[w.Pet.FrameIdx].Data.Height)

	sourceRec := rl.NewRectangle(0, 0, textureWidth, textureHeight)
	if w.Pet.FlipSprite {
		sourceRec.Width *= -1
	}

	destRec := rl.NewRectangle(w.Pet.X-textureWidth*scale/2, w.Pet.Y-textureHeight*scale/2, textureWidth*scale, textureHeight*scale)

	rl.DrawTexturePro(w.Pet.SelectedTextures[w.Pet.FrameIdx].Data, sourceRec, destRec, rl.NewVector2(0, 0), 0, rl.White)
}

func Draw(w *engine.World) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	drawFloor(w.Config.WindowHeight, w.Config.WindowWidth)
	//w.DrawCells()
	drawPet(w)
	drawFoods(w)
	//w.DrawDirts()
	drawStats(w)
	if w.Paused {
		drawPauseMenu(w)
	}
	rl.EndDrawing()
}

func drawPauseMenu(w *engine.World) error {
	// Define the pause menu buttons
	screenWidth := float32(rl.GetScreenWidth())
	buttonPadding := float32(10)
	buttonWidth := float32((screenWidth / 3) - buttonPadding)
	buttonHeight := float32(60)

	button1Rect := rl.NewRectangle(buttonPadding, slidePosition, buttonWidth, buttonHeight)
	button2Rect := rl.NewRectangle(screenWidth/2-buttonWidth/2, slidePosition, buttonWidth, buttonHeight)
	button3Rect := rl.NewRectangle(screenWidth-buttonWidth-buttonPadding, slidePosition, buttonWidth, buttonHeight)

	slideSpeed := float32(5)

	if slidePosition > float32(rl.GetScreenHeight())-1/float32(rl.GetScreenHeight())-buttonHeight-buttonPadding/2 {
		slidePosition -= slideSpeed
	}

	button1Rect.Y = slidePosition
	button2Rect.Y = slidePosition
	button3Rect.Y = slidePosition

	rl.BeginDrawing()
	// Draw the pause menu buttons
	// Calculate the center position of the button
	button1TextX := button1Rect.X + ((button1Rect.Width)-float32(rl.MeasureText("Resume", 20)))/2
	button1TextY := button1Rect.Y + (button1Rect.Height-20)/2

	button2TextX := button2Rect.X + (button2Rect.Width-float32(rl.MeasureText("Restart", 20)))/2
	button2TextY := button2Rect.Y + (button2Rect.Height-20)/2

	button3TextX := button3Rect.X + (button3Rect.Width-float32(rl.MeasureText("Quit", 20)))/2
	button3TextY := button3Rect.Y + (button3Rect.Height-20)/2

	rl.DrawRectangleRec(button1Rect, rl.LightGray)
	rl.DrawRectangleRec(button2Rect, rl.LightGray)
	rl.DrawRectangleRec(button3Rect, rl.LightGray)

	// Draw the text on the buttons at the center
	rl.DrawText("Resume", int32(button1TextX), int32(button1TextY), 20, rl.Black)
	rl.DrawText("Restart", int32(button2TextX), int32(button2TextY), 20, rl.Black)
	rl.DrawText("Quit", int32(button3TextX), int32(button3TextY), 20, rl.Black)

	if w.Paused && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		mousePos := rl.GetMousePosition()

		// Check if the Resume button was clicked
		if rl.CheckCollisionPointRec(mousePos, button1Rect) {
			w.Paused = false
			slidePosition = float32(rl.GetScreenHeight())
		}

		// Check if the Restart button was clicked
		if rl.CheckCollisionPointRec(mousePos, button2Rect) {
			// Implement your restart logic here
			var err error
			if w, err = engine.NewGame(); err != nil {
				return fmt.Errorf("error resetting game", err)
			}
		}

		// Check if the Quit button was clicked
		if rl.CheckCollisionPointRec(mousePos, button3Rect) {
			rl.CloseWindow()
		}
	}
	return nil
}
