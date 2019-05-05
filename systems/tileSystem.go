package systems

import (
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var Spritesheet *common.Spritesheet

var camEntity *common.CameraSystem

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	xPoint int
	yPoint int
}

type TileSystem struct {
	world      *ecs.World
	tileEntity []*Tile
	texture    *common.Texture
}

func (*TileSystem) Remove(ecs.BasicEntity) {}

func (ts *TileSystem) Update(dt float32) {
}

func (ts *TileSystem) New(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())

	ts.world = w
	// 素材シートの読み込み
	loadTxt := "pics/overworld_tileset_grass.png"
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			// ランダムで描画するタイルの種類を変える
			randNum := rand.Intn(10)
			var tileNum int
			switch randNum {
			case 0:
				tileNum = 1
			case 1:
				tileNum = 14
			case 2:
				tileNum = 38
			default:
				tileNum = 0
			}
			// Tileエンティティの作成
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			// 描画位置の指定
			tile.SpaceComponent.Position = engo.Point{
				X: float32(i * 16),
				Y: float32(j * 16),
			}
			// 見た目の設定
			tile.RenderComponent.Drawable = Spritesheet.Cell(tileNum)
			tile.RenderComponent.SetZIndex(0)
			Tiles = append(Tiles, tile)
		}
	}
	for _, system := range ts.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				ts.tileEntity = append(ts.tileEntity, v)
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
	// カメラエンティティの取得
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			camEntity = sys
		}
	}

}
