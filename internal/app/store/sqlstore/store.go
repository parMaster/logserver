package sqlstore

import (
	"database/sql"
	"fmt"

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

// Calculate temperature candle for previous minute and insert data into tempdata table
func (m *Store) CandelizePreviousMinute(sensor string) error {

	q := fmt.Sprintf(
		`INSERT INTO tempdata
	SELECT 
		NEXTVAL('tempdata_id_seq') as id, 
		DATE_PART('year', dt) as year, 
		DATE_PART('month', dt) as month, 
		DATE_PART('day', dt) as day, 
		DATE_PART('hour', dt) as hour, 
		DATE_PART('minute', dt) as minute, 
		MIN(message::float)::numeric(10,2) AS min_temp,
		AVG(message::float)::numeric(10,2) AS avg_temp,
		MAX(message::float)::numeric(10,2) AS max_temp,
		'' as strval,
		topic as sensor
	FROM rawdata 
	WHERE 
		topic = '%s' AND
		date_trunc('minute', dt) = date_trunc('minute', CURRENT_TIMESTAMP - interval '1 minute')
	GROUP BY year, month, day, hour, minute, topic
	ORDER BY year, month, day, hour, minute, topic;`, sensor)

	err := m.db.QueryRow(
		q,
	).Err()

	return err
}
