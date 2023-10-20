package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/go-pkgz/lgr"

	bolt "go.etcd.io/bbolt"
)

// Bolt is a storage implementation (Storer Interface) that uses BoltDB as a backend.
type Bolt struct {
	db            *bolt.DB
	ctx           context.Context
	activeModules map[string]bool
}

func NewBolt(ctx context.Context, dbFile string) (b *Bolt, err error) {

	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second}) // nolint
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		db.Close()
	}()

	return &Bolt{db: db, ctx: ctx, activeModules: make(map[string]bool)}, nil
}

func (b *Bolt) Read(module string) ([]Data, error) {

	result := []Data{}
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(module))
		if b == nil {
			return fmt.Errorf("bucket %q not found", module)
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			val := Data{}
			err := json.Unmarshal(v, &val)
			if err != nil {
				return err
			}
			val.Module = module
			result = append(result, val)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] BoltDB read from [%s]: %v", module, result)
	return result, nil
}

func (b *Bolt) Write(data Data) error {

	if ok, err := b.moduleActive(data.Module); err != nil || !ok {
		return fmt.Errorf("module %s is not active, possible activation error: %e", data.Module, err)
	}

	if data.DateTime == "" {
		data.DateTime = time.Now().Format("2006-01-02 15:04")
	}

	if data.Topic == "" {
		return fmt.Errorf("topic is empty")
	}

	key := data.Topic + "-" + data.DateTime

	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(data.Module))
		if b == nil {
			return fmt.Errorf("bucket %q not found", data.Module)
		}

		jdata, jerr := json.Marshal(Data{Topic: data.Topic, DateTime: data.DateTime, Value: data.Value})
		if jerr != nil {
			return jerr
		}

		return b.Put([]byte(key), jdata)
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] BoltDB saved to [%s]: key: %s,\t v: %s", data.Module, key, data)

	return nil
}

// View returns a map of topics and their values for the given module
// The map is sorted by DateTime and structured as follows:
// map[Topic]map[DateTime]Value
func (b *Bolt) View(module string) (data map[string]map[string]string, err error) {

	data = make(map[string]map[string]string)

	// Get all the records from the bucket module
	records, err := b.Read(module)
	if err != nil {
		return nil, err
	}
	for _, d := range records {
		if _, ok := data[d.Topic]; !ok {
			data[d.Topic] = make(map[string]string)
		}
		data[d.Topic][d.DateTime] = d.Value
	}

	return
}

// Check if the table exists, create if not. Cache the result in the map
func (b *Bolt) moduleActive(module string) (bool, error) {

	if module == "" {
		return false, fmt.Errorf("module name is empty")
	}

	if b.activeModules[module] {
		return true, nil
	}

	if _, ok := b.activeModules[module]; !ok {

		// create a bucket for the module
		err := b.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(module))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})
		if err != nil {
			return false, err
		}
		b.activeModules[module] = true
		log.Printf("[DEBUG] BoltDB bucket created: %s", module)
	}

	return true, nil
}

// CleanUp removes all the data from the storage
func (b *Bolt) CleanUp() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			return tx.DeleteBucket(name)
		})
	})
}
