package memstore_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/parMaster/logserver/internal/app/model"
	"github.com/parMaster/logserver/internal/app/store/memstore"
	"github.com/stretchr/testify/assert"
)

func TestWriteRead(t *testing.T) {
	s := memstore.NewStore()

	mess := model.Message{
		ID:       1,
		DateTime: "2022-01-02 03:04:05",
		Topic:    "Test Topic",
		Message:  "Test Message",
	}

	id, err := s.Write(mess)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	savedMessage, readErr := s.Read(id)
	assert.NoError(t, readErr)
	assert.Equal(t, &mess, savedMessage)

	id, err = s.Write(mess)
	assert.NoError(t, err)
	assert.Equal(t, 2, id)
}
