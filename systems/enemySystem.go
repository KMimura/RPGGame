package systems

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
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
	movingState       int     // 移動の状態(0 => 着地中, 1 => 上に移動中, 2 => 右に移動中, 3 => 下に移動中, 4 => 左に移動中)
	explosionDuration int     // 爆発し始めてからの経過時間
	cellX             int     // セルのX座標
	cellY             int     // セルのY座標
	destinationPoint  float32 // 移動の目標地点の座標
	velocity          float32 // 移動の速度
}

// EnemySystem 敵のシステム
type EnemySystem struct {
	world   *ecs.World
	texture *common.Texture
}

// 敵のエンティティの配列
var enemyEntities []*Enemy

// 被弾した時の画像
var explosion *common.Texture

// Remove 削除する
func (es *EnemySystem) Remove(entity ecs.BasicEntity) {
	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// Update 毎フレームごとに呼び出される
func (es *EnemySystem) Update(dt float32) {
	// カメラとプレーヤーの位置を取得
	camX := camEntity.X()
	camY := camEntity.Y()
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
				if o.cellX == playerInstance.cellX && o.cellY == playerInstance.cellY {
					afflictDamage(es.world)
				}
				// 移動をしていない場合
				if o.movingState == 0 {
					tmpNum := rand.Intn(200) - 195
					if tmpNum > 0 {
						o.movingState = tmpNum
						switch tmpNum {
						case 1:
							if CheckIfPassable(o.cellX, o.cellY-1) {
								o.destinationPoint = o.SpaceComponent.Position.Y - float32(cellLength)
							}
						case 2:
							if CheckIfPassable(o.cellX+1, o.cellY) {
								o.destinationPoint = o.SpaceComponent.Position.X + float32(cellLength)
							}
						case 3:
							if CheckIfPassable(o.cellX, o.cellY+1) {
								o.destinationPoint = o.SpaceComponent.Position.Y + float32(cellLength)
							}
						case 4:
							if CheckIfPassable(o.cellX-1, o.cellY) {
								o.destinationPoint = o.SpaceComponent.Position.X - float32(cellLength)
							}
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
					//左への移動
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

// New 新規作成時に呼び出される
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
			enemy.cellX = i
			enemy.cellY = rand.Intn(30)
			enemy.velocity = 3
			enemy.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X: float32(enemy.cellX * cellLength), Y: float32(enemy.cellY * cellLength)},
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

// removeEnemy 敵のエンティティを削除する
func removeEnemy(enemies []*Enemy, search *Enemy) []*Enemy {
	result := []*Enemy{}
	for _, v := range enemies {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
