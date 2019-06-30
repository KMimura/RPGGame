package systems

// Contains 配列に値が含まれているかどうかを判断
func contains(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// checkIfPassable 障害物があるかどうかを判断
func checkIfPassable(x, y int) bool {
	if ObstaclePoints[x] != nil {
		if contains(ObstaclePoints[x], y) {
			return false
		}
	}
	return true
}

// checkIfPortal 通り道かどうかを判断
func checkIfPortal(x, y int) (bool, *PortalStruct) {
	if PortalPoints[x] != nil {
		if PortalPoints[x][y] != nil {
			return true, PortalPoints[x][y]
		}
	}
	return false, nil
}
