package systems

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/KMimura/RPGGame/utils"
)

// 被弾した際に爆発し続ける時間
var explosionTime = 30

// 敵の画像の大きさ
var enemyRadius float32 = 7

// Enemy 敵
type Enemy struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	// 移動の状態(0 => 着地中, 1 => 上に移動中, 2 => 右に移動中, 3 => 下に移動中, 4 => 左に移動中)
	movingState int
	// 爆発し始めてからの経過時間
	explosionDuration int
	// セルのX座標
	cellX int
	// セルのY座標
	cellY int
	// 移動の目標地点の座標
	destinationPoint float32
	// 移動の速度
	velocity float32
}

// EnemySystem 敵のシステム
type EnemySystem struct {
	world   *ecs.World
	texture *common.Texture
}

var enemyEntities []*Enemy

// 被弾した時の画像
var explosion *common.Texture

func (es *EnemySystem) Remove(entity ecs.BasicEntity) {
	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

func (es *EnemySystem) Update(dt float32) {
	// カメラとプレーヤーの位置を取得
	camX := camEntity.X()
	camY := camEntity.Y()
	// プレーヤーの座標（画像の大きさを加味）
	playerX := int(playerInstance.SpaceComponent.Position.X+playerRadius) / utils.AbstractionValue
	playerY := int(playerInstance.SpaceComponent.Position.Y+playerRadius) / utils.AbstractionValue
	for _, o := range enemyEntities {
		if o.explosionDuration != 0 {
			if o.explosionDuration == 1 {
				// 爆発し始め
				o.RenderComponent.Drawable = explosion
			}
			// 爆発が終わった場合、敵を削除
			if o.explosionDuration >= explosionTime {
				es.world.RemoveEntity(o.BasicEntity)
				enemyEntities = removeEnemy(enemyEntities, o)
			}
			o.explosionDuration++
		} else {
			// 画面に描画されていないオブジェクトは移動処理をしない
			if o.SpaceComponent.Position.X < camX+300 && o.SpaceComponent.Position.X > camX-300 && o.SpaceComponent.Position.Y < camY+300 && o.SpaceComponent.Position.Y > camY-300 {
				// プレーヤーとの当たり判定(画像の大きさを加味)
				if int(o.SpaceComponent.Position.X+enemyRadius)/utils.AbstractionValue == playerX && int(o.SpaceComponent.Position.Y+enemyRadius)/utils.AbstractionValue == playerY {
					AfflictDamage(es.world)
				}
				// 移動をしていない場合
				if o.movingState == 0 {
					tmpNum := rand.Intn(200) - 195
					if tmpNum > 0 {
						o.movingState = tmpNum
						switch tmpNum {
						case 1:
							o.destinationPoint = o.SpaceComponent.Position.Y - float32(cellLength)
						case 2:
							o.destinationPoint = o.SpaceComponent.Position.X + float32(cellLength)
						case 3:
							o.destinationPoint = o.SpaceComponent.Position.Y + float32(cellLength)
						case 4:
							o.destinationPoint = o.SpaceComponent.Position.X - float32(cellLength)
						}
					}
				}
				if o.movingState == 1 {
					// 上への移動処理
					if o.SpaceComponent.Position.Y-o.velocity > o.destinationPoint {
						o.SpaceComponent.Position.Y -= o.velocity
					} else {
						o.SpaceComponent.Position.Y = o.destinationPoint
						o.movingState = 0
						o.cellY--
					}
				} else if o.movingState == 2 {
					//右への移動
					if o.SpaceComponent.Position.X+o.velocity < o.destinationPoint {
						o.SpaceComponent.Position.X += o.velocity
					} else {
						o.SpaceComponent.Position.X = o.destinationPoint
						o.movingState = 0
						o.cellX++
					}
				} else if o.movingState == 3 {
					//下への移動
					if o.SpaceComponent.Position.Y+o.velocity < o.destinationPoint {
						o.SpaceComponent.Position.Y += o.velocity
					} else {
						o.SpaceComponent.Position.Y = o.destinationPoint
						o.movingState = 0
						o.cellY++
					}
				} else if o.movingState == 4 {
					//下への移動
					if o.SpaceComponent.Position.X-o.velocity > o.destinationPoint {
						o.SpaceComponent.Position.X -= o.velocity
					} else {
						o.SpaceComponent.Position.X = o.destinationPoint
						o.movingState = 0
						o.cellX--
					}
				}
			}
		}
	}
}

func (es *EnemySystem) New(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())
	es.world = w
	Enemies := make([]*Enemy, 0)

	// 画像の読み込み
	texture, err := common.LoadedSprite("pics/ghost.png")
	if err != nil {
		fmt.Println("Unable to load texture: " + err.Error())
	}
	// 被弾した時の画像
	explosion, _ = common.LoadedSprite("pics/explosion.png")

	// ランダムで配置
	for i := 0; i < 40; i++ {
		randomNum := rand.Intn(1)
		if randomNum == 0 {
			// 敵の作成
			enemy := Enemy{BasicEntity: ecs.NewBasic()}
			enemy.cellX = i * cellLength
			enemy.cellY = rand.Intn(30) * cellLength
			enemy.velocity = 3
			enemy.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X: float32(enemy.cellX), Y: float32(enemy.cellY)},
				Width:    30,
				Height:   30,
			}
			enemy.RenderComponent = common.RenderComponent{
				Drawable: texture,
				Scale:    engo.Point{X: 1, Y: 1},
			}
			enemy.RenderComponent.SetZIndex(1)
			es.texture = texture
			for _, system := range es.world.Systems() {
				switch sys := system.(type) {
				case *common.RenderSystem:
					sys.Add(&enemy.BasicEntity, &enemy.RenderComponent, &enemy.SpaceComponent)
				}
			}
			Enemies = append(Enemies, &enemy)
		}
		enemyEntities = Enemies
	}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			camEntity = sys
		}
	}
}

func removeEnemy(enemies []*Enemy, search *Enemy) []*Enemy {
	result := []*Enemy{}
	for _, v := range enemies {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
