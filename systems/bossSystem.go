package systems

import (
	"fmt"
	"strconv"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Boss ボスを表す構造体
type Boss struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	direction              int    // 移動の方向
	bulletPicChangeCounter int    // 画像変更のカウンター
	nowDisplaying          int    // 何番目の画像を表示しているか
	life                   int    // ライフ
	cellX                  [2]int // X座標
	cellY                  [2]int // Y座標
}

// BossBar ボスのライフバー
type BossBar struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// BossSystem ボスシステム
type BossSystem struct {
	world        *ecs.World
	bulletEntity *Bullet
	texture      *common.Texture
}

var bossInstance *Boss
var bossBarInstance *BossBar

// ライフバーの画像の配列
var bars []*common.Texture

// Init 初期化
func (bs *BossSystem) New(w *ecs.World) {
	bs.world = w
	// 画像の読み込み
	texture, err := common.LoadedSprite("pics/ghost.png")
	if err != nil {
		fmt.Println("Unable to load texture: " + err.Error())
	}
	// 被弾した時の画像
	explosion, _ = common.LoadedSprite("pics/explosion.png")

	// ライフバーの画像を配列に入れる
	for i := 0; i <= 30; i++ {
		picFile, e := common.LoadedSprite("pics/bars/" + strconv.Itoa(i) + ".png")
		if e != nil {
			fmt.Println(e)
		}
		bars = append(bars, picFile)
	}

	boss := Boss{BasicEntity: ecs.NewBasic()}
	boss.life = 300
	boss.cellX = [2]int{10, 11}
	boss.cellY = [2]int{10, 11}
	boss.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(boss.cellX[0] * cellLength), Y: float32(boss.cellY[0] * cellLength)},
		Width:    64,
		Height:   64,
	}
	boss.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 2, Y: 2},
	}
	boss.RenderComponent.SetZIndex(1)
	bossInstance = &boss
	bs.texture = texture

	bossBar := BossBar{BasicEntity: ecs.NewBasic()}
	bossBar.SpaceComponent = common.SpaceComponent{Position: engo.Point{X: 100, Y: 350}, Width: 453, Height: 24}
	bossBar.RenderComponent = common.RenderComponent{Drawable: bars[0], Scale: engo.Point{X: 1.5, Y: 1.5}}
	bossBar.RenderComponent.SetShader(common.HUDShader)
	bossBar.RenderComponent.SetZIndex(1)
	bossBarInstance = &bossBar

	for _, system := range bs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&boss.BasicEntity, &boss.RenderComponent, &boss.SpaceComponent)
			sys.Add(&bossBar.BasicEntity, &bossBar.RenderComponent, &bossBar.SpaceComponent)
		}
	}
}

// Remove 削除する
func (bs *BossSystem) Remove(entity ecs.BasicEntity) {
	for _, system := range bs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// Update アップデートする
func (bs *BossSystem) Update(dt float32) {
	if bossInstance != nil {
		// ライフバーの変更
		damageDegree := (300 - bossInstance.life) / 10
		if damageDegree <= 30 {
			bossBarInstance.RenderComponent.Drawable = bars[damageDegree]
		} else {
			bs.Remove(bossInstance.BasicEntity)
			bossInstance = nil
			bs.Remove(bossBarInstance.BasicEntity)
			bossBarInstance = nil
		}
	}
}
