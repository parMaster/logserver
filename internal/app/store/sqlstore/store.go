package sqlstore

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/parMaster/logserver/internal/app/model"
	"github.com/parMaster/logserver/internal/app/store"
)

type Store struct {
	db       *sql.DB
	messages map[int]model.Message
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (m *Store) Read(id int) (*model.Message, error) {
	elem, ok := m.messages[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return &elem, nil
}

func (m *Store) Write(msg model.Message) int {
	var id int

	m.db.QueryRow(
		"INSERT INTO rawdata (dt, topic, message) VALUES ($1, $2, $3) RETURNING id",
		time.Now().Format("2006.01.02 15:04:05"),
		msg.Topic,
		msg.Message,
	).Scan(&id)

	return id
}
