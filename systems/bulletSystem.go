package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Bullet 弾を表す構造体
type Bullet struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	direction int
}

// BulletSystem 弾システム
type BulletSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
	bulletEntity *Bullet
}

func (bs *BulletSystem) New(w *ecs.World) {
	bs.world = w
	// プレーヤーの作成
	bullet := Bullet{BasicEntity: ecs.NewBasic()}

	// 初期の配置
	bullet.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(playerInstance.SpaceComponent.Position.X), Y: float32(playerInstance.SpaceComponent.Position.Y)},
		Width:    30,
		Height:   30,
	}
	// 画像の読み込み
	texture, _ := common.LoadedSprite("pics/greenoctocat_top.png")
	bullet.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 0.1, Y: 0.1},
	}
	bullet.RenderComponent.SetZIndex(1)
	bs.bulletEntity = &bullet
	bs.bulletEntity.direction = playerInstance.direction
	bs.texture = texture
	for _, system := range bs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&bullet.BasicEntity, &bullet.RenderComponent, &bullet.SpaceComponent)
		}
	}
}

// Remove 削除する
func (*BulletSystem) Remove(ecs.BasicEntity) {}

// Update アップデートする
func (bs *BulletSystem) Update(dt float32) {
	switch bs.bulletEntity.direction {
	case 0:
		bs.bulletEntity.SpaceComponent.Position.Y -= 10
		if bs.bulletEntity.SpaceComponent.Position.Y < camEntity.Y()-50 {
			bs.Remove(bs.bulletEntity.BasicEntity)
		}
	case 1:
		bs.bulletEntity.SpaceComponent.Position.X += 10
		if bs.bulletEntity.SpaceComponent.Position.X > camEntity.X()+50 {
			bs.Remove(bs.bulletEntity.BasicEntity)
		}
	case 2:
		bs.bulletEntity.SpaceComponent.Position.Y += 10
		if bs.bulletEntity.SpaceComponent.Position.X > camEntity.Y()+50 {
			bs.Remove(bs.bulletEntity.BasicEntity)
		}
	case 3:
		bs.bulletEntity.SpaceComponent.Position.X -= 10
		if bs.bulletEntity.SpaceComponent.Position.X < camEntity.Y()-50 {
			bs.Remove(bs.bulletEntity.BasicEntity)
		}
	}
}
