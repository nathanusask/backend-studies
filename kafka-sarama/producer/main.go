package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	p sarama.AsyncProducer
}

func NewProducer(broker string) (*Producer, error) {
	producer, err := sarama.NewAsyncProducer([]string{broker}, sarama.NewConfig())
	if err != nil {
		return nil, err
	}
	return &Producer{
		p: producer,
	}, nil
}

type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (p *Producer) StartProduce(topic string) {
	for i := 0; i < 100; i++ {
		msg := Message{strconv.FormatInt(time.Now().Unix(), 10), "Message: " + strconv.Itoa(i)}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		select {
		case p.p.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(msgBytes),
		}:
			fmt.Println("Successfully sent message: ", msg)
		case err := <-p.p.Errors():
			fmt.Printf("Failed to send message to kafka, err: %s, msg: %s\n", err, msgBytes)
		}

	}
}

func (p *Producer) Close() error {
	if p != nil {
		return p.p.Close()
	}
	return nil
}

func main() {
	broker := "localhost:9092"
	topic := "an-awesome-topic"

	producer, err := NewProducer(broker)
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	go producer.StartProduce(topic)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println("received signal", <-c)
}
