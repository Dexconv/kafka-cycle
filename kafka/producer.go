package kafka

import (
	"context"
	"crypto/rand"
	"log"
	mathrand "math/rand"

	"github.com/dexconv/kafka-cycle/common"
	kafkago "github.com/segmentio/kafka-go"
)

type Writer struct {
	Writer *kafkago.Writer
}

func NewKafkaWriter() *Writer {
	writer := &kafkago.Writer{
		Addr:     kafkago.TCP("localhost:29092"),
		Topic:    "test",
		Balancer: &kafkago.RoundRobin{},
	}
	return &Writer{
		Writer: writer,
	}
}

func (k *Writer) WriteMessages(
	ctx context.Context,
	sent chan []byte,
	q *common.Queue,
	cooldown chan bool,
) error {
	for {
		nDataLen := mathrand.Intn(512-320) + 320
		// nDataLen := 512
		buf := make([]byte, nDataLen)
		_, err := rand.Read(buf)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			q.Mu.Lock()
			switch {
			case len(q.Buffer)+len(buf) > cap(q.Buffer):
				log.Println("reached cap, emptying queue...")
				res := common.Dequeue(q.Buffer)
				err = k.Writer.WriteMessages(ctx, res[0], res[1], res[2], res[3])
				if err != nil {
					return err
				}
				sent <- q.Buffer
				q.Buffer = make([]byte, 0, 3075)
				// adding new generated data to empty queue
				q.Buffer = append(q.Buffer, buf...)
			case (len(q.Buffer)+len(buf))%1024 == 0:
				log.Println("reached condition, emptying queue...")
				q.Buffer = append(q.Buffer, buf...)
				res := common.Dequeue(q.Buffer)
				err = k.Writer.WriteMessages(ctx, res[0], res[1], res[2], res[3])
				if err != nil {
					return err
				}
				sent <- q.Buffer
				q.Buffer = make([]byte, 0, 3075)
			default:
				q.Buffer = append(q.Buffer, buf...)
				cooldown <- true
			}
			q.Mu.Unlock()
		}
		// entring cooldown
		<-cooldown
	}
}
