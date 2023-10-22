package queue

import (
	"context"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/parMaster/logserver/app/config"
)

// As soon as we have another message queue, we can:
// 1. Create an interface for it
// 2. Create a factory function that returns an instance of the interface
// 3. Replace all references to mqtt.Client with the interface
// 4. Replace the factory function with a switch statement that returns the appropriate instance
// 5. Replace the mqtt.NewClient call with a call to the factory function

// Message represents a message from mqtt queue in strings
type Message struct {
	Topic   string
	Payload string
}

type Subscription struct {
	Topic    string                      // topic to subscribe to (e.g. "croco/cave/#")
	Handler  func(topic, payload string) // handler function to process message
	Messages chan Message                // channel to consume messages from
}

type Client struct {
	mqtt.Client
	subs []Subscription
}

// NewClient creates a new mqtt client and subscribes to the given topics
func NewClient(ctx context.Context, config config.Config, subs ...Subscription) (*Client, error) {

	opts := mqtt.NewClientOptions().AddBroker(config.MqBrokerURL)
	opts.SetUsername(config.MqUser)
	opts.SetPassword(config.MqPassword)
	opts.SetClientID(config.MqClientId)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Second)
	opts.SetResumeSubs(true)

	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("[ERROR] Connection to mqtt broker lost: %s", err)
	})

	opts.SetReconnectingHandler(func(c mqtt.Client, opts *mqtt.ClientOptions) {
		log.Printf("[INFO] Reconnecting to mqtt broker %s as %s", config.MqBrokerURL, config.MqClientId)
	})

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("[INFO] Connected to mqtt broker %s as %s", config.MqBrokerURL, config.MqClientId)

		for i, sub := range subs {
			if token := c.Subscribe(sub.Topic, 0, func(c mqtt.Client, m mqtt.Message) {
				// either Handler or Messages channel can be used
				// if sub.Handler != nil {
				// 	sub.Handler(m.Topic(), string(m.Payload()))
				// }
				if sub.Messages != nil {
					subs[i].Messages <- Message{Topic: m.Topic(), Payload: string(m.Payload())}
				}
			}); token.Wait() && token.Error() != nil {
				log.Printf("[ERROR] failed to subscribe to topic %s: %s", sub.Topic, token.Error())
			}
		}

	})

	mq := mqtt.NewClient(opts)
	if token := mq.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[ERROR] failed to connect to mqtt: %s", token.Error())
		return nil, token.Error()
	}

	go func() {
		<-ctx.Done()
		log.Printf("[INFO] Terminating mqtt client")
		mq.Disconnect(250)
	}()

	Client := Client{
		mq,
		subs,
	}

	return &Client, nil
}
