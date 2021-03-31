package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	network = "tcp" // may need to update upon scaling up
	address = "localhost:9092"

	topic     = "an-awesome-topic" // subject to change
	partition = 0                  // subject to change
	//numPartitions = 5
)

func main() {
	// to produce messages
	ctx := context.Background()

	//conn, err := kafka.DialLeader(ctx, network, address, topic, partition)
	//if err != nil {
	//	log.Fatal("Failed to dial ", network, address, err)
	//}
	//defer func(){
	//	if err := conn.Close(); err != nil {
	//		log.Fatal("failed to close conn: ", err)
	//	}
	//}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println("Entered writer goroutine...")
		defer wg.Done()
		//conn.SetWriteDeadline(time.Now().Add(30 * time.Second))

		w := kafka.Writer{
			Addr:        kafka.TCP(address),
			Topic:       topic,
			Balancer:    &kafka.LeastBytes{},
			Async:       true,
			Compression: kafka.Lz4,
		}
		i := 0
		for {
			msgStr := "Counting " + strconv.Itoa(i)

			err := w.WriteMessages(ctx, kafka.Message{
				Value: []byte(msgStr),
			})
			if err != nil {
				log.Fatal("failed to write messages:", err)
			}

			// Wrap up message and send
			//msg := kafka.Message{Value: []byte("Counting " + strconv.Itoa(i))}
			//err := w.WriteMessages(ctx, msg)
			//if err != nil {
			//	log.Fatal("Failed to write message: ", err.Error())
			//	break
			//}

			fmt.Println("Successfully write message: ", msgStr)
			time.Sleep(time.Second)
			i++
		}

		if err := w.Close(); err != nil {
			log.Fatal("failed to close writer: ", err)
		}
	}()

	wg.Wait()
}
