package systems

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
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

// intermissonState シーンの切り替え中かどうか
var intermissionState bool

// 次のステージの情報
var nextStage *PortalStruct

// IntermissionSystem intermisson
type IntermissionSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
}

// New 新規作成
func (is *IntermissionSystem) New(w *ecs.World) {
	is.world = w
	// 画面を黒く覆う
	shadePic, _ = common.LoadedSprite("pics/black_bk.png")
	shadingProgress = 0
	intermissionState = false
}

// Remove 削除する
func (is *IntermissionSystem) Remove(entity ecs.BasicEntity) {}

// Update アップデートする
func (is *IntermissionSystem) Update(dt float32) {
	if !intermissionState {
		return
	}
	if shadingProgress < 25 {
		camX := camEntity.X()
		camY := camEntity.Y()
		Shades := make([]*Shade, 0)
		for j := 0; j < 38; j++ {
			shade := &Shade{BasicEntity: ecs.NewBasic()}
			// 描画位置の指定
			shade.SpaceComponent.Position = engo.Point{
				X: float32(j*cellLength + int(camX)/2),
				Y: float32(shadingProgress*cellLength + int(camY)/2),
			}
			// 見た目の設定
			shade.RenderComponent = common.RenderComponent{
				Drawable: shadePic,
				Scale:    engo.Point{X: float32(cellLength / 16), Y: float32(cellLength / 16)},
			}
			shade.RenderComponent.SetZIndex(3)
			Shades = append(Shades, shade)
		}
		// シェードの追加
		for _, system := range is.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				for _, s := range Shades {
					sys.Add(&s.BasicEntity, &s.RenderComponent, &s.SpaceComponent)
				}
			}
		}
		time.Sleep(20 * time.Millisecond)
		shadingProgress++
	} else if shadingProgress == 25 {
		// 各種システムの切り替え処理
		stageFileToRead = "./stages/" + nextStage.file + ".json"
		for _, system := range is.world.Systems() {
			switch sys := system.(type) {
			case *SceneSystem:
				for _, tile := range tileEntities {
					sys.Remove(tile.BasicEntity)
				}
				sys.Init(is.world)
			case *PlayerSystem:
				sys.Remove(sys.playerEntity.BasicEntity)
				sys.Init(is.world)
			case *EnemySystem:
				for _, enemy := range enemyEntities {
					sys.Remove(enemy.BasicEntity)
				}
				sys.Init(is.world)
			}
		}
		shadingProgress++
	} else {

	}
}
