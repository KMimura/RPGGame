package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Heart 画面の定位置に表示され続けるエンティティ
type Heart struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

// HeartEntities ハートのエンティティ
var HeartEntities []*Heart

// AddHeart ライフを表すハートの画像を表示
func AddHeart(w *ecs.World) {
	// すでに作成済みのハートの数
	existingHearts := len(HeartEntities)
	hud := Heart{BasicEntity: ecs.NewBasic()}
	hud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(20 + 50*existingHearts), Y: 20},
	}
	texture, _ := common.LoadedSprite("pics/heart.png")
	hud.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 1.6, Y: 1.6},
	}
	hud.RenderComponent.SetShader(common.HUDShader)
	hud.RenderComponent.SetZIndex(1)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
		}
	}
	HeartEntities = append(HeartEntities, &hud)
}

// RemoveHeart ライフの表示を減らす
func RemoveHeart(w *ecs.World) {
	// すでに作成済みのハートの数
	existingHearts := len(HeartEntities)
	if existingHearts < 0 {
		return
	}
	heartToRemove := HeartEntities[existingHearts-1]
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(heartToRemove.BasicEntity)
		}
	}
	result := []*Heart{}
	if existingHearts > 1 {
		for i := 0; i < existingHearts-1; i++ {
			result = append(result, HeartEntities[i])
		}
	}
	HeartEntities = result
}
