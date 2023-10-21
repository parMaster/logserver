package store

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	log "github.com/go-pkgz/lgr"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var debug bool

func init() {
	log.Setup(log.Msec, log.LevelBraces)
	if os.Getenv("DEBUG") != "" {
		// example: DEBUG=1 go test -v ./app/store -run Test_Bolt_Store
		log.Setup(log.Debug, log.Msec, log.LevelBraces)
		debug = true
		log.Printf("[WARN] Debug mode is on")
	}
}

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

/*
DB_FILE=/mnt/ramdisk/test.bolt go test -v  -run ^TestWrite$
*/
func TestWrite(t *testing.T) {

	filename := "/tmp/test.bolt"
	if os.Getenv("DB_FILE") != "" {
		filename = os.Getenv("DB_FILE")
	}

	// Kinda benchmark
	ctx := context.Background()
	s, err := NewBolt(ctx, filename+time.Now().Format("20060102_150405.999999999"))
	assert.NoError(t, err)
	s.CleanUp()

	log.Printf("[INFO] Writing 1000 records")
	for i := 0; i < 1000; i++ {
		s.Write(Data{Module: "bench_write", DateTime: time.Now().Format("2006-01-02 15:04:05.999999999"), Topic: fmt.Sprintf("topic %d", rand.Uint64()), Value: fmt.Sprintf("value %d", rand.Uint64())})
	}
	log.Printf("[INFO] Done")

	log.Printf("[INFO] Reading 100 times")
	for i := 0; i < 100; i++ {
		data, err := s.Read("bench_write")
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	}
	log.Printf("[INFO] Done")

}

/*
DB_FILE=/mnt/ramdisk/test.bolt go test -v  -run ^TestParallelWrite$
*/
func TestParallelWrite(t *testing.T) {

	filename := "/tmp/test.bolt"
	if os.Getenv("DB_FILE") != "" {
		filename = os.Getenv("DB_FILE")
	}

	ctx := context.Background()
	s, err := NewBolt(ctx, filename+time.Now().Format("20060102_150405.999999999"))
	assert.NoError(t, err)
	s.CleanUp()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	N := 100

	go func() {
		log.Printf("[INFO] Writing %d records", N)
		for i := 0; i < N; i++ {
			s.Write(Data{Module: "bench_write", DateTime: time.Now().Format("2006-01-02 15:04:05.999999999"), Topic: fmt.Sprintf("topic %d", rand.Uint64()), Value: fmt.Sprintf("value %d", rand.Uint64())})
		}
		log.Printf("[INFO] Done")
		wg.Done()
	}()

	go func() {
		log.Printf("[INFO] Writing %d records", N)
		for i := 0; i < N; i++ {
			s.Write(Data{Module: "bench_write", DateTime: time.Now().Format("2006-01-02 15:04:05.999999999"), Topic: fmt.Sprintf("topic %d", rand.Uint64()), Value: fmt.Sprintf("value %d", rand.Uint64())})
		}
		log.Printf("[INFO] Done")
		wg.Done()
	}()

	wg.Wait()

}
