package systems

import (
	"math"
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

// 敵が追跡を始める距離
var enragedDistance float64 = 4

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
	life              int     // 敵のHP
	mode              int     // 1 => プレイヤーを追跡中, 0 => そうでない
}

// EnemySystem 敵のシステム
type EnemySystem struct {
	world   *ecs.World
	texture *common.Texture
}

// 敵のエンティティの配列
var enemyEntities []*Enemy

// 通常時の画像
var normalPic *common.Texture

// 被弾した時の画像
var explosion *common.Texture

// プレーヤーを追跡中の時の画像
var enraged *common.Texture

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
					var tmpNum int
					// 敵が追跡中であれば、移動頻度を上げる
					if o.mode == 0 {
						tmpNum = rand.Intn(250) - 245
					} else {
						tmpNum = rand.Intn(100) - 95
					}
					//　敵とプレーヤーのチェビシェフ距離
					distance := math.Abs(float64(playerInstance.cellX - o.cellX))
					if distance < math.Abs(float64(playerInstance.cellY-o.cellY)) {
						distance = math.Abs(float64(playerInstance.cellY - o.cellY))
					}
					if distance < enragedDistance {
						o.mode = 1
					} else {
						o.mode = 0
					}
					if tmpNum > 0 {
						if o.mode == 0 {
							switch tmpNum {
							case 1:
								if checkIfPassable(o.cellX, o.cellY-1) {
									o.movingState = tmpNum
									o.destinationPoint = o.SpaceComponent.Position.Y - float32(cellLength)
								}
							case 2:
								if checkIfPassable(o.cellX+1, o.cellY) {
									o.movingState = tmpNum
									o.destinationPoint = o.SpaceComponent.Position.X + float32(cellLength)
								}
							case 3:
								if checkIfPassable(o.cellX, o.cellY+1) {
									o.movingState = tmpNum
									o.destinationPoint = o.SpaceComponent.Position.Y + float32(cellLength)
								}
							case 4:
								if checkIfPassable(o.cellX-1, o.cellY) {
									o.movingState = tmpNum
									o.destinationPoint = o.SpaceComponent.Position.X - float32(cellLength)
								}
							}
						} else {
							//プレイヤーが敵より右にいるとき
							if playerInstance.cellX-o.cellX > 0 {
								//プレイヤーが敵より下にいるとき
								if playerInstance.cellY-o.cellY >= 0 {
									//右方向の座標差のほうが大きいとき
									if playerInstance.cellX-o.cellX >= playerInstance.cellY-o.cellY {
										if checkIfPassable(o.cellX+1, o.cellY) {
											o.movingState = 2
											o.destinationPoint = o.SpaceComponent.Position.X + float32(cellLength)
										}
										//下方向の座標差のほうが大きいとき
									} else {
										if checkIfPassable(o.cellX, o.cellY+1) {
											o.movingState = 3
											o.destinationPoint = o.SpaceComponent.Position.Y + float32(cellLength)
										}
									}
									//プレイヤーが敵より上にいるとき
								} else {
									//右方向の座標差のほうが大きいとき
									if playerInstance.cellX-o.cellX >= -(playerInstance.cellY - o.cellY) {

										if checkIfPassable(o.cellX+1, o.cellY) {
											o.movingState = 2
											o.destinationPoint = o.SpaceComponent.Position.X + float32(cellLength)
										}
										//上方向の座標差のほうが大きいとき
									} else {
										if checkIfPassable(o.cellX, o.cellY-1) {
											o.movingState = 1
											o.destinationPoint = o.SpaceComponent.Position.Y - float32(cellLength)
										}
									}
								}
								//プレイヤーが敵より左にいるとき
							} else {
								//プレイヤーが敵より下にいるとき
								if playerInstance.cellY-o.cellY >= 0 {
									//左方向の座標差のほうが大きいとき
									if -(playerInstance.cellX - o.cellX) >= playerInstance.cellY-o.cellY {
										if checkIfPassable(o.cellX-1, o.cellY) {
											o.movingState = 4
											o.destinationPoint = o.SpaceComponent.Position.X - float32(cellLength)
										}
										//下方向の座標差のほうが大きいとき
									} else {
										if checkIfPassable(o.cellX, o.cellY+1) {
											o.movingState = 3
											o.destinationPoint = o.SpaceComponent.Position.Y + float32(cellLength)
										}
									}
									//プレイヤーが敵より上にいるとき
								} else {
									//左方向の座標差のほうが大きいとき
									if -(playerInstance.cellX - o.cellX) >= -(playerInstance.cellY - o.cellY) {
										if checkIfPassable(o.cellX-1, o.cellY) {
											o.movingState = 4
											o.destinationPoint = o.SpaceComponent.Position.X - float32(cellLength)
										}
										//上方向の座標差のほうが大きいとき
									} else {
										if checkIfPassable(o.cellX, o.cellY-1) {
											o.movingState = 1
											o.destinationPoint = o.SpaceComponent.Position.Y - float32(cellLength)
										}
									}
								}
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
				// 画像の表示切替
				if o.mode == 0 {
					o.RenderComponent.Drawable = normalPic
				} else {
					o.RenderComponent.Drawable = enraged
				}
			}
		}
	}
}

// New 新規作成時に呼び出される
func (es *EnemySystem) New(w *ecs.World) {
	es.Init(w)
}

// Init 初期化
func (es *EnemySystem) Init(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())
	es.world = w
	Enemies := make([]*Enemy, 0)

	// 通常時の画像
	normalPic, _ = common.LoadedSprite("pics/ghost.png")

	// 被弾した時の画像
	explosion, _ = common.LoadedSprite("pics/explosion.png")

	// プレーヤーを追跡中の画像
	enraged, _ = common.LoadedSprite("pics/ghost_red.png")

	// ランダムで配置
	for _, ep := range EnemyPoints {
		// 敵の作成
		enemy := Enemy{BasicEntity: ecs.NewBasic(), cellX: ep.X, cellY: ep.Y}
		enemy.life = 30 //HPを30に設定
		enemy.velocity = 3
		enemy.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{X: float32(ep.X * cellLength), Y: float32(ep.Y * cellLength)},
			Width:    30,
			Height:   30,
		}
		enemy.RenderComponent = common.RenderComponent{
			Drawable: normalPic,
			Scale:    engo.Point{X: 1.5, Y: 1.5},
		}
		enemy.RenderComponent.SetZIndex(1)
		for _, system := range es.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&enemy.BasicEntity, &enemy.RenderComponent, &enemy.SpaceComponent)
			}
		}
		Enemies = append(Enemies, &enemy)

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
