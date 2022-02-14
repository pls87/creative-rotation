package stats

import (
	"encoding/json"
	"fmt"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/streadway/amqp"
)

type Producer interface {
	Client
	Produce(routingKey string, message Event) error
}

type RabbitProducer struct {
	RabbitClient
}

func (ap *RabbitProducer) Produce(routingKey string, message Event) (err error) {
	var body []byte
	body, err = json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error while publishing: couldn't marshal message: %w", err)
	}

	var ch *amqp.Channel
	if ch, err = ap.openChannel(); err != nil {
		return fmt.Errorf("error while publishing: %w", err)
	}
	defer ch.Close()

	if err = ch.Publish(Exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("error while publishing: %w", err)
	}

	return err
}

func NewProducer(c config.QueueConf) Producer {
	return &RabbitProducer{
		RabbitClient{
			cfg: c,
		},
	}
}
