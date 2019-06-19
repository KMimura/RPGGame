package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

// Shade シーン切り替え時に画面を覆うタイル
type Shade struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// shadingProgressY Y座標方向に、どれだけシェードの描画が進んだか
var shadingProgress int

// shadePic シェードの画像
var shadePic *common.Texture

// IntermissionSystem intermisson
type IntermissionSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
}

// New 新規作成
func (is *IntermissionSystem) New(w *ecs.World) {
	// is.world = w
	// camX := camEntity.X()
	// camY := camEntity.Y()
	// 画面を黒く覆う
	// shadePic, _ = common.LoadedSprite("pics/black_bk.png")
}

// Remove 削除する
func (is *IntermissionSystem) Remove(entity ecs.BasicEntity) {}

// Update アップデートする
func (is *IntermissionSystem) Update(dt float32) {
	// Shades := make([]*Shade, 0)
	// for j := 0; j < 15; j++ {
	// 	shade := &Shade{BasicEntity: ecs.NewBasic()}
	// 	// 描画位置の指定
	// 	shade.SpaceComponent.Position = engo.Point{
	// 		X: float32(j*16*tileMultiply + int(camX)),
	// 		Y: float32(i*16*tileMultiply + int(camY)),
	// 	}
	// 	// 見た目の設定
	// 	shade.RenderComponent = common.RenderComponent{
	// 		Drawable: shadePic,
	// 		Scale:    engo.Point{X: float32(tileMultiply), Y: float32(tileMultiply)},
	// 	}
	// 	shade.RenderComponent.SetZIndex(3)
	// 	Shades = append(Shades, shade)
	// }
	// // シェードの追加
	// for _, system := range is.world.Systems() {
	// 	switch sys := system.(type) {
	// 	case *common.RenderSystem:
	// 		for _, s := range Shades {
	// 			sys.Add(&s.BasicEntity, &s.RenderComponent, &s.SpaceComponent)
	// 		}
	// 	}
	// }
	// time.Sleep(200 * time.Millisecond)
	// fmt.Println("DONE")
}
