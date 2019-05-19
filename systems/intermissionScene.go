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

type intermissonScene struct{}

func (*intermissonScene) Type() string { return "intermissonScene" }

func (*intermissonScene) Preload() {
	engo.Files.Load(
		"pics/overworld_tileset_grass.png",
	)
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	common.SetBackground(color.RGBA{255, 250, 220, 0})
}

func (*intermissonScene) Setup(u engo.Updater) {
	engo.Input.RegisterButton("Space", engo.KeySpace)
	world, _ := u.(*ecs.World)
	world.AddSystem(&IntermissonSystem{})
}

func (*intermissonScene) Exit() {
	engo.Exit()
}
