package main

import (
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	network = "tcp"
	address = "localhost:9092"
)

func main() {
	if len(os.Args) < 2 {
		panic("You need to enter a topic you want to delete!")
	}
	topic := os.Args[1]

	conn, err := kafka.Dial(network, address)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	err = conn.DeleteTopics(topic)
	if err != nil {
		log.Fatal("failed to delete old topics")
	}
	fmt.Println(fmt.Sprintf("Successfully deleted an old topic: %s!", topic))
}
