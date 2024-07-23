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
	conn      *kafka.Conn
	addr      string
	topic     string
	partition int
}

func (n *Notify) Connect(ctx context.Context) error {
	var err error

	n.conn, err = kafka.DialLeader(ctx, "tcp", n.addr, n.topic, n.partition)
	if err != nil {
		return err
	}

	return nil
}

func (n *Notify) Notify(_ context.Context, message notify.Message) error {
	data, err := easyjson.Marshal(message)
	if err != nil {
		return err
	}

	if _, err := n.conn.WriteMessages(kafka.Message{
		Value: data,
	}); err != nil {
		return err
	}

	return nil
}

func (n *Notify) Close() error {
	return n.conn.Close()
}
