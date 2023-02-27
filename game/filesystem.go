package game

import (
	"fmt"
	"os"
)

const MapFileName = "storage/.map_cache.json"

func ReadFile() []byte {
	data, err := os.ReadFile(MapFileName)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func WriteFile(contents []byte) {
	err := os.WriteFile(MapFileName, contents, 0640)
	if err != nil {
		fmt.Println(err)
	}
}
