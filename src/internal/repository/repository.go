package repository

import "sync"

type repository struct {
	games sync.Map
}

func NewRepository() *repository {
	return &repository{}
}
