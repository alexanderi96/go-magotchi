package main

import (
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type World struct {
	Foods                   []*Food
	Dirts                   []*Dirt
	Pet                     *Pet
	WorldWidth, WorldHeight int
	WorldGrid               [][]Cell
}

type Cell struct {
	Obstacle bool
}

func (w *World) SpawnFood() {
	newFood := Food{
		X:         float32(int(0) + rand.Intn(int(screenWidth-foodSize))),
		Y:         float32(rand.Intn(screenHeight - foodSize)),
		Eaten:     false,
		Texture:   foodTexture,
		SpawnTime: time.Now(),
		Energy:    rand.Intn(MaxFoodEnergy-MinFoodEnergy+1) + MinFoodEnergy,
	}
	w.Foods = append(w.Foods, &newFood)
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
	w.DrawPet()
	w.DrawFoods()
	//w.DrawDirts()
	DrawStats()
	rl.EndDrawing()
}

func (w *World) Update() {
	select {
	case <-foodTicker.C:
		world.SpawnFood()
	case <-animationTicker.C:
		world.Pet.Animate()
	case <-timeTicker.C:
		world.Pet.GetOlder()
	default:
	}
	world.Pet.MoveToFood()
}
