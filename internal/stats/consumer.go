package stats

import (
	"encoding/json"
	"fmt"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/streadway/amqp"
)

type Consumer interface {
	Client
	Consume(tag, queue string) (messages chan Event, errors chan error, err error)
}

type RabbitConsumer struct {
	RabbitClient
}

func (nc *RabbitConsumer) openChannel() (ch *amqp.Channel, err error) {
	ch, err = nc.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("couldn't open channel: %w", err)
	}

	return ch, err
}

func (nc *RabbitConsumer) Consume(tag, queue string) (messages chan Event, errors chan error, err error) {
	var ch *amqp.Channel
	if ch, err = nc.openChannel(); err != nil {
		return nil, nil, fmt.Errorf("error while consuming messages: %w", err)
	}

	var deliveries <-chan amqp.Delivery
	deliveries, err = ch.Consume(
		queue, // name
		tag,   // consumerTag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error while consuming messages: %w", err)
	}

	messages = make(chan Event)
	errors = make(chan error)

	go func() {
		defer func() {
			close(messages)
			_ = ch.Close()
			close(errors)
		}()
		var e error
		for d := range deliveries {
			if e = d.Ack(false); e != nil {
				errors <- fmt.Errorf("message couldn't be acknowledged: %w", e)
				continue
			}
			var msg Event
			if e = json.Unmarshal(d.Body, &msg); e != nil {
				errors <- fmt.Errorf("message couldn't be parsed: %w", e)
				continue
			}
			messages <- msg
		}
	}()

	return messages, errors, nil
}

func NewConsumer(c config.QueueConf) Consumer {
	return &RabbitConsumer{
		RabbitClient{
			cfg: c,
		},
	}
}
