package main

import (
	"fmt"
	"log"

	stan "github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "subID")
	if err != nil {
		log.Fatal(err)
	}

	defer sc.Close()

	sub, err := sc.Subscribe("foo", func(m *stan.Msg) {
		fmt.Printf("Received message: %s", string(m.Data))
	}, stan.DeliverAllAvailable())
	if err != nil {
		log.Fatal(err)
	}

	defer sub.Unsubscribe()
}
