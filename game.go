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
		Width:          900,
		Height:         600,
		StandardInputs: true,
		NotResizable:   true,
	}
	engo.Run(opts, &MainScene{})
}

// 現在のステージ
var currentStage string

func (*MainScene) Type() string { return "mainScene" }

func (*MainScene) Preload() {
	engo.Files.Load("pics/hone.png",
		"pics/fire.png",
		"pics/ghost.png",
		"pics/ghost_red.png",
		"pics/overworld_tileset_grass.png",
		"pics/explosion.png",
		"pics/heart.png",
		"pics/black_bk.png",
		"pics/transparent.png",
		"pics/bars/0.png",
		"pics/bars/1.png",
		"pics/bars/2.png",
		"pics/bars/3.png",
		"pics/bars/4.png",
		"pics/bars/5.png",
		"pics/bars/6.png",
		"pics/bars/7.png",
		"pics/bars/8.png",
		"pics/bars/9.png",
		"pics/bars/10.png",
		"pics/bars/11.png",
		"pics/bars/12.png",
		"pics/bars/13.png",
		"pics/bars/14.png",
		"pics/bars/15.png",
		"pics/bars/16.png",
		"pics/bars/17.png",
		"pics/bars/18.png",
		"pics/bars/19.png",
		"pics/bars/20.png",
		"pics/bars/21.png",
		"pics/bars/22.png",
		"pics/bars/23.png",
		"pics/bars/24.png",
		"pics/bars/25.png",
		"pics/bars/26.png",
		"pics/bars/27.png",
		"pics/bars/28.png",
		"pics/bars/29.png",
		"pics/bars/30.png",
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
	world.AddSystem(&systems.IntermissionSystem{})
}

func (*MainScene) Exit() {
	engo.Exit()
}

func main() {
	run()
}
