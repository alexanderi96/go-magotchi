package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Food struct {
	Texture rl.Texture2D
	X, Y    float32
	Eaten   bool
}

type Dirt struct {
	Texture rl.Texture2D
	X, Y    float32
	Eaten   bool
}
