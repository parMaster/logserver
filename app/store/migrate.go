package store

import (
	"context"
	"os"

	log "github.com/go-pkgz/lgr"
)

func Migrate(ctx context.Context) error {

	// read 'DEBUG' flag from env
	if os.Getenv("DEBUG") != "" {
		// example: DEBUG=1 go test -v ./app/store -run Test_Bolt_Store
		log.Setup(log.Debug, log.Msec, log.LevelBraces)
		log.Printf("[WARN] Debug mode is on")
	}

	// init sqlite storage
	s, err := NewSQLite(ctx, "file:./store/mqttdata.db?mode=rwc")
	if err != nil {
		log.Fatalf("[ERROR] Failed to open SQLite storage: %e", err)
	}

	// init bolt storage
	b, err := NewBolt(ctx, "/mnt/ramdisk/mqttdata.bolt")
	if err != nil {
		log.Fatalf("[ERROR] Failed to open Bolt storage: %e", err)
	}

	// prepare bolt to receive data
	b.CleanUp()

	i := 0
	modules := [...]string{"cave"}
	for _, m := range modules {
		data, err := s.View(m)
		if err != nil {
			log.Printf("[ERROR] Failed to read data from SQLite storage: %e", err)
		}
		for topic, values := range data {
			for dt, v := range values {

				// truncate value to 5 chars because i'm sure there are floats that are too long
				if len(v) > 5 {
					v = v[:5]
				}

				err = b.Write(Data{Module: m, DateTime: dt, Topic: topic, Value: v})
				if err != nil {
					log.Printf("[ERROR] Failed to write data to Bolt storage: %e", err)
				}

				i++
				if i%1000 == 0 {
					log.Printf("[DEBUG] %d records migrated", i)
					// os.Exit(0)
				}

				select {
				case <-ctx.Done():
					log.Printf("[DEBUG] Migrate cancelled")
					return nil
				default:
				}
			}
		}
	}

	return nil
}
