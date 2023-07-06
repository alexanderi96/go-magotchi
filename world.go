package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type World struct {
	Foods                             []*Food
	Dirts                             []*Dirt
	Pet                               *Pet
	WorldWidth, WorldHeight, cellSize int
	Paused                            bool
}

type Cell struct {
	Obstacle bool
}

func (w *World) SpawnFood() {
	X := float32(rand.Intn((world.WorldWidth / world.cellSize)) * world.cellSize)
	Y := float32(rand.Intn((world.WorldHeight/world.cellSize))*world.cellSize + world.cellSize)

	for _, food := range w.Foods {
		if food.X == X && food.Y == Y {
			return
		}
	}

	w.Foods = append(w.Foods, NewFood(X, Y))
}

func (w *World) SpawnDirt() {
	// TODO: dirt spawn
}

func (w *World) DrawFoods() {
	for _, food := range w.Foods {
		food.Draw()
	}
}

func (w *World) DrawDirts() {
	for _, dirt := range w.Dirts {
		dirt.Draw()
	}
}

func (w *World) DrawPet() {
	w.Pet.Draw()
}

func (w *World) Draw() {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	DrawFloor()
	//w.DrawCells()
	w.DrawPet()
	w.DrawFoods()
	//w.DrawDirts()
	DrawStats()
	if world.Paused {
		w.DrawPauseMenu()
	}
	rl.EndDrawing()
}

func (w *World) DrawCells() {
	for y := 0; y < screenHeight/w.cellSize; y++ {
		for x := 0; x < screenWidth/w.cellSize; x++ {
			rect := rl.NewRectangle(float32(x*w.cellSize), float32(y*w.cellSize), float32(w.cellSize), float32(w.cellSize))
			rl.DrawRectangleLinesEx(rect, 1, rl.Black)
		}
	}
}

func (w *World) DrawPauseMenu() {
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

	if world.Paused && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		mousePos := rl.GetMousePosition()

		// Check if the Resume button was clicked
		if rl.CheckCollisionPointRec(mousePos, button1Rect) {
			world.Paused = false
			slidePosition = float32(rl.GetScreenHeight())
		}

		// Check if the Restart button was clicked
		if rl.CheckCollisionPointRec(mousePos, button2Rect) {
			// Implement your restart logic here
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
				},
				WorldWidth:  int(screenWidth) / 10,
				WorldHeight: screenHeight / 10,
				Paused:      false,
			}
		}

		// Check if the Quit button was clicked
		if rl.CheckCollisionPointRec(mousePos, button3Rect) {
			rl.CloseWindow()
		}
	}
}

func (w *World) Update() {
	select {
	case <-foodTicker.C:
		if !world.Pet.Dead {
			world.SpawnFood()
		}

	case <-animationTicker.C:
		world.Pet.Animate()

	case <-timeTicker.C:
		if !world.Pet.Dead {
			world.Pet.GetOlder()
		}

	case <-movementTicker.C:
		if world.Pet.WantToMove() {
			world.Pet.MoveToFood()
		}
	default:
	}
	world.Pet.Update()
}
