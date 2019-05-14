package systems

import (
	// "encoding/csv"
	// "fmt"
	// "io"

	"math/rand"
	// "os"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/KMimura/RPGGame/utils"
)

// Spritesheet タイルの画像
var Spritesheet *common.Spritesheet

var camEntity *common.CameraSystem

// ObstaclePoints 障害物のある座標
var ObstaclePoints map[int][]int

// Tile タイル一つ一つを表す構造体
type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	xPoint int
	yPoint int
}

// TileSystem タイルシステム
type TileSystem struct {
	world      *ecs.World
	tileEntity []*Tile
	texture    *common.Texture
}

// tileMultiply タイルを何倍にして表示するか
var tileMultiply = 4

// Remove 削除する
func (*TileSystem) Remove(ecs.BasicEntity) {}

// Update アップデートする
func (ts *TileSystem) Update(dt float32) {
}

// New 作成時に呼び出される
func (ts *TileSystem) New(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())

	// file, err := os.Open("../assets/stages/test.csv")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// reader := csv.NewReader(file)
	// reader.Comma = ','
	// reader.LazyQuotes = true
	// for {
	// 	record, err := reader.Read()
	// 	if err == io.EOF {
	// 		break
	// 	} else if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(record)
	// }

	ts.world = w
	// 素材シートの読み込み
	loadTxt := "pics/overworld_tileset_grass.png"
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	ObstaclePoints = map[int][]int{}
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			// ランダムで描画するタイルの種類を変える
			randNum := rand.Intn(50)
			var tileNum int
			switch randNum {
			case 0:
				tileNum = 1
			case 1:
				tileNum = 95
				// 障害物として座標を記録（曖昧化のために、前後の複数点を記録）
				for x := 0; x < utils.SimpleAbstractionValue; x++ {
					for y := 0; y < utils.SimpleAbstractionValue; y++ {
						ObstaclePoints[i*16*tileMultiply+x] = append(ObstaclePoints[i*16*tileMultiply+x], j*16*tileMultiply+y)
					}
				}
			default:
				tileNum = 0
			}
			// Tileエンティティの作成
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			// 描画位置の指定
			tile.SpaceComponent.Position = engo.Point{
				X: float32(i * 16*tileMultiply),
				Y: float32(j * 16*tileMultiply),
			}
			// 見た目の設定
			tile.RenderComponent = common.RenderComponent{
				Drawable: Spritesheet.Cell(tileNum),
				Scale:    engo.Point{X: float32(tileMultiply), Y: float32(tileMultiply)},
			}

			tile.RenderComponent.SetZIndex(0)
			Tiles = append(Tiles, tile)
		}
	}

	// 障害物座標をutilsにセット（循環参照ができないため）
	utils.SetObstaclePoints(ObstaclePoints)
	for _, system := range ts.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				ts.tileEntity = append(ts.tileEntity, v)
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
	// カメラエンティティの取得
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			camEntity = sys
		}
	}

}
