package rabbitmq

import (
	"context"
	"log"

	"github.com/streadway/amqp"
)

type RabbitCFG struct {
	URI          string
	Exchange     string
	ExchangeType string
	Queue        string
	BindingKey   string
	ConsumerTag  string
}

type RabbitQueue struct {
	rabCfg     RabbitCFG
	conn       *amqp.Connection
	channel    *amqp.Channel
	deliveries <-chan amqp.Delivery
	done       chan error
	ctx        context.Context
}

func (r *RabbitQueue) SendMess(myMes []byte /*string*/) error {
	err := r.channel.Publish(
		r.rabCfg.Exchange,   // publish to an exchange
		r.rabCfg.BindingKey, // routing to 0 or more queues
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            myMes,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	)
	return err
}

func (r *RabbitQueue) Shutdown() {
	r.channel.Cancel(r.rabCfg.ConsumerTag, true)
	r.conn.Close()
}

func (r *RabbitQueue) Handle() {
	for d := range r.deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
	r.done <- nil
}

func CreateQueue(q RabbitCFG, ctx context.Context) (*RabbitQueue, error) {
	c := &RabbitQueue{
		rabCfg:  q,
		conn:    nil,
		channel: nil,
		done:    make(chan error),
		ctx:     ctx,
	}
	var err error
	c.conn, err = amqp.Dial(q.URI)
	if err != nil {
		return nil, err
	}
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	if err = c.channel.ExchangeDeclare(
		q.Exchange,     // name of the exchange
		q.ExchangeType, // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return nil, err
	}

	log.Printf("declared Exchange, declaring Queue %q", q.Queue)
	_, err = c.channel.QueueDeclare(
		q.Queue, // name of the queue
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // noWait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	if err = c.channel.QueueBind(
		q.Queue,      // name of the queue
		q.BindingKey, // bindingKey
		q.Exchange,   // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, err
	}
	return c, err
}
