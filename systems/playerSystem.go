package systems

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Player プレーヤーを表す構造体
type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	direction        int     // 向き (0 => 移動中でない, 1 => 上, 2 => 右, 3 => 下 4 => 左)
	remainingHearts  int     // ライフ
	immunityTime     int     // ダメージを受けない状態の残り時間
	velocity         float32 // 移動の速度
	cellX            int     // セルのX座標
	cellY            int     // セルのY座標
	destinationPoint float32 // 移動の目標地点の座標
	facingDirection  int     // どの方向を向いているか (1 => 上, 2 => 右, 3 => 下 4 => 左)
	movingPic        bool    //移動中の画像を表示するかどうか
}

// PlayerSystem プレーヤーシステム
type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
}

// playerInstance プレーヤーのエンティティのインスタンス
var playerInstance *Player

// 最大の弾の数
var maxBulletCount = 3

// 画像の半径
var playerRadius float32 = 12.5

// それぞれの向きのプレーヤーの画像
var topPicOne *common.Texture
var topPicTwo *common.Texture
var topPicThree *common.Texture
var rightPicOne *common.Texture
var rightPicTwo *common.Texture
var rightPicThree *common.Texture
var bottomPicOne *common.Texture
var bottomPicTwo *common.Texture
var bottomPicThree *common.Texture
var leftPicOne *common.Texture
var leftPicTwo *common.Texture
var leftPicThree *common.Texture

// New 新規作成時に呼び出される
func (ps *PlayerSystem) New(w *ecs.World) {
	ps.world = w
	// プレーヤーの作成
	player := Player{BasicEntity: ecs.NewBasic()}

	// ライフを与える
	player.remainingHearts = 5
	// 移動はしていない
	player.direction = 0
	player.facingDirection = 1
	player.movingPic = false

	playerInstance = &player

	// 初期の配置
	player.cellX = playerInitialPositionX
	player.cellY = playerInitialPositionY
	positionX := cellLength * player.cellX
	positionY := cellLength * player.cellY
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(positionX), Y: float32(positionY)},
		Width:    30,
		Height:   30,
	}
	// 速度
	player.velocity = 4
	// 画像の読み込み
	loadTxt := "pics/hone.png"
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 32, 32, 0, 0)

	topPicTmpOne := Spritesheet.Cell(9)
	topPicTmpTwo := Spritesheet.Cell(11)
	topPicTmpThree := Spritesheet.Cell(10)
	rightPicTmpOne := Spritesheet.Cell(6)
	rightPicTmpTwo := Spritesheet.Cell(8)
	rightPicTmpThree := Spritesheet.Cell(7)
	bottomPicTmpOne := Spritesheet.Cell(0)
	bottomPicTmpTwo := Spritesheet.Cell(2)
	bottomPicTmpThree := Spritesheet.Cell(1)
	leftPicTmpOne := Spritesheet.Cell(3)
	leftPicTmpTwo := Spritesheet.Cell(5)
	leftPicTmpThree := Spritesheet.Cell(4)
	topPicOne = &topPicTmpOne
	topPicTwo = &topPicTmpTwo
	topPicThree = &topPicTmpThree
	rightPicOne = &rightPicTmpOne
	rightPicTwo = &rightPicTmpTwo
	rightPicThree = &rightPicTmpThree
	bottomPicOne = &bottomPicTmpOne
	bottomPicTwo = &bottomPicTmpTwo
	bottomPicThree = &bottomPicTmpThree
	leftPicOne = &leftPicTmpOne
	leftPicTwo = &leftPicTmpTwo
	leftPicThree = &leftPicTmpThree

	player.RenderComponent = common.RenderComponent{
		Drawable: topPicOne,
		Scale:    engo.Point{X: 1, Y: 1},
	}
	player.RenderComponent.SetZIndex(1)
	ps.playerEntity = &player
	ps.texture = topPicOne
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
		// 移動の開始の処理
		if engo.Input.Button("MoveUp").Down() {
			playerInstance.facingDirection = 1
			playerInstance.movingPic = !playerInstance.movingPic
			if checkIfPassable(playerInstance.cellX, playerInstance.cellY-1) {
				playerInstance.direction = 1
				playerInstance.facingDirection = 1
				playerInstance.destinationPoint = playerInstance.SpaceComponent.Position.Y - float32(cellLength)
			} else {
				ifPortal, portalInfo := checkIfPortal(playerInstance.cellX, playerInstance.cellY-1)
				if ifPortal {
					fmt.Println(portalInfo)
				}
			}
		} else if engo.Input.Button("MoveRight").Down() {
			playerInstance.facingDirection = 2
			playerInstance.movingPic = !playerInstance.movingPic
			if checkIfPassable(playerInstance.cellX+1, playerInstance.cellY) {
				playerInstance.direction = 2
				playerInstance.facingDirection = 2
				playerInstance.destinationPoint = playerInstance.SpaceComponent.Position.X + float32(cellLength)
			} else {
				ifPortal, portalInfo := checkIfPortal(playerInstance.cellX+1, playerInstance.cellY)
				if ifPortal {
					fmt.Println(portalInfo)
				}
			}
		} else if engo.Input.Button("MoveDown").Down() {
			playerInstance.facingDirection = 3
			playerInstance.movingPic = !playerInstance.movingPic
			if checkIfPassable(playerInstance.cellX, playerInstance.cellY+1) {
				playerInstance.direction = 3
				playerInstance.facingDirection = 3
				playerInstance.destinationPoint = playerInstance.SpaceComponent.Position.Y + float32(cellLength)
			} else {
				ifPortal, portalInfo := checkIfPortal(playerInstance.cellX, playerInstance.cellY+1)
				if ifPortal {
					fmt.Println(portalInfo)
				}
			}
		} else if engo.Input.Button("MoveLeft").Down() {
			playerInstance.facingDirection = 4
			playerInstance.movingPic = !playerInstance.movingPic
			if checkIfPassable(playerInstance.cellX-1, playerInstance.cellY) {
				playerInstance.direction = 4
				playerInstance.facingDirection = 4
				playerInstance.destinationPoint = playerInstance.SpaceComponent.Position.X - float32(cellLength)
			} else {
				ifPortal, portalInfo := checkIfPortal(playerInstance.cellX-1, playerInstance.cellY)
				if ifPortal {
					fmt.Println(portalInfo)
				}
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
	
	switch ps.playerEntity.facingDirection {
	case 1:
		if playerInstance.direction == 0{
			ps.playerEntity.RenderComponent.Drawable = topPicThree
		} else {
			if ps.playerEntity.movingPic {
				ps.playerEntity.RenderComponent.Drawable = topPicTwo
			} else {
				ps.playerEntity.RenderComponent.Drawable = topPicOne
			}
		}
	case 2:
		if playerInstance.direction == 0{
			ps.playerEntity.RenderComponent.Drawable = rightPicThree
		} else {
			if ps.playerEntity.movingPic {
				ps.playerEntity.RenderComponent.Drawable = rightPicTwo
			} else {
				ps.playerEntity.RenderComponent.Drawable = rightPicOne
			}
		}
	case 3:
		if playerInstance.direction == 0{
			ps.playerEntity.RenderComponent.Drawable = bottomPicThree
		} else {
			if ps.playerEntity.movingPic {
				ps.playerEntity.RenderComponent.Drawable = bottomPicTwo
			} else {
				ps.playerEntity.RenderComponent.Drawable = bottomPicOne
			}
		}
	case 4:
		if playerInstance.direction == 0 {
			ps.playerEntity.RenderComponent.Drawable = leftPicThree
		} else {
			if ps.playerEntity.movingPic {
				ps.playerEntity.RenderComponent.Drawable = leftPicTwo
			} else {
				ps.playerEntity.RenderComponent.Drawable = leftPicOne
			}
		}
	}
	// ダメージを受けない状態のカウントを減らす
	// さらに10フレームごとにプレイヤーキャラ画像点滅させる
	if ps.playerEntity.immunityTime > 0 {
		ps.playerEntity.immunityTime--
		if((ps.playerEntity.immunityTime / 10) % 2 == 0){
			//下の = の右辺の画像は仮のもの
			ps.playerEntity.RenderComponent.Drawable = topPicThree
		}
	}
}

// afflictDamage ライフを減らす
func afflictDamage(w *ecs.World) {
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
