package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	network = "tcp"
	address = "localhost:9092"

	topic = "an-awesome-topic" // subject to change
	//partition = 0                  // subject to change

	consumerGroupID = "consumer-group-id"
)

func main() {
	// to consume messages
	ctx := context.Background()

	var cgID string // consumer group ID
	if len(os.Args) < 2 {
		cgID = consumerGroupID
	} else {
		cgID = os.Args[1]
	}

	//conn, err := kafka.DialLeader(ctx, network, address, topic, partition)
	//if err != nil {
	//	log.Fatal("failed to dial leader:", err)
	//}
	//
	//defer func() {
	//	if err := conn.Close(); err != nil {
	//		log.Fatal("failed to close connection:", err)
	//	}
	//}()

	//fmt.Println("Entered reader goroutine...")
	//conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{address},
		GroupID:  cgID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		// by default, the StartOffset field is set kafka.FirstOffset
		// meaning the reader reads from the beginning
		//StartOffset: [kafka.FirstOffset | kafka.LastOffset],
	})

	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}
	}()

	// consumer group method
	//for {
	//	m, err := r.ReadMessage(ctx)
	//	if err != nil {
	//		break
	//	}
	//	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	//}

	// commit-per-message approach
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		fmt.Println(fmt.Sprintf("message at topic: %s\tpartition: %d\t/offset: %d\t: %s = %s", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value)))
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}

	//first, last, err := conn.ReadOffsets()
	//if err != nil {
	//	log.Fatal("failed to read offsets")
	//}
	//fmt.Println("First:", first)
	//fmt.Println("Last:", last)

	//batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
	//
	//b := make([]byte, 10e3) // 10KB max per message
	//for {
	//	_, err := batch.Read(b)
	//	if err != nil {
	//		break
	//	}
	//	fmt.Println(string(b))
	//	batch.
	//}
	//
	//if err := batch.Close(); err != nil {
	//	log.Fatal("failed to close batch:", err)
	//}
}
