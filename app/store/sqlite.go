package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	DB            *sql.DB
	ctx           context.Context
	activeModules map[string]bool
}

func NewSQLite(ctx context.Context, path string) (*SQLiteStorage, error) {

	sqliteDatabase, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		sqliteDatabase.Close()
	}()

	return &SQLiteStorage{DB: sqliteDatabase, ctx: ctx, activeModules: make(map[string]bool)}, nil
}

func (s *SQLiteStorage) Write(d Data) error {

	if ok, err := s.moduleActive(d.Module); err != nil || !ok {
		return err
	}

	if d.DateTime == "" {
		d.DateTime = time.Now().Format("2006-01-02 15:04")
	}

	if d.Topic == "" {
		return errors.New("topic is empty")
	}

	q := fmt.Sprintf("INSERT INTO `%s` VALUES ($1, $2, $3)", d.Module)

	_, err := s.DB.ExecContext(s.ctx, q, d.DateTime, d.Topic, d.Value)
	return err
}

// Read reads records for the given module from the database
func (s *SQLiteStorage) Read(module string) (data []Data, err error) {

	q := fmt.Sprintf("SELECT * FROM `%s`", module)
	rows, err := s.DB.QueryContext(s.ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		d := Data{Module: module}
		err = rows.Scan(&d.DateTime, &d.Topic, &d.Value)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return
}

// View returns a map of topics and their values for the given module
// The map is sorted by DateTime and structured as follows:
// map[Topic]map[DateTime]Value
func (s *SQLiteStorage) View(module string) (data map[string]map[string]string, err error) {

	data = make(map[string]map[string]string)

	// select distinct topics from module
	q := fmt.Sprintf("SELECT DISTINCT Topic FROM `%s`", module)
	rows, err := s.DB.QueryContext(s.ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var topic string
		err = rows.Scan(&topic)
		if err != nil {
			return nil, err
		}
		data[topic] = make(map[string]string)
	}

	// select all records from module and fill the map
	q = fmt.Sprintf("SELECT DateTime, Topic, AVG(Value) FROM `%s` GROUP BY DateTime, Topic order by DateTime", module)
	rows, err = s.DB.QueryContext(s.ctx, q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		d := Data{Module: module}
		err = rows.Scan(&d.DateTime, &d.Topic, &d.Value)
		if err != nil {
			return nil, err
		}
		data[d.Topic][d.DateTime] = d.Value
	}

	return
}

// Check if the table exists, create if not. Cache the result in the map
func (s *SQLiteStorage) moduleActive(module string) (bool, error) {

	if module == "" {
		return false, errors.New("module name is empty")
	}

	if s.activeModules[module] {
		return true, nil
	}

	if _, ok := s.activeModules[module]; !ok {
		q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (DateTime TEXT, Topic TEXT, Value TEXT)", module)
		_, err := s.DB.ExecContext(s.ctx, q)
		if err != nil {
			return false, err
		}
		s.activeModules[module] = true
	}

	return true, nil
}

// Cleanup removes the table for the given module
func (s *SQLiteStorage) Cleanup(module string) {
	q := fmt.Sprintf("DROP TABLE `%s`", module)
	s.DB.Exec(q)
	delete(s.activeModules, module)
}
