package map_block

func CoordsAreWithinLimits(coords [2]int8, edgesX [2]int8, edgesY [2]int8) bool {
	return (coords[0] >= edgesX[0] && coords[0] <= edgesX[1]) && (coords[1] >= edgesY[0] && coords[1] <= edgesY[1])
}

func ConnectionToCoords(connection int8, currentX int8, currentY int8) [2]int8 {
	coords := [2]int8{NoConnection, NoConnection}
	switch connection {
	case 3: //To left.
		coords = [2]int8{currentX - 1, currentY}
		break
	case 2: //TO bottom.
		coords = [2]int8{currentX, currentY + 1}
		break
	case 1:
		coords = [2]int8{currentX + 1, currentY}
		break
	case 0:
		coords = [2]int8{currentX, currentY - 1}
		break
	}
	return coords
}
