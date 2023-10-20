package store

import (
	"context"
	"time"

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

	return &Bolt{db: db, ctx: ctx, activeModules: make(map[string]bool)}, nil
}

func (b *Bolt) Read(module string) (data []Data, err error) {
	return data, nil
}

func (b *Bolt) Write(data Data) error {
	return nil
}

// View returns a map of topics and their values for the given module
// The map is sorted by DateTime and structured as follows:
// map[Topic]map[DateTime]Value
func (b *Bolt) View(module string) (data map[string]map[string]string, err error) {

	data = make(map[string]map[string]string)

	return
}
