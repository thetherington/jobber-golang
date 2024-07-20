package rabbitmq

import (
	"time"

	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer func() error

type RabbitQueueAdapter struct {
	URL string

	errorChannel chan *amqp.Error
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	Closed       bool

	Consumers []Consumer
}

func NewRabbitMQAdapter(url string) *RabbitQueueAdapter {
	q := new(RabbitQueueAdapter)

	q.URL = url
	q.Consumers = make([]Consumer, 0)

	q.Connect()
	go q.Reconnector()

	return q
}

func (q *RabbitQueueAdapter) Reconnector() {
	for {
		err := <-q.errorChannel
		if !q.Closed {
			slog.With("error", err).Error("RabbitMQ reconnecting after connection closed")

			q.Connect()

			for _, fn := range q.Consumers {
				fn()
			}
		}
	}
}

func (q *RabbitQueueAdapter) Connect() {
	for {
		slog.Info("Connecting to RabbitMQ", "address", q.URL)

		conn, err := amqp.Dial(q.URL)
		if err == nil {
			q.Connection = conn
			q.errorChannel = make(chan *amqp.Error)
			q.Connection.NotifyClose(q.errorChannel)

			slog.Info("RabbitMQ connection established!")

			q.openChannel()

			return
		}

		slog.With("error", err).Warn("Connection to RabbitMQ failed. Retrying in 1 sec.")
		time.Sleep(1000 * time.Millisecond)
	}
}

func (q *RabbitQueueAdapter) openChannel() {
	channel, err := q.Connection.Channel()
	if err != nil {
		slog.With("error", err).Error("RabbitMQ Opening channel failed")
		return
	}

	q.Channel = channel
}

func (q *RabbitQueueAdapter) Close() {
	slog.Info("RabbitMQ closing connection")

	q.Closed = true
	q.Channel.Close()
	q.Connection.Close()
}

func (q *RabbitQueueAdapter) AddConsumer(fn Consumer) error {
	q.Consumers = append(q.Consumers, fn)

	return fn()
}
