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
	direction              int // 移動の方向
	bulletPicChangeCounter int // 画像変更のカウンター
	nowDisplaying          int // 何番目の画像を表示しているか
}

// BulletSystem 弾システム
type BulletSystem struct {
	world        *ecs.World
	bulletEntity *Bullet
	texture      *common.Texture
}

// 弾のエンティティの配列
var bulletEntities []*Bullet

// 弾のシステムのインスタンス
var bulletSystemInstance *BulletSystem

// 弾の画像の半径
var bulletRadius float32 = 12.5

// 弾の画像
var bulletPics []*common.Texture

// New 新しく作成する
func (bs *BulletSystem) New(w *ecs.World) {
	bs.world = w
	bulletSystemInstance = bs
	loadTxt := "pics/fire.png"
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 128, 128, 0, 0)
	for i := 0; i < 8; i++ {
		tmp := Spritesheet.Cell(i)
		bulletPics = append(bulletPics, &tmp)
	}
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
	for _, bullet := range bulletEntities {
		bullet.bulletPicChangeCounter++
		bulletPicIndex := bullet.bulletPicChangeCounter / 5
		if bulletPicIndex > 7 {
			bs.Remove(bullet.BasicEntity)
			bulletEntities = removeBullet(bulletEntities, bullet)
			continue
		}
		bullet.RenderComponent.Drawable = bulletPics[bulletPicIndex]
		switch bullet.direction {
		case 1:
			if checkIfPassable(int(bullet.SpaceComponent.Position.X)/cellLength, (int(bullet.SpaceComponent.Position.Y)-10)/cellLength) && bullet.SpaceComponent.Position.Y >= camEntity.Y()-250 {
				bullet.SpaceComponent.Position.Y -= 10
			} else {
				bs.Remove(bullet.BasicEntity)
				bulletEntities = removeBullet(bulletEntities, bullet)
			}
		case 2:
			if checkIfPassable((int(bullet.SpaceComponent.Position.X)+10)/cellLength, int(bullet.SpaceComponent.Position.Y)/cellLength) && bullet.SpaceComponent.Position.X <= camEntity.X()+250 {
				bullet.SpaceComponent.Position.X += 10
			} else {
				bs.Remove(bullet.BasicEntity)
				bulletEntities = removeBullet(bulletEntities, bullet)
			}
		case 3:
			if checkIfPassable(int(bullet.SpaceComponent.Position.X)/cellLength, (int(bullet.SpaceComponent.Position.Y)+10)/cellLength) && bullet.SpaceComponent.Position.Y <= camEntity.Y()+250 {
				bullet.SpaceComponent.Position.Y += 10
			} else {
				bs.Remove(bullet.BasicEntity)
				bulletEntities = removeBullet(bulletEntities, bullet)
			}
		case 4:
			if checkIfPassable((int(bullet.SpaceComponent.Position.X)-10)/cellLength, int(bullet.SpaceComponent.Position.Y)/cellLength) && bullet.SpaceComponent.Position.X >= camEntity.X()-250 {
				bullet.SpaceComponent.Position.X -= 10
			} else {
				bs.Remove(bullet.BasicEntity)
				bulletEntities = removeBullet(bulletEntities, bullet)
			}
		}
		// 弾のセル座標(自身の画像の大きさを加味)
		bulletX := int(bullet.SpaceComponent.Position.X) / cellLength
		bulletY := int(bullet.SpaceComponent.Position.Y) / cellLength
		// 当たり判定は、敵の画像の大きさを加味して行う
		for _, e := range enemyEntities {
			if e.cellX == bulletX && e.cellY == bulletY {
				// 爆発中でないかチェック
				if e.explosionDuration == 0 {
					e.explosionDuration = 1
					// 敵に命中した弾はワールドから削除
					bs.Remove(bullet.BasicEntity)
					bulletEntities = removeBullet(bulletEntities, bullet)
				}
			}
		}
	}
}

// addBullet 弾を作成する
func (bs *BulletSystem) addBullet(x, y float32, dir int) {
	// プレーヤーの作成
	bullet := Bullet{BasicEntity: ecs.NewBasic()}
	bullet.nowDisplaying = 0
	bullet.bulletPicChangeCounter = 0

	// 初期の配置
	bullet.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(playerInstance.SpaceComponent.Position.X), Y: float32(playerInstance.SpaceComponent.Position.Y)},
		Width:    30,
		Height:   30,
	}
	bullet.RenderComponent = common.RenderComponent{
		Drawable: bulletPics[0],
		Scale:    engo.Point{X: 0.3, Y: 0.3},
	}
	bullet.RenderComponent.SetZIndex(1)
	bullet.direction = dir
	bulletEntities = append(bulletEntities, &bullet)
	// bs.texture = texture
	for _, system := range bs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&bullet.BasicEntity, &bullet.RenderComponent, &bullet.SpaceComponent)
		}
	}
}

// removeBullet 弾を削除する
func removeBullet(bullets []*Bullet, search *Bullet) []*Bullet {
	result := []*Bullet{}
	for _, v := range bullets {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
