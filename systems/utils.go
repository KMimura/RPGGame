package systems

// Contains 配列に値が含まれているかどうかを判断
func Contains(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// CheckIfPassable 障害物があるかどうかを判断
func CheckIfPassable(x, y int) bool {
	if ObstaclePoints[x] != nil {
		if Contains(ObstaclePoints[x], y) {
			return false
		}
	}
	return true
}
