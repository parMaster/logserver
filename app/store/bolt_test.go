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

	s.CleanUp()

	testData := Data{
		Module:   "TestModuleName",
		DateTime: "2022-01-02 03:04:05",
		Topic:    "TestTopicString",
		Value:    "TestValueString",
	}
	err = s.Write(testData)
	assert.NoError(t, err)

	savedMessage, readErr := s.Read(testData.Module)
	assert.NoError(t, readErr)
	assert.NotEmpty(t, savedMessage)
	assert.Equal(t, testData, savedMessage[0])

	// Write the same record again
	err = s.Write(testData)
	assert.NoError(t, err)

	// Write update
	testData.Value = "TestValueStringUpdated"
	err = s.Write(testData)
	assert.NoError(t, err)

	// Read
	savedMessage, readErr = s.Read(testData.Module)
	assert.NoError(t, readErr)
	assert.NotEmpty(t, savedMessage)
	assert.Equal(t, testData, savedMessage[0])

}
