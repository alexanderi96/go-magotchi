package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/alexanderi96/go-magotchi/engine"
	"github.com/alexanderi96/go-magotchi/gui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	game *engine.World
)

func init() {

	game, err := engine.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if game.Config.IsResizable {
		rl.SetConfigFlags(rl.FlagWindowResizable)
	}

	rl.SetTargetFPS(game.Config.TargetFPS)
	rl.InitWindow(game.Config.WindowWidth, game.Config.WindowHeight, "Go-Magotchi")

}

func main() {

	if game.Config.ShouldBeProfiled {
		f, err := os.Create("cpu.pprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeyP) {
			game.Paused = !game.Paused
		}

		if !game.Paused {
			game.Update()
		}

		gui.Draw(game)
	}

	performCloseTasks()

	rl.CloseWindow()
}

func performCloseTasks() {
	gui.UnloadTextures(game)
}
