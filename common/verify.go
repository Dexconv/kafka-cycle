package common

import (
	"bytes"
	"context"
	"fmt"
	"log"
)

func Verify(ctx context.Context, sent, recieved chan []byte, cooldown chan bool) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			a := <-sent
			b := <-recieved
			cooldown <- true
			if bytes.Equal(a, b) {
				log.Println("message is verified")
				log.Printf("sent: %x...%x , recieved: %x...%x", a[:5], a[len(a)-5:], b[:5], b[len(b)-5:])
			} else {
				return fmt.Errorf("message could not be verified")
			}
		}
	}
}
