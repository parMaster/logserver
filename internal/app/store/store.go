package store

import "github.com/parMaster/logserver/internal/app/model"

type Storer interface {
	Read() model.Message
	Write(msg model.Message) int
}
