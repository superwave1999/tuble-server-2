package game

import (
	"fmt"
	"os"
)

const MapFileName = ".map_cache.json"

func ReadFile() []byte {
	data, err := os.ReadFile(MapFileName)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func WriteFile(contents []byte) {
	err := os.WriteFile(MapFileName, contents, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
