package server

import (
	"log"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/parMaster/logserver/app/store"
)

func (l *LogServer) newMqClient() (mqtt.Client, error) {

	opts := mqtt.NewClientOptions().AddBroker(l.config.MqBrokerURL)
	opts.SetUsername(l.config.MqUser)
	opts.SetPassword(l.config.MqPassword)
	opts.SetClientID(l.config.MqClientId)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Second)
	opts.SetResumeSubs(true)

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("[INFO] Connected to mqtt broker %s as %s", l.config.MqBrokerURL, l.config.MqClientId)
		l.SubscribeAndHandle()
	})

	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("[ERROR] Connection to mqtt broker lost: %s", err)
	})

	opts.SetReconnectingHandler(func(c mqtt.Client, opts *mqtt.ClientOptions) {
		log.Printf("[INFO] Reconnecting to mqtt broker %s as %s", l.config.MqBrokerURL, l.config.MqClientId)
	})

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[ERROR] failed to connect to mqtt: %s", token.Error())
		return nil, token.Error()
	}

	go func() {
		<-l.ctx.Done()
		log.Printf("[INFO] Terminating mqtt client")
		c.Disconnect(250)
	}()

	return c, nil
}

func (l *LogServer) Subscribe(topic string, handlerFunc mqtt.MessageHandler) {
	if token := l.mq.Subscribe(topic, 0, handlerFunc); token.Wait() && token.Error() != nil {
		log.Printf("[ERROR] failed to subscribe to topic %s: %s", topic, token.Error())
	}
}

func (l *LogServer) SubscribeAndHandle() {
	// Croco cave logs
	log.Printf("[INFO] Subscribing to croco/cave/#")
	l.Subscribe("croco/cave/#", func(c mqtt.Client, m mqtt.Message) {
		log.Printf("DEBUG [%s] \t %s\r\n", m.Topic(), m.Payload())
		// croco/cave/temperature
		// croco/cave/targetTemperature
		// croco/cave/heater
		// croco/cave/light

		topicParts := strings.Split(m.Topic(), "/")
		switch topicParts[2] {
		case "temperature":
			l.store.Write(store.Data{Module: "cave", Topic: "temp", Value: string(m.Payload())})
		case "targetTemperature":
			l.store.Write(store.Data{Module: "cave", Topic: "targetTemp", Value: string(m.Payload())})
		case "heater":
			l.store.Write(store.Data{Module: "cave", Topic: "heater", Value: string(m.Payload())})
		case "light":
			l.store.Write(store.Data{Module: "cave", Topic: "light", Value: string(m.Payload())})
		}

	})

	// ESP32 probes raw logs
	l.Subscribe("ESP32-A473F53A7D80/p/ds18b20/#", func(c mqtt.Client, m mqtt.Message) {
		// ESP32-A473F53A7D80/p/ds18b20/1	23.75
		// ESP32-A473F53A7D80/p/ds18b20/2	24.00
		if !l.config.CollectRaw {
			return
		}

		topicParts := strings.Split(m.Topic(), "/")
		if len(topicParts) == 4 && topicParts[3] > "0" {
			val, err := strconv.ParseFloat(string(m.Payload()), 64)
			if err != nil {
				log.Printf("ERROR [%s] \t %s \t %e \r\n", m.Topic(), m.Payload(), err)
				return
			}
			// Ignore invalid values
			if val < 0 || val > 100 {
				return
			}

			log.Printf("DEBUG [%s] \t %s\r\n", m.Topic(), m.Payload())
			l.store.Write(store.Data{Module: "probes", Topic: "ds18b20" + "/" + topicParts[3], Value: string(m.Payload())})
		}
	})

}
