package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/KMimura/RPGGame/utils"
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
	world          *ecs.World
	playerEntity   *Player
	texture        *common.Texture
	bulletEntities []*Bullet
}

var bulletSystemInstance *BulletSystem

func (bs *BulletSystem) New(w *ecs.World) {
	bs.world = w
	bulletSystemInstance = bs
}

// Remove 削除する
func (bs *BulletSystem) Remove(entity ecs.BasicEntity) {
	for _, system := range bs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// Update アップデートする
func (bs *BulletSystem) Update(dt float32) {
	for _, bullet := range bs.bulletEntities {
		switch bullet.direction {
		case 0:
			if utils.CheckIfPassable(int(bullet.SpaceComponent.Position.X), int(bullet.SpaceComponent.Position.Y)-10) && bullet.SpaceComponent.Position.Y >= camEntity.Y()-250 {
				bullet.SpaceComponent.Position.Y -= 10
			} else {
				bs.Remove(bullet.BasicEntity)
				bs.bulletEntities = removeBullet(bs.bulletEntities, bullet)
			}
		case 1:
			if utils.CheckIfPassable(int(bullet.SpaceComponent.Position.X)+10, int(bullet.SpaceComponent.Position.Y)) && bullet.SpaceComponent.Position.X <= camEntity.X()+250 {
				bullet.SpaceComponent.Position.X += 10
			} else {
				bs.Remove(bullet.BasicEntity)
				bs.bulletEntities = removeBullet(bs.bulletEntities, bullet)
			}
		case 2:
			if utils.CheckIfPassable(int(bullet.SpaceComponent.Position.X), int(bullet.SpaceComponent.Position.Y)+10) && bullet.SpaceComponent.Position.Y <= camEntity.Y()+250 {
				bullet.SpaceComponent.Position.Y += 10
			} else {
				bs.Remove(bullet.BasicEntity)
				bs.bulletEntities = removeBullet(bs.bulletEntities, bullet)
			}
		case 3:
			if utils.CheckIfPassable(int(bullet.SpaceComponent.Position.X)-10, int(bullet.SpaceComponent.Position.Y)) && bullet.SpaceComponent.Position.X >= camEntity.X()-250 {
				bullet.SpaceComponent.Position.X -= 10
			} else {
				bs.Remove(bullet.BasicEntity)
				bs.bulletEntities = removeBullet(bs.bulletEntities, bullet)
			}
		}
		// 弾の座標(曖昧化するために割る)
		bulletX := int(bullet.SpaceComponent.Position.X) / utils.AbstractionValue
		bulletY := int(bullet.SpaceComponent.Position.Y) / utils.AbstractionValue
		for _, system := range bs.world.Systems() {
			switch sys := system.(type) {
			case *EnemySystem:
				for _, e := range sys.enemyEntity {
					if bulletX == int(e.SpaceComponent.Position.X)/50 {
						if bulletY == int(e.SpaceComponent.Position.Y)/50 {
							// 爆発中でないかチェック
							if e.explosionDuration == 0 {
								e.explosionDuration = 1
								// 敵に命中した弾はワールドから削除
								bs.Remove(bullet.BasicEntity)
								bs.bulletEntities = removeBullet(bs.bulletEntities, bullet)
							}
						}
					}
				}
			}
		}
	}
}

func (bs *BulletSystem) addBullet(x, y float32, dir int) {
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
	bullet.direction = dir
	bs.bulletEntities = append(bs.bulletEntities, &bullet)
	bs.texture = texture
	for _, system := range bs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&bullet.BasicEntity, &bullet.RenderComponent, &bullet.SpaceComponent)
		}
	}
}

func removeBullet(bullets []*Bullet, search *Bullet) []*Bullet {
	result := []*Bullet{}
	for _, v := range bullets {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
