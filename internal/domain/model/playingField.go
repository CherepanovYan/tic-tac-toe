package model

const (
	Empty = iota
	PlayerX
	PlayerO
)

const SizeField = 3

type PlayingField [SizeField][SizeField]int

func NewPlayingField() PlayingField {
	return PlayingField{}
}

func (p *PlayingField) Set(i, j int, player int) bool {
	if i >= SizeField || i < 0 || j < 0 || j >= SizeField {
		return false
	}
	if p[i][j] != Empty {
		return false
	}
	p[i][j] = player
	return true
}

func (p PlayingField) Copy() PlayingField {
	var copy PlayingField
	for i := 0; i < SizeField; i++ {
		for j := 0; j < SizeField; j++ {
			copy[i][j] = p[i][j]
		}
	}
	return copy
}
