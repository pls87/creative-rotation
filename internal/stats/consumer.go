package stats

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
)

type Consumer interface {
	Init(cfg config.QueueConf) error
	Dispose() error
	Consume(topic string) (chan Event, error)
}

type KafkaConsumer struct {
	c    *kafka.Consumer
	logg *logger.Logger
	cfg  config.QueueConf
}

func (kc *KafkaConsumer) Init() (err error) {
	kc.c, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": net.JoinHostPort(kc.cfg.Host, strconv.Itoa(kc.cfg.Port)),
	})
	if err != nil {
		return fmt.Errorf("couldn't initialize consumer: %w", err)
	}
	return nil
}

func (kc *KafkaConsumer) Dispose() error {
	return kc.c.Close()
}

func (kc *KafkaConsumer) Consume(topic string) (events chan Event, err error) {
	err = kc.c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't subscribe topic %s: %w", topic, err)
	}
	events = make(chan Event)
	var ev Event
	for e := range kc.c.Events() {
		msg := e.String()
		err = json.Unmarshal([]byte(msg), &ev)
		if err != nil {
			kc.logg.Errorf("couldn't unmarshal the message: %s: %s", msg, err)
			continue
		}
		events <- ev
	}

	return events, nil
}
