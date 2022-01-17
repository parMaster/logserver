package store

import "github.com/parMaster/logserver/internal/app/model"

type Storer interface {
	Read(id int) (*model.Message, error)
	Write(msg model.Message) (int, error)
}
