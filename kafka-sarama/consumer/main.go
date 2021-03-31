package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathanusask/backend-studies/kafka-sarama/consumer/consumer"
)

func main() {
	var broker = "localhost:9092"
	var topic = "an-awesome-topic"

	//multiBatchConsumer, err := consumer.StartMultiBatchConsumer(broker, topic)
	//if err != nil {
	//	panic(err)
	//}
	//defer multiBatchConsumer.Close()

	syncConsumer, err := consumer.StartSyncConsumer(broker, topic)
	if err != nil {
		panic(err)
	}
	defer syncConsumer.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println("received signal", <-c)
}
