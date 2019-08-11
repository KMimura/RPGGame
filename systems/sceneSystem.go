package systems

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"math/rand"
	"os"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Spritesheet タイルの画像
var Spritesheet *common.Spritesheet

// camEntity カメラシステムのエンティティ
var camEntity *common.CameraSystem

// 敵の画像の大きさ
var tileRadius float32 = 7

// ObstaclePoints 障害物のある座標
var ObstaclePoints map[int][]int

// playerInitialPositionX,Y プレーヤーの初期位置
var playerInitialPositionX int
var playerInitialPositionY int

// cameraInitialPositionX,Y カメラの初期位置
var cameraInitialPositionX int
var cameraInitialPositionY int

// EnemyPoints 敵を出現させる座標に関する情報
var EnemyPoints []*EnemyStruct

// PortalPoints 他のステージへのポータルの情報に関する情報
var PortalPoints map[int]map[int]*PortalStruct

// cellLength セル一辺のピクセル数（必ず16の倍数にすること）
var cellLength = 32

// Tile タイル一つ一つを表す構造体
type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// SceneSystem シーンシステム
type SceneSystem struct {
	world   *ecs.World
	texture *common.Texture
}

// EnemyStruct 敵の座標情報を持つ構造体
type EnemyStruct struct {
	X  int // X座標
	Y  int // Y座標
	id int // 敵のid
}

// PortalStruct 他のステージへの入り口の情報を持つ構造体
type PortalStruct struct {
	X        int    // X座標
	Y        int    // Y座標
	position string // 初期位置
	file     string // 移動後のファイル
}

// タイルシステムのエンティティのインスタンス
var tileEntities []*Tile

// 読み込むべきステージファイル
var stageFileToRead string

// Remove 削除する
func (ss *SceneSystem) Remove(entity ecs.BasicEntity) {
	for _, system := range ss.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// Update アップデートする
func (ss *SceneSystem) Update(dt float32) {
}

// New 作成時に呼び出される
func (ss *SceneSystem) New(w *ecs.World) {
	stageFileToRead = "./stages/main.json"
	ss.Init(w)
}

// Init 初期化
func (ss *SceneSystem) Init(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())

	ss.world = w
	// 素材シートの読み込み
	loadTxt := "pics/overworld_tileset_grass.png"
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	ObstaclePoints = map[int][]int{}
	PortalPoints = make(map[int]map[int]*PortalStruct)
	file, err := os.Open(stageFileToRead)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var sceneJSON map[string]interface{}
	json.Unmarshal([]byte(byteValue), &sceneJSON)
	// プレーヤーの初期位置
	playerInitialPositionX = int(sceneJSON["meta-data"].(map[string]interface{})["player-initial-positions"].(map[string]interface{})["A"].(map[string]interface{})["X"].(float64))
	playerInitialPositionY = int(sceneJSON["meta-data"].(map[string]interface{})["player-initial-positions"].(map[string]interface{})["A"].(map[string]interface{})["Y"].(float64))
	cameraInitialPositionX = int(sceneJSON["meta-data"].(map[string]interface{})["camera-initial-positions"].(map[string]interface{})["A"].(map[string]interface{})["X"].(float64))
	cameraInitialPositionY = int(sceneJSON["meta-data"].(map[string]interface{})["camera-initial-positions"].(map[string]interface{})["A"].(map[string]interface{})["Y"].(float64))
	i := 0
	for _, r := range sceneJSON["cell-data"].([]interface{}) {
		j := 0
		for _, c := range r.([]interface{}) {
			tileNum := c.(map[string]interface{})["cell"].(float64)
			if c.(map[string]interface{})["obstacle"].(bool) == true {
				// 障害物としてセル座標を記録
				ObstaclePoints[j] = append(ObstaclePoints[j], i)
			}
			// Tileエンティティの作成
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			// 描画位置の指定
			tile.SpaceComponent.Position = engo.Point{
				X: float32(j * cellLength),
				Y: float32(i * cellLength),
			}
			// 見た目の設定
			tile.RenderComponent = common.RenderComponent{
				Drawable: Spritesheet.Cell(int(tileNum)),
				Scale:    engo.Point{X: float32(cellLength / 16), Y: float32(cellLength / 16)}, // cellLengthが画像の元の大きさ（16ピクセル）の何倍であるかを算出し、設定
			}
			tile.RenderComponent.SetZIndex(0)
			Tiles = append(Tiles, tile)
			// 敵を出現させるべきか判定
			if c.(map[string]interface{})["enemy"].(bool) == true {
				enemyStruct := EnemyStruct{X: j, Y: i, id: int(c.(map[string]interface{})["enemy-data"].(map[string]interface{})["id"].(float64))}
				EnemyPoints = append(EnemyPoints, &enemyStruct)
			}
			// 他のステージへの通り道であった場合、記憶しておく
			if c.(map[string]interface{})["portal"].(bool) == true {
				portalStruct := PortalStruct{X: j, Y: i, position: c.(map[string]interface{})["portal-data"].(map[string]interface{})["position"].(string), file: c.(map[string]interface{})["portal-data"].(map[string]interface{})["file"].(string)}
				if PortalPoints[j] == nil {
					PortalPoints[j] = make(map[int]*PortalStruct)
				}
				PortalPoints[j][i] = &portalStruct
			}
			j++
		}
		i++
	}

	for _, system := range ss.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				tileEntities = append(tileEntities, v)
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
