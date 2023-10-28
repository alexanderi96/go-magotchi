package engine

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/alexanderi96/go-magotchi/config"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Texture struct {
	Data rl.Texture2D
}

type TextureSet struct {
	// https://lucapixel.itch.io/free-food-pixel-art-45-icons
	FoodTexture *Texture

	//https://elthen.itch.io/2d-pixel-art-fox-sprites
	IdleTextures     []*Texture
	MovingTextures   []*Texture
	SleepingTextures []*Texture

	// https://www.flaticon.com/free-icon/pharmacy_2695914
	// https://www.flaticon.com/free-icon/hunger_4968451
	// https://www.flaticon.com/free-icon/happy_1023758
	// https://www.flaticon.com/free-icon/flash_2511629
	// https://www.flaticon.com/free-icon/time_3240587
	HealthIcon    Texture
	HungerIcon    Texture
	HappinessIcon Texture
	EnergyIcon    Texture
	AgeIcon       Texture
}

type World struct {
	Foods    []*Food
	Dirts    []*Dirt
	Pet      *Pet
	Config   *config.Config
	Paused   bool
	Textures *TextureSet

	foodTicker      *time.Ticker
	animationTicker *time.Ticker
	timeTicker      *time.Ticker
	movementTicker  *time.Ticker
}

var (
	//go:embed assets/*
	assets embed.FS
)

func NewGame() (*World, error) {

	config, err := config.ReadConfig("./config.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	world := &World{
		Foods: []*Food{},
		Dirts: []*Dirt{},
		Pet: &Pet{
			X:         float32(config.ViewportX / 2),
			Y:         float32(config.ViewportY / 2),
			Health:    100,
			Hunger:    50,
			Happiness: 50,
			Energy:    100,
			Age:       0,
			FrameIdx:  0,
			Sleeping:  false,
			Dead:      false,
		},
		Config: config,
		Paused: false,

		Textures: &TextureSet{},

		foodTicker:      time.NewTicker(10 * time.Second),
		animationTicker: time.NewTicker(150 * time.Millisecond),
		timeTicker:      time.NewTicker(time.Second),
		movementTicker:  time.NewTicker(time.Millisecond),
	}

	world.loadTextures()

	return world, nil
}

func (t *Texture) loadAnimation(path string) error {
	textureData, err := assets.ReadFile(path)
	log.Println("loading animation: ", path)
	if err != nil {
		return fmt.Errorf("error loading animation %s: %x", path, err)
	}
	log.Println("animation loaded")
	t.Data = rl.LoadTextureFromImage(rl.LoadImageFromMemory(".png", textureData, int32(len(textureData))))
	log.Println("returning nil error")
	return nil
}

func loadAnimations(t []*Texture, path string) error {
	log.Println("loading animations: ", path)
	framePath, err := assets.ReadDir(path)
	if err != nil {
		return fmt.Errorf("invalid path: %x", err)
	}
	for i := 1; i <= len(framePath); i++ {
		log.Println("loading frame: ", i)
		texture := &Texture{}

		animationFrame := fmt.Sprintf(path+"/%d.png", i)
		if err := texture.loadAnimation(animationFrame); err != nil {
			return fmt.Errorf("error loading animations: %x", err)
		}
		log.Println("appending frame to the list")
		t = append(t, texture)
	}
	return nil
}

func (w *World) loadTextures() {

	loadAnimations(w.Textures.IdleTextures, "assets/pet/fox/still")
	loadAnimations(w.Textures.MovingTextures, "assets/pet/fox/walking")

	// for i := 1; i <= 5; i++ {
	// 	texture := rl.LoadTexture(fmt.Sprintf("assets/pet/fox/sleeping/%d.png", i))
	// 	w.Textures.SleepingTextures = append(w.Textures.SleepingTextures, texture)
	// }

	w.Textures.FoodTexture.loadAnimation("assets/food/food.png")
	w.Textures.HealthIcon.loadAnimation("assets/hud/health.png")
	w.Textures.HungerIcon.loadAnimation("assets/hud/hunger.png")
	w.Textures.HappinessIcon.loadAnimation("assets/hud/happiness.png")
	w.Textures.EnergyIcon.loadAnimation("assets/hud/energy.png")
	w.Textures.AgeIcon.loadAnimation("assets/hud/age.png")
}

func (w *World) spawnFood() {
	X := float32(rand.Intn(int(w.Config.ViewportX)))
	Y := float32(rand.Intn(int(w.Config.ViewportY)))

	for _, food := range w.Foods {
		if food.X == X && food.Y == Y {
			return
		}
	}

	w.Foods = append(w.Foods, NewFood(X, Y, w.Textures.FoodTexture))
}

func (w *World) spawnDirt() {
	// TODO: dirt spawn
}

func (w *World) Update() {
	select {
	case <-w.foodTicker.C:
		if !w.Pet.Dead {
			w.spawnFood()
		}

	case <-w.animationTicker.C:
		w.Pet.Animate(w.Textures)

	case <-w.timeTicker.C:
		if !w.Pet.Dead {
			w.Pet.GetOlder()
		}

	case <-w.movementTicker.C:
		if w.Pet.WantToMove() {
			w.Pet.MoveToFood(float32(w.Config.ViewportX), float32(w.Config.ViewportY), w.Foods)
		}
	default:
	}
	w.Pet.Update()
}
