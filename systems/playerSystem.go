package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/KMimura/RPGGame/utils"
)

// Player プレーヤーを表す構造体
type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	// 向き (0 => 上, 1 => 右, 2 => 下, 3 => 左)
	direction int
	// ライフ
	remainingHearts int
	// ダメージを受けない状態の残り時間
	immunityTime int
}

// PlayerSystem プレーヤーシステム
type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
}

var playerInstance *Player
var playerSystemInstance *PlayerSystem

// 最大の弾の数
var maxBulletCount = 3

// それぞれの向きのプレーヤーの画像
var topPic *common.Texture
var rightPic *common.Texture
var bottomPic *common.Texture
var leftPic *common.Texture

func (ps *PlayerSystem) New(w *ecs.World) {
	playerSystemInstance = ps
	ps.world = w
	// プレーヤーの作成
	player := Player{BasicEntity: ecs.NewBasic()}

	// ライフを与える
	player.remainingHearts = 5

	playerInstance = &player

	// 初期の配置
	positionX := int(engo.WindowWidth() / 2)
	positionY := int(engo.WindowHeight() - 88)
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(positionX), Y: float32(positionY)},
		Width:    30,
		Height:   30,
	}
	// 画像の読み込み
	topPic, _ = common.LoadedSprite("pics/greenoctocat_top.png")
	rightPic, _ = common.LoadedSprite("pics/greenoctocat_right.png")
	bottomPic, _ = common.LoadedSprite("pics/greenoctocat_bottom.png")
	leftPic, _ = common.LoadedSprite("pics/greenoctocat_left.png")

	player.RenderComponent = common.RenderComponent{
		Drawable: topPic,
		Scale:    engo.Point{X: 0.1, Y: 0.1},
	}
	player.RenderComponent.SetZIndex(1)
	ps.playerEntity = &player
	ps.texture = topPic
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}
	common.CameraBounds = engo.AABB{
		Min: engo.Point{X: 0, Y: 0},
		Max: engo.Point{X: 1200, Y: 1200},
	}
}

// Remove 削除する
func (*PlayerSystem) Remove(ecs.BasicEntity) {}

// Update アップデートする
func (ps *PlayerSystem) Update(dt float32) {
	camX := camEntity.X()
	camY := camEntity.Y()
	// ダメージを受けてすぐであったら、新たにダメージは受けないようにする
	if ps.playerEntity.immunityTime > 0 {
		ps.playerEntity.immunityTime--
	}
	if engo.Input.Button("MoveRight").Down() {
		if ps.playerEntity.direction != 1 {
			ps.playerEntity.direction = 1
		} else {
			// 移動先のブロックに障害物がないか確認(プレーヤーのいるタイルの判別は、画像の中心部)
			if utils.CheckIfPassable(int(ps.playerEntity.SpaceComponent.Position.X+12+5)/(16*tileMultiply), int(ps.playerEntity.SpaceComponent.Position.Y+12)/(16*tileMultiply)) {
				if camX < 4000 {
					ps.playerEntity.SpaceComponent.Position.X += 5
					if ps.playerEntity.SpaceComponent.Position.X-camX > 100 {
						engo.Mailbox.Dispatch(common.CameraMessage{
							Axis:        common.XAxis,
							Value:       5,
							Incremental: true,
						})
					}
				}
			}
		}
	} else if engo.Input.Button("MoveLeft").Down() {
		if ps.playerEntity.direction != 3 {
			ps.playerEntity.direction = 3
		} else {
			// 移動先のブロックに障害物がないか確認
			if utils.CheckIfPassable(int(ps.playerEntity.SpaceComponent.Position.X+12-5)/(16*tileMultiply), int(ps.playerEntity.SpaceComponent.Position.Y+12)/(16*tileMultiply)) {
				if camX > 200 {
					ps.playerEntity.SpaceComponent.Position.X -= 5
					if camX-ps.playerEntity.SpaceComponent.Position.X > 100 {
						engo.Mailbox.Dispatch(common.CameraMessage{
							Axis:        common.XAxis,
							Value:       -5,
							Incremental: true,
						})
					}
				} else if ps.playerEntity.SpaceComponent.Position.X > 5 {
					ps.playerEntity.SpaceComponent.Position.X -= 5
				}
			}
		}
	} else if engo.Input.Button("MoveUp").Down() {
		if ps.playerEntity.direction != 0 {
			ps.playerEntity.direction = 0
		} else {
			// 移動先のブロックに障害物がないか確認
			if utils.CheckIfPassable(int(ps.playerEntity.SpaceComponent.Position.X+12)/(16*tileMultiply), int(ps.playerEntity.SpaceComponent.Position.Y+12-5)/(16*tileMultiply)) {
				if camY > 200 {
					ps.playerEntity.SpaceComponent.Position.Y -= 5
					if camY-ps.playerEntity.SpaceComponent.Position.Y > 100 {
						engo.Mailbox.Dispatch(common.CameraMessage{
							Axis:        common.YAxis,
							Value:       -5,
							Incremental: true,
						})
					}
				} else if ps.playerEntity.SpaceComponent.Position.Y > 5 {
					ps.playerEntity.SpaceComponent.Position.Y -= 5
				}
			}
		}
	} else if engo.Input.Button("MoveDown").Down() {
		if ps.playerEntity.direction != 2 {
			ps.playerEntity.direction = 2
		} else {
			// 移動先のブロックに障害物がないか確認
			if utils.CheckIfPassable(int(ps.playerEntity.SpaceComponent.Position.X+12)/(16*tileMultiply), int(ps.playerEntity.SpaceComponent.Position.Y+12+5)/(16*tileMultiply)) {
				if camY < 4000 {
					ps.playerEntity.SpaceComponent.Position.Y += 5
					if ps.playerEntity.SpaceComponent.Position.Y-camY > 100 {
						engo.Mailbox.Dispatch(common.CameraMessage{
							Axis:        common.YAxis,
							Value:       5,
							Incremental: true,
						})
					}
				}
			}
		}
	} else if engo.Input.Button("Space").JustPressed() {
		if len(bulletSystemInstance.bulletEntities) < maxBulletCount {
			bulletSystemInstance.addBullet(ps.playerEntity.SpaceComponent.Position.X, ps.playerEntity.SpaceComponent.Position.Y, ps.playerEntity.direction)
		}
	}
	switch ps.playerEntity.direction {
	case 0:
		ps.playerEntity.RenderComponent.Drawable = topPic
	case 1:
		ps.playerEntity.RenderComponent.Drawable = rightPic
	case 2:
		ps.playerEntity.RenderComponent.Drawable = bottomPic
	case 3:
		ps.playerEntity.RenderComponent.Drawable = leftPic
	}
}

// Damage ライフを減らす
func (ps *PlayerSystem) Damage() {
	if ps.playerEntity.immunityTime != 0 {
		return
	}
	ps.playerEntity.remainingHearts--
	if ps.playerEntity.remainingHearts < 0 {
	} else {
		RemoveHeart(ps.world)
		ps.playerEntity.immunityTime = 100
	}

}
