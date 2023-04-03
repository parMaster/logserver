package logserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-chi/chi/v5"
	"github.com/parMaster/logserver/config"
	"github.com/parMaster/logserver/internal/app/store"
)

type LogServer struct {
	ctx    context.Context
	config config.Config
	mq     mqtt.Client
	store  store.Storer
}

func NewLogServer(ctx context.Context, config config.Config) *LogServer {
	l := &LogServer{
		ctx:    ctx,
		config: config,
	}

	// Inititalize message queue
	var err error
	l.mq, err = l.newMqClient()
	if err != nil {
		log.Fatalf("Can't configure mqtt client %e", err)
	}

	// Initialize database
	err = store.Load(ctx, config, &l.store)
	if err != nil {
		log.Fatalf("Can't configure database %e", err)
	}

	// db, err := newDB(config.DatabaseURL)
	// if err != nil {
	// 	return err
	// }

	// go s.CandelizeMinutely()

	return l
}

func (l *LogServer) newMqClient() (mqtt.Client, error) {

	opts := mqtt.NewClientOptions().AddBroker(l.config.MqBrokerURL)
	opts.SetUsername(l.config.MqUser)
	opts.SetPassword(l.config.MqPassword)
	opts.SetClientID(l.config.MqClientId)
	opts.SetCleanSession(true)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[ERROR] failed to connect to mqtt: %s", token.Error())
		return nil, token.Error()
	}

	// subscribe to root topic and all subtopics
	if token := c.Subscribe(l.config.MqRootTopic, 1, l.HandleMessage); token.Wait() && token.Error() != nil {
		log.Printf("[ERROR] failed to subscribe: %s", token.Error())
		return nil, token.Error()
	}

	log.Printf("[INFO] Successfuly connected to mqtt")
	return c, nil
}

func (l *LogServer) Start() error {
	httpServer := &http.Server{
		Addr:              l.config.DatabaseURL,
		Handler:           l.router(),
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Second,
	}

	httpServer.ListenAndServe()

	<-l.ctx.Done()
	log.Printf("[INFO] Terminating http server")

	if err := httpServer.Close(); err != nil {
		log.Printf("[ERROR] failed to close http server, %v", err)
	}
	return nil
}

func (l *LogServer) router() http.Handler {
	router := chi.NewRouter()

	router.Get("/api/v1/check", l.HandleCheck)

	return router
}

func (l *LogServer) HandleCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /api/v1/check")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		log.Printf("[ERROR] %s", err.Error())
	}
}

func (l *LogServer) HandleMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("INFO [%s] \t %s\r\n", msg.Topic(), msg.Payload())

	// parse message and save to database

}
