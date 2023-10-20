package store

import (
	"context"
	"fmt"
	"os"
	"testing"

	log "github.com/go-pkgz/lgr"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_Bolt_Store(t *testing.T) {

	// read 'DEBUG' flag from env
	if os.Getenv("DEBUG") != "" {
		// example: DEBUG=1 go test -v ./app/store -run Test_Bolt_Store
		log.Setup(log.Debug, log.Msec, log.LevelBraces)
		log.Printf("[WARN] Debug mode is on")
	}

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

	// Read not found
	noSuchBucketMessages, noSuchBucketErr := s.Read("No Such Bucket")
	assert.Error(t, noSuchBucketErr)
	assert.Empty(t, noSuchBucketMessages)
	assert.Equal(t, noSuchBucketErr, fmt.Errorf("bucket %q not found", "No Such Bucket"))

	// Write another record
	rec2 := Data{
		Module:   "TestModuleName",
		DateTime: "2022-01-02 03:04:05",
		Topic:    "TestTopicString2",
		Value:    "TestValueString2",
	}

	err = s.Write(rec2)
	assert.NoError(t, err)

	// View
	viewMessages, viewErr := s.View(rec2.Module)
	assert.NoError(t, viewErr)
	assert.NotEmpty(t, viewMessages)
	assert.Equal(t, testData.Value, viewMessages[testData.Topic][testData.DateTime])
	assert.Equal(t, rec2.Value, viewMessages[rec2.Topic][rec2.DateTime])

}
