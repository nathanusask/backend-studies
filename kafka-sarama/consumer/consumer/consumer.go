package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GroupHandler interface {
	sarama.ConsumerGroupHandler
	WaitReady()
	Reset()
}

type Group struct {
	cg sarama.ConsumerGroup
}

func NewConsumerGroup(broker string, topics []string, group string, handler GroupHandler) (*Group, error) {
	ctx := context.Background()
	cfg := sarama.NewConfig()
	cfg.Version = sarama.MaxVersion
	//cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	client, err := sarama.NewConsumerGroup([]string{broker}, group, cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := client.Consume(ctx, topics, handler)
			if err != nil {
				if err == sarama.ErrClosedConsumerGroup {
					break
				} else {
					panic(err)
				}
			}
			if ctx.Err() != nil {
				return
			}
			handler.Reset()
		}
	}()

	handler.WaitReady() // Await till the consumer has been set up

	return &Group{
		cg: client,
	}, nil
}

func (c *Group) Close() error {
	return c.cg.Close()
}

type SessionMessage struct {
	Session sarama.ConsumerGroupSession
	Message *sarama.ConsumerMessage
}

func decodeMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func StartSyncConsumer(broker, topic string) (*Group, error) {
	handler := NewSyncConsumerGroupHandler(func(data []byte) error {
		if msg, err := decodeMessage(data); err != nil {
			fmt.Println(fmt.Sprintf("Successfully read message %v", msg))
			return err
		}
		return nil
	})
	consumer, err := NewConsumerGroup(broker, []string{topic}, "sync-consumer-"+fmt.Sprintf("%d", time.Now().Unix()), handler)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func StartMultiAsyncConsumer(broker, topic string) (*Group, error) {
	var bufChan = make(chan *SessionMessage, 1000)
	for i := 0; i < 8; i++ {
		go func() {
			for message := range bufChan {
				if msg, err := decodeMessage(message.Message.Value); err == nil {
					message.Session.MarkMessage(message.Message, "")
					fmt.Println("Successfully read message: ", msg)
				} else {
					fmt.Println("Failed to read message with error: ", err)
				}
			}
		}()
	}
	handler := NewMultiAsyncConsumerGroupHandler(&MultiAsyncConsumerConfig{
		BufChan: bufChan,
	})
	consumer, err := NewConsumerGroup(broker, []string{topic}, "multi-async-consumer-"+fmt.Sprintf("%d", time.Now().Unix()), handler)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func StartMultiBatchConsumer(broker, topic string) (*Group, error) {
	var bufChan = make(chan batchMessages, 100)
	for i := 0; i < 8; i++ {
		go func() {
			for messages := range bufChan {
				for _, message := range messages {
					if msg, err := decodeMessage(message.Message.Value); err == nil {
						//message.Session.Commit()
						message.Session.MarkMessage(message.Message, "")
						fmt.Printf("Successfully read from partition %d at offset %d, with message: %v\n",
							message.Message.Partition,
							message.Message.Offset,
							msg,
						)
					}
				}
			}
		}()
	}
	handler := NewMultiBatchConsumerGroupHandler(&MultiBatchConsumerConfig{
		MaxBufSize: 100,
		BufChan:    bufChan,
	})
	consumer, err := NewConsumerGroup(broker, []string{topic}, "multi-batch-consumer-"+fmt.Sprintf("%d", time.Now().Unix()), handler)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
