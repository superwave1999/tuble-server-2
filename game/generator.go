package game

import (
	"encoding/json"
	mapBlock "main/game/map-block"
	mapBuilder "main/game/map-builder"
	"strconv"
	"time"
)

var MapJson string

func mapStorageStruct(builtMap [][]mapBlock.Block) MapFile {
	utc := time.Now().UTC()
	return MapFile{
		Date:      utc.Format("2006-01-02"),
		Timestamp: strconv.FormatInt(utc.UnixMilli(), 10),
		Map:       builtMap,
	}
}

func NewMap() {
	newMap := mapStorageStruct(mapBuilder.New())
	marshal, _ := json.Marshal(newMap)
	WriteFile(marshal)
	MapJson = string(marshal)
}

type MapFile struct {
	Date      string
	Timestamp string
	Map       [][]mapBlock.Block
}
