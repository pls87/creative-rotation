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

type Producer interface {
	Init() error
	Dispose()
	Produce(topic string, msg Event) error
}

type KafkaProducer struct {
	p    *kafka.Producer
	cfg  config.QueueConf
	logg *logger.Logger
}

func (kp *KafkaProducer) Init() (err error) {
	kp.p, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": net.JoinHostPort(kp.cfg.Host, strconv.Itoa(kp.cfg.Port)),
	})
	if err != nil {
		return fmt.Errorf("couldn't initialize producer: %w", err)
	}

	e := <-kp.p.Events()
	if er, ok := e.(kafka.Error); ok {
		return fmt.Errorf("couldn't initialize producer: %w", er)
	}

	go func() {
		for e := range kp.p.Events() {
			if ev, ok := e.(*kafka.Message); ok {
				if ev.TopicPartition.Error != nil {
					kp.logg.Errorf("kafka delivery failed: %v", ev.TopicPartition)
				}
			}
		}
	}()

	return nil
}

func (kp *KafkaProducer) Dispose() {
	kp.p.Close()
}

func (kp *KafkaProducer) Produce(topic string, msg Event) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal message %v: %w", msg, err)
	}
	err = kp.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          body,
	}, nil)

	if err != nil {
		return fmt.Errorf("couldn't produce the message %v: %w", msg, err)
	}

	return nil
}

func NewProducer(logger *logger.Logger, cfg config.QueueConf) Producer {
	return &KafkaProducer{
		cfg:  cfg,
		logg: logger,
	}
}
