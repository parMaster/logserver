package memstore

import (
	"github.com/parMaster/logserver/internal/app/model"
	"github.com/parMaster/logserver/internal/app/store"
)

type Store struct {
	messages map[int]*model.Message
}

func NewStore() *Store {
	return &Store{
		messages: make(map[int]*model.Message),
	}
}

func (m *Store) Read(id int) (*model.Message, error) {
	elem, ok := m.messages[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return elem, nil
}

func (m *Store) Write(msg model.Message) (int, error) {
	msg.ID = len(m.messages) + 1
	m.messages[msg.ID] = &msg
	return msg.ID, nil
}

func (m *Store) CandelizePreviousMinute(sensor string) error {
	return store.ErrPretendToCandelize
}
