package stats

import (
	"fmt"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/streadway/amqp"
)

const (
	Exchange        = "stats"
	ImpressionQueue = "impressions"
	ConversionQueue = "conversions"
	ImpressionKey   = "new_impression"
	ConversionKey   = "new_conversion"
)

type Client interface {
	Init() error
	Dispose() error
}

type RabbitClient struct {
	conn *amqp.Connection
	cfg  config.QueueConf
}

func (sc *RabbitClient) Init() (err error) {
	sc.conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		sc.cfg.User, sc.cfg.Password, sc.cfg.Host, sc.cfg.Port))
	if err != nil {
		return fmt.Errorf("couldn't connect to queue: %w", err)
	}

	var ch *amqp.Channel
	ch, err = sc.conn.Channel()
	if err != nil {
		return fmt.Errorf("couldn't open channel: %w", err)
	}
	defer ch.Close()

	if err = ch.ExchangeDeclare(Exchange, "direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create exchange %s: %w", Exchange, err)
	}

	if _, err = ch.QueueDeclare(ImpressionQueue,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create queue %s: %w", ImpressionQueue, err)
	}

	if err = ch.QueueBind(ImpressionQueue, ImpressionKey, Exchange, false, nil); err != nil {
		return fmt.Errorf("error binding queue='%s' to exchange='%s' with routing key='%s': %w",
			ImpressionQueue, Exchange, ImpressionKey, err)
	}

	if _, err = ch.QueueDeclare(ConversionQueue,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("couldn't create queue %s: %w", ConversionQueue, err)
	}

	if err = ch.QueueBind(ConversionQueue, ConversionKey, Exchange, false, nil); err != nil {
		return fmt.Errorf("error binding queue='%s' to exchange='%s' with routing key='%s': %w",
			ConversionQueue, Exchange, ConversionKey, err)
	}

	return nil
}

func (sc *RabbitClient) Dispose() (err error) {
	if sc.conn != nil {
		return sc.conn.Close()
	}
	return nil
}
