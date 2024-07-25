package kafkanotify

import (
	"context"
	"net"
	"strconv"

	"github.com/mailru/easyjson"
	"github.com/segmentio/kafka-go"
	"github.com/trb1maker/gobannercircle/internal/notify"
)

func NewKafkaNotify(host string, port int, topic string, partition int) *Notify {
	return &Notify{
		addr:      net.JoinHostPort(host, strconv.Itoa(port)),
		topic:     topic,
		partition: partition,
	}
}

type Notify struct {
	wr        *kafka.Writer
	addr      string
	topic     string
	partition int
}

func (n *Notify) Connect() {
	n.wr = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{n.addr},
		Topic:    n.topic,
		Balancer: &kafka.LeastBytes{},
	})
}

func (n *Notify) Notify(ctx context.Context, message notify.Message) error {
	data, err := easyjson.Marshal(message)
	if err != nil {
		return err
	}

	return n.wr.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}

func (n *Notify) Close() error {
	return n.wr.Close()
}
