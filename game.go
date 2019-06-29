package main

import (
	"bytes"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/KMimura/RPGGame/systems"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

type MainScene struct{}

func run() {
	opts := engo.RunOptions{
		Title:          "RPGGame",
		Width:          600,
		Height:         400,
		StandardInputs: true,
		NotResizable:   true,
	}
	engo.Run(opts, &MainScene{})
}

// 現在のステージ
var currentStage string

func (*MainScene) Type() string { return "mainScene" }

func (*MainScene) Preload() {
	engo.Files.Load("pics/characters.png",
		"pics/greenoctocat_top.png",
		"pics/ghost.png",
		"pics/overworld_tileset_grass.png",
		"pics/explosion.png",
		"pics/heart.png",
		"pics/black_bk.png",
	)
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	common.SetBackground(color.RGBA{255, 250, 220, 0})
}

func (*MainScene) Setup(u engo.Updater) {
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("MoveUp", engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton("MoveDown", engo.KeyS, engo.KeyArrowDown)
	engo.Input.RegisterButton("Space", engo.KeySpace)
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&systems.SceneSystem{})
	world.AddSystem(&systems.PlayerSystem{})
	world.AddSystem(&systems.EnemySystem{})
	world.AddSystem(&systems.BulletSystem{})
}

func (*MainScene) Exit() {
	engo.Exit()
}

func main() {
	run()
}
