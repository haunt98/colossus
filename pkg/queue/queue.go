package queue

import (
	"encoding/json"
	"fmt"

	"github.com/mailru/easyjson"

	"github.com/streadway/amqp"
)

//TODO wrap error

type Queue struct {
	amqpChannel *amqp.Channel
	amqpQueue   amqp.Queue
}

func NewQueue(channel *amqp.Channel, name string) (*Queue, error) {
	queue, err := channel.QueueDeclare(name,
		false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Queue{
		amqpChannel: channel,
		amqpQueue:   queue,
	}, nil
}

const jsonContentType = "application/json"

func (q *Queue) PublishJSON(value interface{}) error {
	var data []byte
	var err error
	if easyValue, ok := value.(easyjson.Marshaler); ok {
		data, err = easyjson.Marshal(easyValue)
		if err != nil {
			return fmt.Errorf("json failed to marshal %v: %w", value, err)
		}
	} else {
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("json failed to marshal %v: %w", value, err)
		}
	}

	if err := q.amqpChannel.Publish("", q.amqpQueue.Name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  jsonContentType,
			Body:         data,
		}); err != nil {
		return err
	}

	return nil
}

func (q *Queue) Consume() (<-chan amqp.Delivery, error) {
	deliveries, err := q.amqpChannel.Consume(q.amqpQueue.Name,
		"", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}
