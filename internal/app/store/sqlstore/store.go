package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/parMaster/logserver/internal/app/model"
	"github.com/parMaster/logserver/internal/app/store"
)

type Store struct {
	db       *sql.DB
	messages map[int]*model.Message
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:       db,
		messages: make(map[int]*model.Message),
	}
}

func (m *Store) Read(id int) (*model.Message, error) {

	mess := &model.Message{}

	if err := m.db.QueryRow(
		"SELECT id, dt, topic, message FROM rawdata WHERE id = $1",
		id,
	).Scan(
		&mess.ID,
		&mess.DateTime,
		&mess.Topic,
		&mess.Message,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return mess, nil
}

func (m *Store) Write(msg model.Message) (int, error) {
	var id int

	err := m.db.QueryRow(
		"INSERT INTO rawdata (dt, topic, message) VALUES ($1, $2, $3) RETURNING id",
		msg.DateTime,
		msg.Topic,
		msg.Message,
	).Scan(&id)

	return id, err
}
