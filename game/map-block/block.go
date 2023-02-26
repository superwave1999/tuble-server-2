package map_block

import (
	"github.com/mroth/weightedrand"
	"math/rand"
)

const TypeNormal string = "n"
const TypeEndpoint string = "e"
const TypePenaltyTime string = "x"
const TypePenaltyMoves string = "z"
const TypeBenefitTime string = "a"
const TypeBenefitMoves string = "b"
const NoConnection int8 = -1

type Block struct {
	Type        string
	X           int8
	Y           int8
	Connections [2]int8
}

func Create(coordinateX int8, coordinateY int8) Block {
	return Block{
		Type: TypeNormal,
		X:    coordinateX,
		Y:    coordinateY,
		Connections: [2]int8{
			NoConnection,
			NoConnection,
		},
	}
}

func (b *Block) SetStartingBlock(nextCoords [2]int8) {
	b.Type = TypeEndpoint
	b.SetFirstConnection(NoConnection)
	b.SetConnectionToCoords(nextCoords, true)
}

func (b *Block) SetConnectionToCoords(nextCoords [2]int8, second bool) {
	f := b.SetFirstConnection
	if second {
		f = b.SetSecondConnection
	}

	diffX := b.X - nextCoords[0]
	diffY := b.Y - nextCoords[1]

	if diffX > 0 {
		f(3) //Incoming from left.
	}
	if diffX < 0 {
		f(1) //Incoming from right.
	}
	if diffY > 0 {
		f(0) //Incoming from top.
	}
	if diffY < 0 {
		f(2) //Incoming from bottom.
	}
}

func (b *Block) SetRandomSpecialType() {
	chooser, _ := weightedrand.NewChooser(
		weightedrand.Choice{Item: TypeNormal, Weight: 36},
		weightedrand.Choice{Item: TypePenaltyTime, Weight: 4},
		weightedrand.Choice{Item: TypePenaltyMoves, Weight: 4},
		weightedrand.Choice{Item: TypeBenefitTime, Weight: 3},
		weightedrand.Choice{Item: TypeBenefitMoves, Weight: 3},
	)
	b.Type = chooser.Pick().(string)
}

func (b *Block) SetRandomConnections() {
	a := rand.Perm(4)
	b.SetFirstConnection(int8(a[0]))
	b.SetSecondConnection(int8(a[1]))
}

func (b *Block) incrementConnections() {
	for i, val := range b.Connections {
		n := val + 1
		if n > 3 {
			n = 0
		}
		b.Connections[i] = n
	}
}

func (b *Block) RandomRotate() {
	incr := rand.Intn(4) //0-3
	for t := 0; t < incr; t++ {
		b.incrementConnections()
	}
}

func (b *Block) SetFirstConnection(c int8) {
	b.Connections[0] = c
}

func (b *Block) SetSecondConnection(c int8) {
	b.Connections[1] = c
}

func (b *Block) SetType(t string) {
	b.Type = t
}

func (b *Block) IsEmpty() bool {
	return b.Connections[0] == NoConnection && b.Connections[1] == NoConnection
}

func (b *Block) IsMovable() bool {
	return b.Connections[0] != NoConnection && b.Connections[1] != NoConnection
}

//Map verification.

func (b *Block) NextConnectedBlockCoords(existingCoords [2]int8) [2]int8 {
	//Pick the connection that isn't the existing one.
	use := ConnectionToCoords(b.Connections[0], b.X, b.Y)
	if existingCoords == use {
		use = ConnectionToCoords(b.Connections[1], b.X, b.Y)
	}
	return use
}

func (b *Block) IsConnectedFrom(prevCoords [2]int8) bool {
	//Must be able to handle both orientations of pipes.
	return (prevCoords == ConnectionToCoords(b.Connections[0], b.X, b.Y)) || (prevCoords == ConnectionToCoords(b.Connections[1], b.X, b.Y))
}
