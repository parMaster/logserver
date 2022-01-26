package logserver

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-pkgz/lgr"
	"github.com/gorilla/mux"
	"github.com/parMaster/logserver/internal/app/model"
	"github.com/parMaster/logserver/internal/app/store"
	"github.com/parMaster/logserver/internal/app/store/sqlstore"
)

type LogServer struct {
	logger *lgr.Logger
	router *mux.Router
	mq     *mqtt.Client
	store  store.Storer
}

func NewServer(store store.Storer, config Config) *LogServer {
	s := &LogServer{
		logger: lgr.New(),
		router: mux.NewRouter(),
		store:  store,
	}

	var err error
	if s.mq, err = s.configureMqClient(&config); err != nil {
		os.Exit(1)
	}

	s.router.HandleFunc("/check", s.HandleCheck())

	go s.CandelizeMinutely()

	return s
}

func (l *LogServer) CandelizeMinutely() {

	ticker := time.NewTicker(1 * time.Minute)
	for _ = range ticker.C {
		l.logger.Logf("Candelizing...")
		if err := l.store.CandelizePreviousMinute("croco/cave/temperature"); err != nil {
			l.logger.Logf("ERROR %s", err.Error())
		}
		if err := l.store.CandelizePreviousMinute("croco/cave/targetTemperature"); err != nil {
			l.logger.Logf("ERROR %s", err.Error())
		}
	}
}

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := sqlstore.NewStore(db)

	s := NewServer(store, *config)

	if err := http.ListenAndServe(config.BindAddr, s.router); err != nil {
		return err
	}
	return nil
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	db.Query("SET TIMEZONE TO 'Europe/Kiev';")

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (l *LogServer) configureMqClient(config *Config) (*mqtt.Client, error) {

	opts := mqtt.NewClientOptions().AddBroker(config.MqBrokerURL)
	opts.SetUsername(config.MqUser)
	opts.SetPassword(config.MqPassword)
	opts.SetClientID(config.MqClientId)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		l.logger.Logf("FATAL failed to connect to mqtt: %s", token.Error())
		return nil, token.Error()
	}

	if token := c.Subscribe("croco/#", 1, l.HandleMessage); token.Wait() && token.Error() != nil {
		l.logger.Logf("FATAL failed to subscribe: %s", token.Error())
		return nil, token.Error()
	}
	l.logger.Logf("INFO Successfuly connected to mqtt")
	return &c, nil
}

func (l *LogServer) HandleCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.logger.Logf("INFO HandleCheck called")
	}
}

func (l *LogServer) HandleMessage(client mqtt.Client, msg mqtt.Message) {
	l.logger.Logf("INFO [%s] \t %s\r\n", msg.Topic(), msg.Payload())

	if l.store != nil {
		l.store.Write(model.Message{
			ID:       0,
			DateTime: time.Now().Format("2006.01.02 15:04:05"),
			Topic:    msg.Topic(),
			Message:  string(msg.Payload()),
		})
	}

}
