package model

const (
    Empty   = 0
    PlayerX = 1 
    PlayerO = 2 
)

type playingField struct{
	Field [3][3]int `json:"field"`
}

func NewPlayingField() *playingField{
        return &playingField{}
}