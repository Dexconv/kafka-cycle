package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

type Reader struct {
	Reader *kafkago.Reader
}

var (
	reader1 = NewKafkaReader(0)
	reader2 = NewKafkaReader(1)
	reader3 = NewKafkaReader(2)
	reader4 = NewKafkaReader(3)
)

func NewKafkaReader(partition int) *Reader {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     "test",
		Partition: partition,
	})

	return &Reader{
		Reader: reader,
	}
}

func (k *Reader) FetchMessage(ctx context.Context) (kafkago.Message, error) {

	message, err := k.Reader.FetchMessage(ctx)
	if err != nil {
		return kafkago.Message{}, err
	}
	select {
	case <-ctx.Done():
		return kafkago.Message{}, ctx.Err()
	default:
		// log.Printf("message fetched from partition: %v \n", k.Reader.Config().Partition)
		return message, nil
	}

}

func ReadMessages(ctx context.Context, recieved chan<- []byte) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m1, err := reader1.FetchMessage(ctx)
			if err != nil {
				return err
			}
			m2, err := reader2.FetchMessage(ctx)
			if err != nil {
				return err
			}
			m3, err := reader3.FetchMessage(ctx)
			if err != nil {
				return err
			}
			m4, err := reader4.FetchMessage(ctx)
			if err != nil {
				return err
			}
			x1 := append(m1.Value, m2.Value...)
			x2 := append(m3.Value, m4.Value...)
			recieved <- append(x1, x2...)
			if err != nil {
				return err
			}
		}
	}
}
