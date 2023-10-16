package main

import (
	"context"
	"log"

	"github.com/dexconv/kafka-cycle/common"
	"github.com/dexconv/kafka-cycle/kafka"
	"golang.org/x/sync/errgroup"
)

var (
	queue = common.MakeQueue()
)

func main() {
	writer := kafka.NewKafkaWriter()

	ctx := context.Background()
	sent := make(chan []byte, 3075)
	recieved := make(chan []byte, 3075)
	cooldown := make(chan bool, 1)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return common.Verify(ctx, sent, recieved, cooldown)
	})

	g.Go(func() error {
		return writer.WriteMessages(ctx, sent, queue, cooldown)
	})

	g.Go(func() error {
		return kafka.ReadMessages(ctx, recieved)
	})

	err := g.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
