package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

// Boss ボスを表す構造体
type Boss struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	direction              int // 移動の方向
	bulletPicChangeCounter int // 画像変更のカウンター
	nowDisplaying          int // 何番目の画像を表示しているか
}

// BossSystem ボスシステム
type BossSystem struct {
	world        *ecs.World
	bulletEntity *Bullet
	texture      *common.Texture
}

// New 新しく作成する
func (bs *BossSystem) New(w *ecs.World) {
	bs.Init(w)
}

// Init 初期化
func (bs *BossSystem) Init(w *ecs.World) {
}

// Remove 削除する
func (bs *BossSystem) Remove(entity ecs.BasicEntity) {
}

// Update アップデートする
func (bs *BossSystem) Update(dt float32) {
}
