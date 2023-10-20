package store

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_Bolt_Store(t *testing.T) {

	ctx := context.Background()

	s, err := NewBolt(ctx, "/tmp/test.bolt")
	assert.NoError(t, err)

	testData := Data{
		Module:   "Test Module",
		DateTime: "2022-01-02 03:04:05",
		Topic:    "Test Topic",
		Value:    "Test Value",
	}
	err = s.Write(testData)
	assert.NoError(t, err)

	// savedMessage, readErr := s.Read("Test Module")
	// assert.NoError(t, readErr)
	// assert.NotEmpty(t, savedMessage)
	// assert.Equal(t, mess, savedMessage[0])

}
