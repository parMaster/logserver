package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Returns system temp dir (i.e. /tmp on Linux, no trailing slash).
// If TEMP_DIR environment variable is set, it is returned instead
func tempDir() string {

	if os.Getenv("TEMP_DIR") != "" {
		return os.Getenv("TEMP_DIR")
	}

	return os.TempDir()
}

func Test_Sqlite_Store(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	store, err := NewSQLite(ctx, fmt.Sprintf("file:%s/test.db?cache=shared&mode=rwc", tempDir()))
	if err != nil {
		log.Printf("[ERROR] Failed to open SQLite storage: %e", err)
	}

	testRecord := Data{
		Module:   "testModule",
		DateTime: "2019-01-01 00:00",
		Topic:    "testTopic",
		Value:    "testValue",
	}

	store.Cleanup(testRecord.Module)

	// write a record
	err = store.Write(testRecord)
	assert.Nil(t, err)

	// read the record
	data, err := store.Read(testRecord.Module)
	assert.Equal(t, 1, len(data))
	assert.Nil(t, err)
	assert.Equal(t, data[0], testRecord)

	// write another record
	err = store.Write(testRecord)
	assert.Nil(t, err)
	data, err = store.Read(testRecord.Module)
	assert.Equal(t, 2, len(data))
	assert.Nil(t, err)
	assert.Equal(t, data[1], testRecord)

	// test if the module is not active (no such table)
	data, err = store.Read("notable")
	assert.Nil(t, data)
	assert.Error(t, err)

	// empty topic is not allowed
	err = store.Write(Data{Module: "testModule", Topic: "", Value: "testValue"})
	assert.Error(t, err)

	// empty value is allowed
	err = store.Write(Data{Module: "testModule", Topic: "testTopic", Value: ""})
	assert.NoError(t, err)

	// Test if the date time is set to the current time if it is not set.
	dt := time.Now().Format("2006-01-02 15:04")
	err = store.Write(Data{Module: "testModule", Topic: "testTopic", Value: "testValue"})
	assert.NoError(t, err)
	savedValues, err := store.Read("testModule")
	assert.NoError(t, err)
	assert.Equal(t, dt, savedValues[len(savedValues)-1].DateTime)

	v, _ := store.Read("testModule")
	n := 100
	// n records
	for i := 0; i < n; i++ {
		err = store.Write(testRecord)
		assert.NoError(t, err)
	}
	vn, _ := store.Read("testModule")
	assert.Equal(t, n, len(vn)-len(v))
}

func Test_SqliteStorage_View(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	store, err := NewSQLite(ctx, fmt.Sprintf("file:%s/test_view.db?mode=rwc", tempDir()))
	if err != nil {
		log.Printf("[ERROR] Failed to open SQLite storage: %e", err)
	}

	records := []Data{
		{Module: "view", DateTime: time.Now().Format("2006-01-02") + " 00:00", Topic: "temp", Value: "36000"},
		{Module: "view", DateTime: time.Now().Format("2006-01-02") + " 00:01", Topic: "temp", Value: "36100"},
		{Module: "view", DateTime: time.Now().Format("2006-01-02") + " 00:02", Topic: "temp", Value: "36200"},
		{Module: "view", DateTime: time.Now().Format("2006-01-02") + " 00:00", Topic: "rpm", Value: "100"},
		{Module: "view", DateTime: time.Now().Format("2006-01-02") + " 00:01", Topic: "rpm", Value: "200"},
		{Module: "view", DateTime: time.Now().Format("2006-01-02") + " 00:02", Topic: "rpm", Value: "300"},
	}

	for _, r := range records {
		err = store.Write(r)
		assert.NoError(t, err)
	}

	// test if the view is created

	viewExpected := map[string]map[string]string{
		"temp": {
			time.Now().Format("2006-01-02") + " 00:00": "36000",
			time.Now().Format("2006-01-02") + " 00:01": "36100",
			time.Now().Format("2006-01-02") + " 00:02": "36200",
		},
		"rpm": {
			time.Now().Format("2006-01-02") + " 00:00": "100",
			time.Now().Format("2006-01-02") + " 00:01": "200",
			time.Now().Format("2006-01-02") + " 00:02": "300",
		},
	}

	view, err := store.View("view") // create the view
	assert.NoError(t, err)
	assert.Equal(t, viewExpected, view)

}

func Test_SqliteStorage_readOnly(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	store, err := NewSQLite(ctx, fmt.Sprintf("file:%s/test_view.db?mode=ro", tempDir()))
	if err != nil {
		log.Printf("[ERROR] Failed to open SQLite storage: %e", err)
	}

	err = store.Write(Data{Module: "testModule", Topic: "testTopic", Value: "testValue"})
	assert.Error(t, err)
	assert.Equal(t, "attempt to write a readonly database", err.Error())

	s, err := NewSQLite(ctx, "file:/tmp/test_notcreated.db?mode=ro")
	assert.NotNil(t, s)
	assert.NoError(t, err)
	err = s.Write(Data{Module: "testModule", Topic: "testTopic", Value: "testValue"})
	assert.Error(t, err)
	assert.Equal(t, "unable to open database file: no such file or directory", err.Error())
}
