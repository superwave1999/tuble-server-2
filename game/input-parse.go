package game

import (
	"encoding/json"
	mapBlock "main/game/map-block"
)

func HttpInputToMap(bytes []byte) (MapInput, error) {
	var mapInput MapInput
	err := json.Unmarshal(bytes, &mapInput)
	if err != nil {
		return mapInput, err
	}
	return mapInput, nil
}

type MapInput struct {
	Map    [][]mapBlock.Block
	TimeMs int
	Moves  int
}
