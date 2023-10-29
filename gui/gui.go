package gui

import (
	"embed"
	"fmt"
	"log"

	"github.com/alexanderi96/go-magotchi/engine"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type texture struct {
	data rl.Texture2D
	path string
}

type guiTextureSet struct {
	// https://www.flaticon.com/free-icon/pharmacy_2695914
	healthIcon *texture
	// https://www.flaticon.com/free-icon/hunger_4968451
	hungerIcon *texture
	// https://www.flaticon.com/free-icon/happy_1023758
	happinessIcon *texture
	// https://www.flaticon.com/free-icon/flash_2511629
	energyIcon *texture
	// https://www.flaticon.com/free-icon/time_3240587
	ageIcon *texture

	dirtTexture *texture

	// https://lucapixel.itch.io/free-food-pixel-art-45-icons
	foodTexture *texture

	//https://elthen.itch.io/2d-pixel-art-fox-sprites
	idlePetFrames     []*texture
	movingPetFrames   []*texture
	sleepingPetFrames []*texture

	selectedPetFrames []*texture
}

type GuiContext struct {
	guiTextureSet *guiTextureSet

	slidePosition float32
}

// Beepberry screen dimensions: 400x240
var (
	//go:embed assets/*
	assets embed.FS
)

func NewGuiContext() (*GuiContext, error) {
	c := &GuiContext{
		slidePosition: float32(rl.GetScreenHeight()),
		guiTextureSet: &guiTextureSet{},
	}

	if err := c.loadTextures(); err != nil {
		return nil, err
	}

	return c, nil
}

func loadFrame(path string) (rl.Texture2D, error) {
	texturedata, err := assets.ReadFile(path)
	log.Println("loading animation: ", path)
	if err != nil {
		return rl.Texture2D{}, fmt.Errorf("error loading animation %s: %x", path, err)
	}
	log.Println("animation loaded")

	return rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", texturedata, int32(len(texturedata)))), nil
}

func (c *GuiContext) animatePet(p *engine.Pet) {

	if p.Moving {

		c.guiTextureSet.selectedPetFrames = c.guiTextureSet.movingPetFrames
		// } else if p.Sleeping {
		// 	c.guiTextureSet.selectedPetFrames = t.SleepingFrames
	} else {
		c.guiTextureSet.selectedPetFrames = c.guiTextureSet.idlePetFrames
	}

	p.FrameIdx = (p.FrameIdx + 1) % len(c.guiTextureSet.selectedPetFrames)

}

func loadFrames(t []*texture, path string) error {
	log.Println("loading animations: ", path)
	framepath, err := assets.ReadDir(path)
	if err != nil {
		return fmt.Errorf("invalid path: %x", err)
	}
	for i := 1; i <= len(framepath); i++ {
		log.Println("loading frame: ", i)

		animationFramepath := fmt.Sprintf(path+"/%d.png", i)
		if frame, err := loadFrame(animationFramepath); err != nil {
			return fmt.Errorf("error loading animations: %x", err)
		} else {
			t = append(t, &texture{data: frame, path: animationFramepath})
		}
		log.Println("appending frame to the list")

	}
	return nil
}

func (c *GuiContext) loadTextures() error {

	loadFrames(c.guiTextureSet.idlePetFrames, "assets/pet/fox/still")
	loadFrames(c.guiTextureSet.movingPetFrames, "assets/pet/fox/walking")

	// for i := 1; i <= 5; i++ {
	// 	texture := rl.LoadTexture(fmt.Sprintf("assets/pet/fox/sleeping/%d.png", i))
	// 	c.guiTextureSet.sleepingPetFrames = append(c.guiTextureSet.sleepingPetFrames, texture)
	// }
	var err error
	if c.guiTextureSet.foodTexture.data, err = loadFrame("assets/food/food.png"); err != nil {
		return fmt.Errorf("error loading texture %x", err)
	}
	if c.guiTextureSet.healthIcon.data, err = loadFrame("assets/hud/health.png"); err != nil {
		return fmt.Errorf("error loading texture %x", err)
	}
	if c.guiTextureSet.hungerIcon.data, err = loadFrame("assets/hud/hunger.png"); err != nil {
		return fmt.Errorf("error loading texture %x", err)
	}
	if c.guiTextureSet.happinessIcon.data, err = loadFrame("assets/hud/happiness.png"); err != nil {
		return fmt.Errorf("error loading texture %x", err)
	}
	if c.guiTextureSet.energyIcon.data, err = loadFrame("assets/hud/energy.png"); err != nil {
		return fmt.Errorf("error loading texture %x", err)
	}
	if c.guiTextureSet.ageIcon.data, err = loadFrame("assets/hud/age.png"); err != nil {
		return fmt.Errorf("error loading texture %x", err)
	}

	return nil
}

func (c *GuiContext) UnloadTextures(w *engine.World) {
	for _, frame := range c.guiTextureSet.idlePetFrames {
		rl.UnloadTexture(frame.data)
	}
	for _, frame := range c.guiTextureSet.movingPetFrames {
		rl.UnloadTexture(frame.data)
	}
	for _, frame := range c.guiTextureSet.sleepingPetFrames {
		rl.UnloadTexture(frame.data)
	}
	rl.UnloadTexture(c.guiTextureSet.foodTexture.data)

	rl.UnloadTexture(c.guiTextureSet.healthIcon.data)
	rl.UnloadTexture(c.guiTextureSet.hungerIcon.data)
	rl.UnloadTexture(c.guiTextureSet.happinessIcon.data)
	rl.UnloadTexture(c.guiTextureSet.energyIcon.data)
	rl.UnloadTexture(c.guiTextureSet.ageIcon.data)
}

func (c *GuiContext) drawStats(w *engine.World) {
	statsY := int32(10)
	statsX := int32(10)

	rl.DrawTextureV(c.guiTextureSet.healthIcon.data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += c.guiTextureSet.healthIcon.data.Width + 5

	stats := fmt.Sprintf(":%.0f", w.Pet.Health)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(c.guiTextureSet.hungerIcon.data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += c.guiTextureSet.hungerIcon.data.Width + 5

	stats = fmt.Sprintf(":%.0f", w.Pet.Hunger)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(c.guiTextureSet.happinessIcon.data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += c.guiTextureSet.happinessIcon.data.Width + 5

	stats = fmt.Sprintf(":%.0f", w.Pet.Happiness)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(c.guiTextureSet.energyIcon.data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += c.guiTextureSet.energyIcon.data.Width + 5

	stats = fmt.Sprintf(":%.0f", w.Pet.Energy)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)
	statsX += 40

	rl.DrawTextureV(c.guiTextureSet.ageIcon.data, rl.NewVector2(float32(statsX), float32(statsY)), rl.White)
	statsX += c.guiTextureSet.ageIcon.data.Width + 5

	stats = fmt.Sprintf(":%d", w.Pet.Age)
	rl.DrawText(stats, statsX, statsY, 20, rl.Black)

}

func drawFood(x, y int32, f *engine.Food, t *texture) {
	if f.X >= float32(x) || f.X < float32(0) || f.Y > float32(y) || f.Y < float32(0) {
		log.Printf("food out of bounds: %f %f ", f.X, f.Y)
		log.Printf("viewport bounds: %f %f ", float32(x), float32(y))
	} else if !f.Eaten {
		textureWidth := float32(t.data.Width)
		textureHeight := float32(t.data.Height)
		scale := float32(f.Energy) / f.Size()
		x := f.X - (textureWidth*scale)/2
		y := f.Y - (textureHeight*scale)/2

		rl.DrawTextureEx(t.data, rl.Vector2{X: x, Y: y}, 0, scale, rl.White)
	}
}

func drawDirt(x, y int32, d *engine.Dirt, t *texture) {
	if d.X >= float32(x) || d.X < float32(0) || d.Y > float32(y) || d.Y < float32(0) {
		log.Printf("dirt out of bounds: %f %f ", d.X, d.Y)
		log.Printf("viewport bounds: %f %f ", float32(x), float32(y))
	} else {
		textureWidth := float32(t.data.Width)
		textureHeight := float32(t.data.Height)
		scale := d.Size()
		x := d.X - (textureWidth*scale)/2
		y := d.Y - (textureHeight*scale)/2

		rl.DrawTextureEx(t.data, rl.Vector2{X: x, Y: y}, 0, scale, rl.White)
	}
}

func drawFloor(x, y int32) {
	rl.DrawRectangle(0, 0, x, y, rl.DarkGreen)
}

func (c *GuiContext) drawFoods(w *engine.World) {
	for _, food := range w.Foods {
		drawFood(w.Config.ViewportX, w.Config.ViewportY, food, c.guiTextureSet.foodTexture)
	}
}

func (c *GuiContext) drawDirts(w *engine.World) {
	for _, dirt := range w.Dirts {
		drawDirt(w.Config.ViewportX, w.Config.ViewportY, dirt, c.guiTextureSet.dirtTexture)
	}
}

func (c *GuiContext) drawPet(w *engine.World) {
	if c.guiTextureSet.selectedPetFrames == nil {
		c.guiTextureSet.selectedPetFrames = c.guiTextureSet.idlePetFrames
	}

	scale := float32(1)
	maxScale := float32(3)
	maxAge := float32(60)

	if w.Pet.Age <= int(maxAge) {
		scale = 1 + ((maxScale - 1) * (float32(w.Pet.Age) / maxAge))
	} else {
		scale = maxScale
	}

	textureWidth := float32(c.guiTextureSet.selectedPetFrames[w.Pet.FrameIdx].data.Width)
	textureHeight := float32(c.guiTextureSet.selectedPetFrames[w.Pet.FrameIdx].data.Height)

	sourceRec := rl.NewRectangle(0, 0, textureWidth, textureHeight)
	if w.Pet.FlipSprite {
		sourceRec.Width *= -1
	}

	destRec := rl.NewRectangle(w.Pet.X-textureWidth*scale/2, w.Pet.Y-textureHeight*scale/2, textureWidth*scale, textureHeight*scale)

	rl.DrawTexturePro(c.guiTextureSet.selectedPetFrames[w.Pet.FrameIdx].data, sourceRec, destRec, rl.NewVector2(0, 0), 0, rl.White)

	c.checkPetFoodCollision(w)
}

func (c *GuiContext) checkPetFoodCollision(w *engine.World) {
	closestFoodIdx := w.Pet.GetClosestFoodDistance(w.Config.ViewportX, w.Config.ViewportY, w.Foods)
	food := w.Foods[closestFoodIdx]
	if rl.CheckCollisionRecs(rl.NewRectangle(w.Pet.X, w.Pet.Y, float32(c.guiTextureSet.selectedPetFrames[w.Pet.FrameIdx].data.Width), float32(c.guiTextureSet.selectedPetFrames[w.Pet.FrameIdx].data.Height)), rl.NewRectangle(food.X, food.Y, float32(c.guiTextureSet.foodTexture.data.Width), float32(c.guiTextureSet.foodTexture.data.Height))) {
		food.Eaten = true
		w.Foods = append(w.Foods[:closestFoodIdx], w.Foods[closestFoodIdx+1:]...)
		w.Pet.Energy += food.Energy
	}
}

func (c *GuiContext) Draw(w *engine.World) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	drawFloor(w.Config.WindowHeight, w.Config.WindowWidth)
	//w.DrawCells()
	c.drawPet(w)
	select {
	case <-w.AnimationTicker.C:
		c.animatePet(w.Pet)
	}
	c.drawFoods(w)
	//w.DrawDirts()
	c.drawStats(w)
	if w.Paused {
		c.drawPauseMenu(w)
	}
	rl.EndDrawing()
}

func (c *GuiContext) drawPauseMenu(w *engine.World) error {
	// Define the pause menu buttons
	screenWidth := float32(rl.GetScreenWidth())
	buttonPadding := float32(10)
	buttonWidth := float32((screenWidth / 3) - buttonPadding)
	buttonHeight := float32(60)

	button1Rect := rl.NewRectangle(buttonPadding, c.slidePosition, buttonWidth, buttonHeight)
	button2Rect := rl.NewRectangle(screenWidth/2-buttonWidth/2, c.slidePosition, buttonWidth, buttonHeight)
	button3Rect := rl.NewRectangle(screenWidth-buttonWidth-buttonPadding, c.slidePosition, buttonWidth, buttonHeight)

	slideSpeed := float32(5)

	if c.slidePosition > float32(rl.GetScreenHeight())-1/float32(rl.GetScreenHeight())-buttonHeight-buttonPadding/2 {
		c.slidePosition -= slideSpeed
	}

	button1Rect.Y = c.slidePosition
	button2Rect.Y = c.slidePosition
	button3Rect.Y = c.slidePosition

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
			c.slidePosition = float32(rl.GetScreenHeight())
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
