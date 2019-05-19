package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Tile タイル一つ一つを表す構造体
type IntermissionTile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// IntermissonSystem インターミッションシステム
type IntermissonSystem struct {
	world   *ecs.World
	texture *common.Texture
}

// Remove 削除する
func (*IntermissonSystem) Remove(ecs.BasicEntity) {
}

// Update アップデートする
func (*IntermissonSystem) Update(dt float32) {
	if engo.Input.Button("Space").JustPressed() {
		engo.SetScene(&MainScene{}, true)
	}
}

// New 作成時に呼び出される
func (is *IntermissonSystem) New(w *ecs.World) {
	is.world = w
	// 素材シートの読み込み
	loadTxt := "pics/overworld_tileset_grass.png"
	spritesheet := common.NewSpritesheetWithBorderFromFile(loadTxt, 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	for i := 0; i < 30; i++ {
		for j := 0; j < 30; j++ {
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			// 描画位置の指定
			tile.SpaceComponent.Position = engo.Point{
				X: float32(i * 16 * tileMultiply),
				Y: float32(j * 16 * tileMultiply),
			}
			// 見た目の設定
			tile.RenderComponent = common.RenderComponent{
				Drawable: spritesheet.Cell(61),
				Scale:    engo.Point{X: float32(tileMultiply), Y: float32(tileMultiply)},
			}
			tile.RenderComponent.SetZIndex(150)
			Tiles = append(Tiles, tile)
		}
	}

	for _, system := range is.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}
