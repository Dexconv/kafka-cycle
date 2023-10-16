package common

import (
	"log"
	"math"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Queue struct {
	Mu     sync.Mutex
	Buffer []byte
}

func MakeQueue() *Queue {
	return &Queue{
		Mu: sync.Mutex{},
		//1
		Buffer: make([]byte, 0, 3075),
	}
}

func Dequeue(q []byte) (res []kafka.Message) {
	size := int(math.Ceil(float64(len(q)) / 4))
	log.Printf("queue length %d devided to 4 parts with max length of %d ", len(q), size)
	var j int
	for i := 0; i < len(q); i += size {
		j += size
		if j > len(q) {
			j = len(q)
		}
		res = append(res, kafka.Message{Value: q[i:j]})
	}
	return
}
