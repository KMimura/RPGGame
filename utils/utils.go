package utils

// AbstractionValue 抽象化するために用いる値
var AbstractionValue = 50

// SimpleAbstractionValue 単純に足し/引きで抽象化するときに用いる値
var SimpleAbstractionValue = 50

var obstaclePoints map[int][]int

// Contains 配列に値が含まれているかどうかを判断
func Contains(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// SetObstaclePoints 障害物のある座標をセットする
func SetObstaclePoints(val map[int][]int) {
	obstaclePoints = val
}

// CheckIfPassable 障害物があるかどうかを判断
func CheckIfPassable(x, y int) bool {
	if obstaclePoints[x] != nil {
		if Contains(obstaclePoints[x], y) {
			return false
		}
	}
	return true
}
