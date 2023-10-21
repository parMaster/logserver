package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/parMaster/logserver/app/config"
	"github.com/parMaster/logserver/app/queue"
	"github.com/parMaster/logserver/app/store"
)

type Service struct {
	store.Storer
	q *queue.Client
}

// RunService consumes messages from mqtt queue and writes them to database
// It is intended to be run as a service/daemon
func RunService(ctx context.Context, config config.Config) {

	s := &Service{}

	// Initialize database
	var err error
	err = store.Load(ctx, config, &s.Storer)
	if err != nil {
		log.Fatalf("Can't configure database %e", err)
	}

	// Describe subscriptions
	var subs []queue.Subscription
	subs = append(subs, queue.Subscription{
		Topic:   "croco/cave/#",
		Handler: s.crocoCaveLogs,
	})

	if config.CollectRaw {
		subs = append(subs, queue.Subscription{
			Topic:   "ESP32-A473F53A7D80/p/ds18b20/#",
			Handler: s.crocoCaveRaw,
		})
	}

	// Inititalize message queue, subscribe to topics
	s.q, err = queue.NewClient(ctx, config, subs...)
	if err != nil {
		log.Fatalf("Can't configure mqtt client %e", err)
	}

	<-ctx.Done()
	log.Printf("[INFO] Terminating service")
}

func (s *Service) crocoCaveLogs(topic, payload string) {
	log.Printf("DEBUG [%s] \t %s\r\n", topic, payload)
	// croco/cave/temperature
	// croco/cave/targetTemperature
	// croco/cave/heater
	// croco/cave/light

	topicParts := strings.Split(topic, "/")
	switch topicParts[2] {
	case "temperature":
		s.Write(store.Data{Module: "cave", Topic: "temp", Value: payload})
	case "targetTemperature":
		s.Write(store.Data{Module: "cave", Topic: "targetTemp", Value: payload})
	case "heater":
		s.Write(store.Data{Module: "cave", Topic: "heater", Value: payload})
	case "light":
		s.Write(store.Data{Module: "cave", Topic: "light", Value: payload})
	}
}

func (s *Service) crocoCaveRaw(topic, payload string) {
	// ESP32 probes raw logs
	// ESP32-A473F53A7D80/p/ds18b20/#
	// ESP32-A473F53A7D80/p/ds18b20/1	23.75
	// ESP32-A473F53A7D80/p/ds18b20/2	24.00

	topicParts := strings.Split(topic, "/")
	if len(topicParts) == 4 && topicParts[3] > "0" {
		val, err := strconv.ParseFloat(payload, 64)
		if err != nil {
			log.Printf("ERROR [%s] \t %s \t %e \r\n", topic, payload, err)
			return
		}
		// Ignore invalid values
		if val < 0 || val > 100 {
			return
		}

		log.Printf("DEBUG [%s] \t %s\r\n", topic, payload)
		s.Write(store.Data{Module: "probes", Topic: "ds18b20" + "/" + topicParts[3], Value: payload})
	}
}
