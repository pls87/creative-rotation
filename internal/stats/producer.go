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
	ch *amqp.Channel
}

func (ap *RabbitProducer) Init() (err error) {
	if err = ap.RabbitClient.Init(); err != nil {
		return err
	}
	ap.ch, err = ap.openChannel()
	return err
}

func (ap *RabbitProducer) Produce(routingKey string, message Event) (err error) {
	var body []byte
	body, err = json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error while publishing: couldn't marshal message: %w", err)
	}

	if err = ap.ch.Publish(Exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}); err != nil {
		return fmt.Errorf("error while publishing: %w", err)
	}

	return err
}

func NewProducer(c config.QueueConf) Producer {
	return &RabbitProducer{
		RabbitClient: RabbitClient{
			cfg: c,
		},
	}
}
