package store

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_Memory_Store(t *testing.T) {

	s := NewMemoryStore()

	// Read empty
	emptyMessages, emptyErr := s.Read("Test Module")
	assert.Error(t, emptyErr)
	assert.Empty(t, emptyMessages)
	assert.ErrorIs(t, emptyErr, ErrRecordNotFound)

	// Write
	testData := Data{
		Module:   "Test Module",
		DateTime: "2022-01-02 03:04:05",
		Topic:    "Test Topic",
		Value:    "Test Value",
	}
	err := s.Write(testData)
	assert.NoError(t, err)

	// Read
	savedMessages, readErr := s.Read("Test Module")
	assert.NoError(t, readErr)
	assert.NotEmpty(t, savedMessages)
	assert.Equal(t, testData, savedMessages[0])

	// Read not found
	noSuchBucketMessages, noSuchBucketErr := s.Read("No Such Bucket")
	assert.Error(t, noSuchBucketErr)
	assert.Empty(t, noSuchBucketMessages)
	assert.ErrorIs(t, noSuchBucketErr, ErrRecordNotFound)

	// View
	viewMessages, viewErr := s.View("Test Module")
	assert.NoError(t, viewErr)
	assert.NotEmpty(t, viewMessages)
	assert.Equal(t, testData.Value, viewMessages[testData.Topic][testData.DateTime])
}
