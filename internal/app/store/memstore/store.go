package memstore

import "github.com/parMaster/logserver/internal/app/model"

type Store struct {
	messages map[int]model.Message
}

func (m *Store) Read() model.Message {
	return m.messages[0]
}

func (m *Store) Write(msg model.Message) int {
	m.messages[msg.ID] = msg
	return msg.ID
}
