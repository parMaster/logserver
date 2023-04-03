package store

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/parMaster/logserver/config"
	"github.com/parMaster/logserver/internal/app/model"
	"github.com/parMaster/logserver/internal/app/store/sqlite"
)

type Storer interface {
	// Read reads records for the given module from the database.
	Read(context.Context, string) ([]model.Data, error)
	// Write writes the data to the database.
	Write(context.Context, model.Data) error
	// View returns the data for the given module in the format that is suitable for the web view.
	View(context.Context, string) (map[string]map[string]string, error)
}

func Load(ctx context.Context, cfg config.Config, s *Storer) error {
	var err error
	switch cfg.DatabaseKind {
	case "sqlite":
		*s, err = sqlite.NewStorage(ctx, cfg.DatabaseURL)
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
