package systems

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Enemy struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	// 移動の状態(0 => 着地中, 1 => 上に移動中, 2 => 右に移動中, 3 => 下に移動中, 4 => 左に移動中)
	movingState int
	// 移動の残り時間
	movingDuration int
}

type EnemySystem struct {
	world       *ecs.World
	enemyEntity []*Enemy
	texture     *common.Texture
}

func (es *EnemySystem) Remove(entity ecs.BasicEntity) {
	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

func (es *EnemySystem) Update(dt float32) {
	rand.Seed(time.Now().UnixNano())
	// カメラとプレーヤーの位置を取得
	camX := camEntity.X()
	camY := camEntity.Y()
	playerX := playerInstance.SpaceComponent.Position.X
	playerY := playerInstance.SpaceComponent.Position.Y
	for _, o := range es.enemyEntity {
		// 画面に描画されていないオブジェクトは移動処理をしない
		if o.SpaceComponent.Position.X < camX+300 && o.SpaceComponent.Position.X > camX-300 && o.SpaceComponent.Position.Y < camY+300 && o.SpaceComponent.Position.Y > camY-300 {
			// プレーヤーとの当たり判定
			if o.SpaceComponent.Position.X == playerX && o.SpaceComponent.Position.Y == playerY {
				fmt.Println("damaged")
			}
			// 移動をしていない場合
			if o.movingState == 0 {
				tmpNum := rand.Intn(200) - 195
				if tmpNum > 0 {
					o.movingState = tmpNum
					jumpTemp := rand.Intn(3)
					switch jumpTemp {
					case 0:
						o.movingDuration = 10
					case 1:
						o.movingDuration = 15
					case 2:
						o.movingDuration = 20
					}
				}
			}
			if o.movingState == 1 {
				// 上への移動処理
				if o.movingDuration > 0 {
					o.SpaceComponent.Position.X -= 3
					o.movingDuration -= 1
				} else {
					o.movingState = 0
				}
			} else if o.movingState == 2 {
				//右への移動
				if o.movingDuration > 0 {
					o.SpaceComponent.Position.Y += 3
					o.movingDuration -= 1
				} else {
					o.movingState = 0
				}
			} else if o.movingState == 3 {
				//下への移動
				if o.movingDuration > 0 {
					o.SpaceComponent.Position.X += 3
					o.movingDuration -= 1
				} else {
					o.movingState = 0
				}
			} else if o.movingState == 4 {
				//下への移動
				if o.movingDuration > 0 {
					o.SpaceComponent.Position.Y -= 3
					o.movingDuration -= 1
				} else {
					o.movingState = 0
				}
			}
		}
	}
}

func (es *EnemySystem) New(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())
	es.world = w
	Enemies := make([]*Enemy, 0)
	// ランダムで配置
	for i := 0; i < 4000; i++ {
		randomNum := rand.Intn(400)
		if randomNum == 0 {
			// 敵の作成
			enemy := Enemy{BasicEntity: ecs.NewBasic()}
			enemy.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X: float32(i), Y: float32(rand.Intn(300))},
				Width:    30,
				Height:   30,
			}
			// 画像の読み込み
			texture, err := common.LoadedSprite("pics/ghost.png")
			if err != nil {
				fmt.Println("Unable to load texture: " + err.Error())
			}
			enemy.RenderComponent = common.RenderComponent{
				Drawable: texture,
				Scale:    engo.Point{X: 1.1, Y: 1.1},
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
		es.enemyEntity = Enemies
	}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			camEntity = sys
		}
	}
}
