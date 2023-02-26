package config

import (
	"os"
	"strconv"
)

var Game *GameConfig

type GameConfig struct {
	ForceEdgeFinish    bool
	Size               uint8
	MinConnected       uint8
	ScoreMoveBenefit   uint32
	ScoreMovePenalty   uint32
	ScoreTimeBenefitMs uint32
	ScoreTimePenaltyMs uint32
}

func GetGameConfig() {
	c := GameConfig{
		ForceEdgeFinish:    true,
		Size:               6,
		MinConnected:       13,
		ScoreMoveBenefit:   3,
		ScoreMovePenalty:   2,
		ScoreTimeBenefitMs: 5000,
		ScoreTimePenaltyMs: 12000,
	}
	parseBool, err := strconv.ParseBool(os.Getenv("MAP_EDGE_FINISH"))
	if err == nil {
		c.ForceEdgeFinish = parseBool
	}
	vint8, err := strconv.ParseUint(os.Getenv("MAP_SIZE"), 10, 8)
	if err == nil {
		c.Size = uint8(vint8)
	}
	vint8, err = strconv.ParseUint(os.Getenv("MAP_MIN_CONNECTED_BLOCKS"), 10, 8)
	if err == nil {
		c.MinConnected = uint8(vint8)
	}
	vint32, err := strconv.ParseUint(os.Getenv("SCORE_MOVE_BENEFIT"), 10, 32)
	if err == nil {
		c.ScoreMoveBenefit = uint32(vint32)
	}
	vint32, err = strconv.ParseUint(os.Getenv("SCORE_MOVE_PENALTY"), 10, 32)
	if err == nil {
		c.ScoreMovePenalty = uint32(vint32)
	}
	vint32, err = strconv.ParseUint(os.Getenv("SCORE_TIME_BENEFIT_MS"), 10, 32)
	if err == nil {
		c.ScoreTimeBenefitMs = uint32(vint32)
	}
	vint32, err = strconv.ParseUint(os.Getenv("SCORE_TIME_PENALTY_MS"), 10, 32)
	if err == nil {
		c.ScoreTimePenaltyMs = uint32(vint32)
	}
	checkConfig(c)
	Game = &c
}

func checkConfig(c GameConfig) {
	if c.Size < 2 {
		panic("PANIC: Configured map size is too small!")
	}
	if c.MinConnected > (c.Size * c.Size) {
		panic("PANIC: Min connected blocks cannot be greater than size^2 !")
	}
}
