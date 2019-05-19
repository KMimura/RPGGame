package main

import (
	"github.com/EngoEngine/engo"
	"github.com/KMimura/RPGGame/systems"
)

func run() {
	opts := engo.RunOptions{
		Title:          "RPGGame",
		Width:          600,
		Height:         400,
		StandardInputs: true,
		NotResizable:   true,
	}
	engo.Run(opts, &systems.MainScene{})
}

func main() {
	run()
}
