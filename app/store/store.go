package store

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/parMaster/logserver/config"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrPretendToCandelize = errors.New("pretending to candelize data")
)

type Data struct {
	Module   string
	DateTime string
	Topic    string
	Value    string
}

type Message struct {
	ID       int
	DateTime string
	Topic    string
	Message  string
}

type Storer interface {
	// Read reads records for the given module from the database.
	Read(string) ([]Data, error)
	// Write writes the data to the database.
	Write(Data) error
	// View returns the data for the given module in the format that is suitable for the web view.
	View(string) (map[string]map[string]string, error)
}

func Load(ctx context.Context, cfg config.Config, s *Storer) error {
	var err error
	switch cfg.DatabaseKind {
	case "bolt":
		*s, err = NewBolt(ctx, cfg.DatabaseURL)
		if err != nil {
			return fmt.Errorf("failed to init SQLite storage: %e", err)
		}
	case "sqlite":
		*s, err = NewSQLite(ctx, cfg.DatabaseURL)
		if err != nil {
			return fmt.Errorf("failed to init SQLite storage: %e", err)
		}
	case "":
		log.Printf("[DEBUG] Storage is not configured")
		return errors.New("storage is not configured")
	default:
		return fmt.Errorf("storage type %s is not supported", cfg.DatabaseKind)
	}
	return err
}
