package sqlstore_test

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/parMaster/logserver/internal/app/model"
	"github.com/parMaster/logserver/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost user=db_user password=pwd dbname=mqttdata sslmode=disable"
	}
	os.Exit(m.Run())
}

func TestWriteRead(t *testing.T) {

	db, _ := sqlstore.TestDB(t, databaseURL)

	// todo - setup and teardown
	// db, teardown := sqlstore.TestDB(t, databaseURL)
	// defer teardown("rawdata")

	s := sqlstore.NewStore(db)

	mess := model.Message{
		ID:       1,
		DateTime: "2022-01-02T03:04:05Z",
		Topic:    "Test Topic",
		Message:  "Test Message",
	}

	id, err := s.Write(mess)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)

	savedMessage, readErr := s.Read(id)
	assert.NoError(t, readErr)
	assert.Equal(t, mess.DateTime, savedMessage.DateTime)
	assert.Equal(t, mess.Topic, savedMessage.Topic)
	assert.Equal(t, mess.Message, savedMessage.Message)

}
