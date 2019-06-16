package systems

// test
import (
	"bytes"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

type MainScene struct{}

// 読み込むcsvファイル
const (
	MAIN      = "./stages/main.csv"
	SECONDARY = "./stages/secondary.csv"
)

// 現在のステージ
var currentStage string

func (*MainScene) Type() string { return "mainScene" }

func (*MainScene) Preload() {
	engo.Files.Load("pics/greenoctocat_top.png",
		"pics/greenoctocat_left.png",
		"pics/greenoctocat_right.png",
		"pics/greenoctocat_bottom.png",
		"pics/ghost.png",
		"pics/overworld_tileset_grass.png",
		"pics/explosion.png",
		"pics/heart.png")
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	common.SetBackground(color.RGBA{255, 250, 220, 0})
}

func (*MainScene) Setup(u engo.Updater) {
	// とりあえずメインのステージを読み込む
	currentStage = MAIN
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("MoveUp", engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton("MoveDown", engo.KeyS, engo.KeyArrowDown)
	engo.Input.RegisterButton("Space", engo.KeySpace)
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&TileSystem{})
	world.AddSystem(&PlayerSystem{})
	world.AddSystem(&EnemySystem{})
	world.AddSystem(&BulletSystem{})
}

func (*MainScene) Exit() {
	engo.Exit()
}
