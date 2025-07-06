package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Admin struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func MustNewAdmin(rabbitMqConf RabbitConf) *Admin {
	var admin Admin
	conn, err := amqp.Dial(GetRabbitURL(rabbitMqConf))
	if err != nil {
		log.Fatalf("failed to connect rabbitmq, error: %v", err)
	}

	admin.conn = conn
	channel, err := admin.conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel, error: %v", err)
	}

	admin.channel = channel
	return &admin
}

func (q *Admin) DeclareExchange(conf ExchangeConf, args amqp.Table) error {
	return q.channel.ExchangeDeclare(
		conf.ExchangeName,
		conf.Type,
		conf.Durable,
		conf.AutoDelete,
		conf.Internal,
		conf.NoWait,
		args,
	)
}

func (q *Admin) DeclareQueue(conf QueueConf, args amqp.Table) error {
	_, err := q.channel.QueueDeclare(
		conf.Name,
		conf.Durable,
		conf.AutoDelete,
		conf.Exclusive,
		conf.NoWait,
		args,
	)

	return err
}

func (q *Admin) Bind(queueName string, routeKey string, exchange string, notWait bool, args amqp.Table) error {
	return q.channel.QueueBind(
		queueName,
		routeKey,
		exchange,
		notWait,
		args,
	)
}

func (q *Admin) Init(name string, t string) {
	exchangeConf := ExchangeConf{
		ExchangeName: name,
		Type:         t,
		Durable:      true,
		AutoDelete:   false,
		Internal:     false,
		NoWait:       false,
	}
	err := q.DeclareExchange(exchangeConf, nil)
	if err != nil {
		log.Fatal(err)
	}
	queueConf := QueueConf{
		Name:       name,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	}
	err = q.DeclareQueue(queueConf, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = q.Bind(name, name, name, false, nil)
	if err != nil {
		log.Fatal(err)
	}
}
