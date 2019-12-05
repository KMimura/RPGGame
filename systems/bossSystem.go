package systems

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Boss ボスを表す構造体
type Boss struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	movingState            int     // 移動の状態(0 => 着地中, 1 => 上に移動中, 2 => 右に移動中, 3 => 下に移動中, 4 => 左に移動中)
	direction              int     // 移動の方向
	bulletPicChangeCounter int     // 画像変更のカウンター
	nowDisplaying          int     // 何番目の画像を表示しているか
	life                   int     // ライフ
	cellX                  [2]int  // X座標
	cellY                  [2]int  // Y座標
	destinationPoint       float32 // 移動の目標地点の座標
	velocity               float32 // 移動の速度
}

// BossBar ボスのライフバー
type BossBar struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// BossBullet ボスの出す弾
type BossBullet struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	direction              int // 弾の進む方向（0から7まで、時計回り）
	bulletPicChangeCounter int // 画像変更のカウンター
	nowDisplaying          int // 何番目の画像を表示しているか
}

// BossSystem ボスシステム
type BossSystem struct {
	world        *ecs.World
	bulletEntity *Bullet
	texture      *common.Texture
}

var bossInstance *Boss
var bossBarInstance *BossBar
var bossBulletEntities []*BossBullet // ボスの弾の配列

// ライフバーの画像の配列
var bars []*common.Texture

// New 初期化
func (bs *BossSystem) New(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())
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
	boss.velocity = 6
	boss.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(boss.cellX[0] * cellLength), Y: float32(boss.cellY[0] * cellLength)},
		Width:    64,
		Height:   64,
	}
	boss.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 3, Y: 3},
	}
	boss.RenderComponent.SetZIndex(1)
	bossInstance = &boss
	bs.texture = texture

	bossBar := BossBar{BasicEntity: ecs.NewBasic()}
	bossBar.SpaceComponent = common.SpaceComponent{Position: engo.Point{X: 180, Y: 520}, Width: 453, Height: 24}
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
		// たまに弾を出す
		tmpNum := rand.Intn(150)
		if tmpNum == 50 {
			for i := 0; i < 8; i++ {
				bs.addBossBullet(i)
			}
		}

		// 移動
		if bossInstance.movingState == 0 {
			// 現在移動中でない場合
			tmpNum -= 145
			if tmpNum > 0 {
				switch tmpNum {
				case 1:
					if checkIfPassable(bossInstance.cellX[0], bossInstance.cellY[0]-2) {
						bossInstance.movingState = tmpNum
						bossInstance.destinationPoint = bossInstance.SpaceComponent.Position.Y - float32(cellLength)*2
					}
				case 2:
					if checkIfPassable(bossInstance.cellX[0]+2, bossInstance.cellY[0]) {
						bossInstance.movingState = tmpNum
						bossInstance.destinationPoint = bossInstance.SpaceComponent.Position.X + float32(cellLength)*2
					}
				case 3:
					if checkIfPassable(bossInstance.cellX[0], bossInstance.cellY[0]+2) {
						bossInstance.movingState = tmpNum
						bossInstance.destinationPoint = bossInstance.SpaceComponent.Position.Y + float32(cellLength)*2
					}
				case 4:
					if checkIfPassable(bossInstance.cellX[0]-2, bossInstance.cellY[0]) {
						bossInstance.movingState = tmpNum
						bossInstance.destinationPoint = bossInstance.SpaceComponent.Position.X - float32(cellLength)*2
					}
				}
			}
		} else {
			if bossInstance.movingState == 1 {
				// 上への移動処理
				if bossInstance.SpaceComponent.Position.Y-bossInstance.velocity > bossInstance.destinationPoint {
					bossInstance.SpaceComponent.Position.Y -= bossInstance.velocity
				} else {
					bossInstance.SpaceComponent.Position.Y = bossInstance.destinationPoint
					bossInstance.movingState = 0
					bossInstance.cellY[0] = bossInstance.cellY[0] - 2
					bossInstance.cellY[1] = bossInstance.cellY[1] - 2
				}
			} else if bossInstance.movingState == 2 {
				//右への移動
				if bossInstance.SpaceComponent.Position.X+bossInstance.velocity < bossInstance.destinationPoint {
					bossInstance.SpaceComponent.Position.X += bossInstance.velocity
				} else {
					bossInstance.SpaceComponent.Position.X = bossInstance.destinationPoint
					bossInstance.movingState = 0
					bossInstance.cellX[0] = bossInstance.cellX[0] + 2
					bossInstance.cellX[1] = bossInstance.cellX[1] + 2
				}
			} else if bossInstance.movingState == 3 {
				//下への移動
				if bossInstance.SpaceComponent.Position.Y+bossInstance.velocity < bossInstance.destinationPoint {
					bossInstance.SpaceComponent.Position.Y += bossInstance.velocity
				} else {
					bossInstance.SpaceComponent.Position.Y = bossInstance.destinationPoint
					bossInstance.movingState = 0
					bossInstance.cellY[0] = bossInstance.cellY[0] + 2
					bossInstance.cellY[1] = bossInstance.cellY[1] + 2
				}
			} else if bossInstance.movingState == 4 {
				//左への移動
				if bossInstance.SpaceComponent.Position.X-bossInstance.velocity > bossInstance.destinationPoint {
					bossInstance.SpaceComponent.Position.X -= bossInstance.velocity
				} else {
					bossInstance.SpaceComponent.Position.X = bossInstance.destinationPoint
					bossInstance.movingState = 0
					bossInstance.cellX[0] = bossInstance.cellX[0] - 2
					bossInstance.cellX[1] = bossInstance.cellX[1] - 2
				}
			}
		}
	}
	// 弾の移動
	for _, bullet := range bossBulletEntities {
		bullet.bulletPicChangeCounter++
		bulletPicIndex := bullet.bulletPicChangeCounter / 5
		if bulletPicIndex > 7 {
			bs.Remove(bullet.BasicEntity)
			bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			continue
		}
		bullet.RenderComponent.Drawable = bulletPics[bulletPicIndex]
		switch bullet.direction {
		case 0:
			if checkIfPassable(int(bullet.SpaceComponent.Position.X)/cellLength, (int(bullet.SpaceComponent.Position.Y)-5)/cellLength) && bullet.SpaceComponent.Position.Y >= camEntity.Y()-250 {
				bullet.SpaceComponent.Position.Y -= 5
			} else {
				bs.Remove(bullet.BasicEntity)
				bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			}
		case 1:
			if checkIfPassable((int(bullet.SpaceComponent.Position.X)+5)/cellLength, (int(bullet.SpaceComponent.Position.Y)-5)/cellLength) && bullet.SpaceComponent.Position.X <= camEntity.X()+250 {
				bullet.SpaceComponent.Position.X += 5
				bullet.SpaceComponent.Position.Y -= 5
			} else {
				bs.Remove(bullet.BasicEntity)
				bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			}
		case 2:
			if checkIfPassable((int(bullet.SpaceComponent.Position.X)+5)/cellLength, int(bullet.SpaceComponent.Position.Y)/cellLength) && bullet.SpaceComponent.Position.Y <= camEntity.Y()+250 {
				bullet.SpaceComponent.Position.X += 5
			} else {
				bs.Remove(bullet.BasicEntity)
				bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			}
		case 3:
			if checkIfPassable((int(bullet.SpaceComponent.Position.X)+5)/cellLength, (int(bullet.SpaceComponent.Position.Y)+5)/cellLength) && bullet.SpaceComponent.Position.X >= camEntity.X()-250 {
				bullet.SpaceComponent.Position.X += 5
				bullet.SpaceComponent.Position.Y += 5
			} else {
				bs.Remove(bullet.BasicEntity)
				bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			}
		case 4:
			if checkIfPassable(int(bullet.SpaceComponent.Position.X)/cellLength, (int(bullet.SpaceComponent.Position.Y)+5)/cellLength) && bullet.SpaceComponent.Position.Y >= camEntity.Y()-250 {
				bullet.SpaceComponent.Position.Y += 5
			} else {
				bs.Remove(bullet.BasicEntity)
				bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			}
		case 5:
			if checkIfPassable((int(bullet.SpaceComponent.Position.X)-5)/cellLength, (int(bullet.SpaceComponent.Position.Y)+5)/cellLength) && bullet.SpaceComponent.Position.X <= camEntity.X()+250 {
				bullet.SpaceComponent.Position.X -= 5
				bullet.SpaceComponent.Position.Y += 5
			} else {
				bs.Remove(bullet.BasicEntity)
				bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			}
		case 6:
			if checkIfPassable((int(bullet.SpaceComponent.Position.X)-5)/cellLength, int(bullet.SpaceComponent.Position.Y)/cellLength) && bullet.SpaceComponent.Position.Y <= camEntity.Y()+250 {
				bullet.SpaceComponent.Position.X -= 5
			} else {
				bs.Remove(bullet.BasicEntity)
				bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			}
		case 7:
			if checkIfPassable((int(bullet.SpaceComponent.Position.X)-5)/cellLength, (int(bullet.SpaceComponent.Position.Y)-5)/cellLength) && bullet.SpaceComponent.Position.X >= camEntity.X()-250 {
				bullet.SpaceComponent.Position.X -= 5
				bullet.SpaceComponent.Position.Y -= 5
			} else {
				bs.Remove(bullet.BasicEntity)
				bossBulletEntities = removeBossBullet(bossBulletEntities, bullet)
			}
		}
		// 弾のセル座標(自身の画像の大きさを加味)
		bulletX := int(bullet.SpaceComponent.Position.X) / cellLength
		bulletY := int(bullet.SpaceComponent.Position.Y) / cellLength
		// 当たり判定
		if bulletX == playerInstance.cellX && bulletY == playerInstance.cellY {
			afflictDamage(bs.world)
		}
	}

}

// addBullet 弾を作成する
func (bs *BossSystem) addBossBullet(dir int) {
	// 弾の作成
	bullet := BossBullet{BasicEntity: ecs.NewBasic()}
	bullet.nowDisplaying = 0
	bullet.bulletPicChangeCounter = 0

	// 初期の配置
	bullet.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(bossInstance.cellX[1] * cellLength), Y: float32(bossInstance.cellY[1] * cellLength)},
		Width:    30,
		Height:   30,
	}
	bullet.RenderComponent = common.RenderComponent{
		Drawable: bulletPics[0],
		Scale:    engo.Point{X: 0.4, Y: 0.4},
	}
	bullet.RenderComponent.SetZIndex(1)
	bullet.direction = dir
	bossBulletEntities = append(bossBulletEntities, &bullet)
	// bs.texture = texture
	for _, system := range bs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&bullet.BasicEntity, &bullet.RenderComponent, &bullet.SpaceComponent)
		}
	}
}

// removeBullet 弾を削除する
func removeBossBullet(bullets []*BossBullet, search *BossBullet) []*BossBullet {
	result := []*BossBullet{}
	for _, v := range bullets {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
