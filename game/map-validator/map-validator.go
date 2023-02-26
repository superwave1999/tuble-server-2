package map_validator

import (
	"errors"
	"main/config"
	"main/game"
	"main/game/extensions"
	mapBlock "main/game/map-block"
)

const MapSizeHardLimit = 20

type Validation struct {
	TimeMs            int
	Moves             int
	NetTimeMs         int
	NetMoves          int
	MovePenaltyBlocks uint
	MoveBenefitBlocks uint
	TimePenaltyBlocks uint
	TimeBenefitBlocks uint
}

func findStartingPoint(gameMap [][]mapBlock.Block) *mapBlock.Block {
	var activeBlock *mapBlock.Block
out:
	for valueX := 0; valueX < len(gameMap); valueX++ {
		for valueY := 0; valueY < len(gameMap[valueX]); valueY++ {
			if gameMap[valueX][valueY].Type == mapBlock.TypeEndpoint && gameMap[valueX][valueY].Connections[0] == mapBlock.NoConnection {
				//Select starting endpoint.
				activeBlock = &gameMap[valueX][valueY]
				break out
			}
		}
	}
	return activeBlock
}

func getLengthsWithLimit(gameMap [][]mapBlock.Block) (int8, int8) {
	lenX := len(gameMap)
	currMinLenY := 0
	for i, coordsY := range gameMap {
		length := len(coordsY)
		if i == 0 || length < currMinLenY {
			currMinLenY = length
		}
	}
	return int8(lenX), int8(currMinLenY)
}

func Validate(input game.MapInput) (Validation, error) {
	validation := Validation{
		TimeMs:            input.TimeMs,
		Moves:             input.Moves,
		NetTimeMs:         input.TimeMs,
		NetMoves:          input.Moves,
		TimeBenefitBlocks: 0,
		TimePenaltyBlocks: 0,
		MoveBenefitBlocks: 0,
		MovePenaltyBlocks: 0,
	}

	lenX, lenY := getLengthsWithLimit(input.Map)
	if lenX > MapSizeHardLimit || lenY > MapSizeHardLimit {
		return validation, errors.New("submitted map is too large")
	}
	edgesX := [2]int8{0, lenX - 1}
	edgesY := [2]int8{0, lenY - 1}

	//Find starting point.
	prevCoords := [2]int8{mapBlock.NoConnection, mapBlock.NoConnection}
	activeBlock := findStartingPoint(input.Map)
	if activeBlock == nil {
		return validation, errors.New("submitted map has no starting point")
	}
	coords := [2]int8{activeBlock.X, activeBlock.Y}
	var troddenPath [][2]int8
	for {
		isFirstBlock := prevCoords == [2]int8{mapBlock.NoConnection, mapBlock.NoConnection}
		if !mapBlock.CoordsAreWithinLimits(coords, edgesX, edgesY) || (!isFirstBlock && !activeBlock.IsConnectedFrom(prevCoords)) || (troddenPath != nil && (extensions.ArraySearch(coords, troddenPath) >= 0)) {
			return validation, errors.New("submitted map is not connected")
		}
		//Current block is in limits and connected to previous.
		switch activeBlock.Type {
		case mapBlock.TypeEndpoint:
			if !isFirstBlock {
				if activeBlock.Connections[1] == mapBlock.NoConnection {
					//Both ends meet, all OK.
					return validation, nil
				} else {
					//Second connection is not -1. Something fishy.
					return validation, errors.New("submitted map has ghosts")
				}
			}
			break
		case mapBlock.TypeBenefitMoves:
			validation.NetMoves -= int(config.Game.ScoreMoveBenefit)
			validation.MoveBenefitBlocks++
			break
		case mapBlock.TypeBenefitTime:
			validation.NetTimeMs -= int(config.Game.ScoreTimeBenefitMs)
			validation.TimeBenefitBlocks++
			break
		case mapBlock.TypePenaltyMoves:
			validation.NetMoves += int(config.Game.ScoreMovePenalty)
			validation.MovePenaltyBlocks++
			break
		case mapBlock.TypePenaltyTime:
			validation.NetTimeMs += int(config.Game.ScoreTimePenaltyMs)
			validation.TimePenaltyBlocks++
			break
		}
		troddenPath = append(troddenPath, coords)
		prevCoordinateForConnection := prevCoords
		prevCoords = coords
		coords = activeBlock.NextConnectedBlockCoords(prevCoordinateForConnection)
		activeBlock = &input.Map[coords[0]][coords[1]]
	}
}
