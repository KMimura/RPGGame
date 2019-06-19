package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Player プレーヤーを表す構造体
type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	// 向き (0 => 移動中でない, 1 => 上, 2 => 右, 3 => 下 4 => 左)
	direction int
	// ライフ
	remainingHearts int
	// ダメージを受けない状態の残り時間
	immunityTime int
	// 移動の速度
	velocity float32
	// セルのX座標
	cellX int
	// セルのY座標
	cellY int
	// 移動の目標地点の座標
	destinationPoint float32
	// どの方向を向いているか (1 => 上, 2 => 右, 3 => 下 4 => 左)
	facingDirection int
}

// PlayerSystem プレーヤーシステム
type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
}

var playerInstance *Player

// 最大の弾の数
var maxBulletCount = 3

// 画像の半径
var playerRadius float32 = 12.5

// それぞれの向きのプレーヤーの画像
var topPic *common.Texture
var rightPic *common.Texture
var bottomPic *common.Texture
var leftPic *common.Texture

func (ps *PlayerSystem) New(w *ecs.World) {
	ps.world = w
	// プレーヤーの作成
	player := Player{BasicEntity: ecs.NewBasic()}

	// ライフを与える
	player.remainingHearts = 5
	// 移動はしていない
	player.direction = 0
	player.facingDirection = 1

	playerInstance = &player

	// 初期の配置
	player.cellX = 2
	player.cellY = 3
	positionX := cellLength * player.cellX
	positionY := cellLength * player.cellY
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(positionX), Y: float32(positionY)},
		Width:    30,
		Height:   30,
	}
	// 速度
	player.velocity = 5
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

	for i := 0; i < 5; i++ {
		AddHeart(w)
	}
}

// Remove 削除する
func (*PlayerSystem) Remove(ecs.BasicEntity) {}

// Update アップデートする
func (ps *PlayerSystem) Update(dt float32) {
	camX := camEntity.X()
	camY := camEntity.Y()
	switch playerInstance.direction {
	case 0:
		// 移動の処理
		if engo.Input.Button("MoveUp").Down() {
			if CheckIfPassable(playerInstance.cellX, playerInstance.cellY-1) {
				playerInstance.direction = 1
				playerInstance.facingDirection = 1
				playerInstance.destinationPoint = playerInstance.SpaceComponent.Position.Y - float32(cellLength)
			}
		} else if engo.Input.Button("MoveRight").Down() {
			if CheckIfPassable(playerInstance.cellX+1, playerInstance.cellY) {
				playerInstance.direction = 2
				playerInstance.facingDirection = 2
				playerInstance.destinationPoint = playerInstance.SpaceComponent.Position.X + float32(cellLength)
			}
		} else if engo.Input.Button("MoveDown").Down() {
			if CheckIfPassable(playerInstance.cellX, playerInstance.cellY+1) {
				playerInstance.direction = 3
				playerInstance.facingDirection = 3
				playerInstance.destinationPoint = playerInstance.SpaceComponent.Position.Y + float32(cellLength)
			}
		} else if engo.Input.Button("MoveLeft").Down() {
			if CheckIfPassable(playerInstance.cellX-1, playerInstance.cellY) {
				playerInstance.direction = 4
				playerInstance.facingDirection = 4
				playerInstance.destinationPoint = playerInstance.SpaceComponent.Position.X - float32(cellLength)
			}
		} else if engo.Input.Button("Space").JustPressed() {
			if len(bulletEntities) < maxBulletCount {
				bulletSystemInstance.addBullet(ps.playerEntity.SpaceComponent.Position.X, ps.playerEntity.SpaceComponent.Position.Y, ps.playerEntity.facingDirection)
			}
		}
	case 1:
		// 上への移動処理
		// カメラを動かす距離
		camMoveLen := playerInstance.velocity * -1
		if playerInstance.SpaceComponent.Position.Y-playerInstance.velocity > playerInstance.destinationPoint {
			// まるまるワンフレーム動き続けることができる場合
			playerInstance.SpaceComponent.Position.Y -= playerInstance.velocity
		} else if playerInstance.SpaceComponent.Position.Y-playerInstance.velocity == playerInstance.destinationPoint {
			// まるまる移動して移動が終わるとき
			playerInstance.SpaceComponent.Position.Y -= playerInstance.velocity
			playerInstance.direction = 0
			playerInstance.cellY--
		} else {
			// ワンフレームまるまるは動けない場合
			camMoveLen = playerInstance.destinationPoint - playerInstance.SpaceComponent.Position.Y
			playerInstance.SpaceComponent.Position.Y = playerInstance.destinationPoint
			playerInstance.direction = 0
			playerInstance.cellY--
		}
		// カメラの移動
		if camY-ps.playerEntity.SpaceComponent.Position.Y > 100 {
			engo.Mailbox.Dispatch(common.CameraMessage{
				Axis:        common.YAxis,
				Value:       camMoveLen,
				Incremental: true,
			})
		}
	case 2:
		// 右への移動処理
		// カメラを動かす距離
		camMoveLen := playerInstance.velocity
		if playerInstance.SpaceComponent.Position.X+playerInstance.velocity < playerInstance.destinationPoint {
			// まるまるワンフレーム動き続けることができる場合
			playerInstance.SpaceComponent.Position.X += playerInstance.velocity
		} else if playerInstance.SpaceComponent.Position.X+playerInstance.velocity == playerInstance.destinationPoint {
			// まるまる移動して移動が終わるとき
			playerInstance.SpaceComponent.Position.X += playerInstance.velocity
			playerInstance.direction = 0
			playerInstance.cellX++
		} else {
			// ワンフレームまるまるは動けない場合
			camMoveLen = playerInstance.destinationPoint - playerInstance.SpaceComponent.Position.X
			playerInstance.SpaceComponent.Position.X = playerInstance.destinationPoint
			playerInstance.direction = 0
			playerInstance.cellX++
		}
		// カメラの移動
		if camX-ps.playerEntity.SpaceComponent.Position.X < 100 {
			engo.Mailbox.Dispatch(common.CameraMessage{
				Axis:        common.XAxis,
				Value:       camMoveLen,
				Incremental: true,
			})
		}
	case 3:
		// 下への移動処理
		// カメラを動かす距離
		camMoveLen := playerInstance.velocity
		if playerInstance.SpaceComponent.Position.Y+playerInstance.velocity < playerInstance.destinationPoint {
			// まるまるワンフレーム動き続けることができる場合
			playerInstance.SpaceComponent.Position.Y += playerInstance.velocity
		} else if playerInstance.SpaceComponent.Position.Y+playerInstance.velocity == playerInstance.destinationPoint {
			// まるまる移動して移動が終わるとき
			playerInstance.SpaceComponent.Position.Y += playerInstance.velocity
			playerInstance.direction = 0
			playerInstance.cellY++
		} else {
			// ワンフレームまるまるは動けない場合
			camMoveLen = playerInstance.destinationPoint - playerInstance.SpaceComponent.Position.Y
			playerInstance.SpaceComponent.Position.Y = playerInstance.destinationPoint
			playerInstance.direction = 0
			playerInstance.cellY++
		}
		// カメラの移動
		if camY-ps.playerEntity.SpaceComponent.Position.Y < 100 {
			engo.Mailbox.Dispatch(common.CameraMessage{
				Axis:        common.YAxis,
				Value:       camMoveLen,
				Incremental: true,
			})
		}
	case 4:
		// 左への移動処理
		// カメラを動かす距離
		camMoveLen := playerInstance.velocity * -1
		if playerInstance.SpaceComponent.Position.X-playerInstance.velocity > playerInstance.destinationPoint {
			// まるまるワンフレーム動き続けることができる場合
			playerInstance.SpaceComponent.Position.X -= playerInstance.velocity
		} else if playerInstance.SpaceComponent.Position.X-playerInstance.velocity == playerInstance.destinationPoint {
			// まるまる移動して移動が終わるとき
			playerInstance.SpaceComponent.Position.X -= playerInstance.velocity
			playerInstance.direction = 0
			playerInstance.cellX--
		} else {
			// ワンフレームまるまるは動けない場合
			camMoveLen = playerInstance.destinationPoint - playerInstance.SpaceComponent.Position.X
			playerInstance.SpaceComponent.Position.X = playerInstance.destinationPoint
			playerInstance.direction = 0
			playerInstance.cellX--
		}
		// カメラの移動
		if ps.playerEntity.SpaceComponent.Position.X-camX < 100 {
			engo.Mailbox.Dispatch(common.CameraMessage{
				Axis:        common.XAxis,
				Value:       camMoveLen,
				Incremental: true,
			})
		}
	}
	// ダメージを受けない状態のカウントを減らす
	if ps.playerEntity.immunityTime > 0 {
		ps.playerEntity.immunityTime--
	}
	switch ps.playerEntity.direction {
	case 1:
		ps.playerEntity.RenderComponent.Drawable = topPic
	case 2:
		ps.playerEntity.RenderComponent.Drawable = rightPic
	case 3:
		ps.playerEntity.RenderComponent.Drawable = bottomPic
	case 4:
		ps.playerEntity.RenderComponent.Drawable = leftPic
	}
}

// Damage ライフを減らす
func AfflictDamage(w *ecs.World) {
	if playerInstance.immunityTime != 0 {
		return
	}
	playerInstance.remainingHearts--
	if playerInstance.remainingHearts < 0 {
	} else {
		RemoveHeart(w)
		playerInstance.immunityTime = 100
	}

}
