package systems

import (
	"fmt"
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

// shadeXCells X軸方向にいくつシェードの画像を描画するかを規定
const shadeXCells = 38

// shadeYCells Y軸方向にいくつシェードの画像を描画するかを規定
const shadeYCells = 25

// shadingProgressY Y座標方向に、どれだけシェードの描画が進んだか
var shadingProgress int

// shadePic シェードの画像
var shadePic *common.Texture

// intermissonState シーンの切り替え中かどうか
var intermissionState bool

// 次のステージの情報
var nextStage *PortalStruct

// シェード画像の配列
var shadesArray [][]*Shade

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
func (is *IntermissionSystem) Remove(entity ecs.BasicEntity) {
	for _, system := range is.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// Update アップデートする
func (is *IntermissionSystem) Update(dt float32) {
	if !intermissionState {
		return
	}
	if shadingProgress < shadeYCells {
		camX := camEntity.X()
		camY := camEntity.Y()
		Shades := make([]*Shade, 0)
		for j := 0; j < shadeXCells; j++ {
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
		shadesArray = append(shadesArray, Shades)
		time.Sleep(20 * time.Millisecond)
		shadingProgress++
	} else if shadingProgress == shadeYCells {
		// 各種システムのデータ削除・初期化処理
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
		fmt.Println(float32(cameraInitialPositionX))
		fmt.Println(float32(cameraInitialPositionY))
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:        common.XAxis,
			Value:       float32(cameraInitialPositionX),
			Incremental: false,
		})
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:        common.YAxis,
			Value:       float32(cameraInitialPositionY),
			Incremental: false,
		})
		shadingProgress++
	} else {
		// シェードの削除
		for _, shades := range shadesArray {
			for _, shade := range shades {
				is.Remove(shade.BasicEntity)
			}
			time.Sleep(20 * time.Millisecond)
		}
		intermissionState = false
	}
}
