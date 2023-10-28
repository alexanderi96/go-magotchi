package config

import (
	"github.com/spf13/viper"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Config struct {
	FullScreen       bool
	WindowWidth      int32
	WindowHeight     int32
	SidebarWidth     int32
	ViewportX        int32
	ViewportY        int32
	TargetFPS        int32
	IsResizable      bool
	ScaleFactor      float32
	ShouldBeProfiled bool
}

func ReadConfig(filepath string) (*Config, error) {
	viper.SetConfigFile(filepath)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		FullScreen:       viper.GetBool("full_screen"),
		WindowWidth:      viper.GetInt32("window_width"),
		WindowHeight:     viper.GetInt32("window_height"),
		TargetFPS:        viper.GetInt32("target_fps"),
		IsResizable:      viper.GetBool("is_resizable"),
		ScaleFactor:      float32(viper.GetFloat64("scale_factor")),
		ShouldBeProfiled: viper.GetBool("should_be_profiled"),
	}

	return config, nil
}

func (c *Config) UpdateWindowSettings() {

	currentWidth := int32(rl.GetScreenWidth())
	currentHeight := int32(rl.GetScreenHeight())

	c.WindowWidth = currentWidth
	c.WindowHeight = currentHeight

	if c.FullScreen {
		rl.ToggleFullscreen()
	}
}

func (c *Config) ResizeViewport(X, Y int32) {
	c.ViewportX = c.WindowWidth + X
	c.ViewportY = c.WindowHeight + Y
}
