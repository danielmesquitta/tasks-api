package clibroker

import (
	"fmt"
	"sync"

	"github.com/danielmesquitta/tasks-api/internal/provider/broker"
)

type CLIMessageBroker struct {
	subscribers map[broker.Topic][]broker.Handler
	mutex       sync.Mutex
}

func NewCLIMessageBroker() *CLIMessageBroker {
	return &CLIMessageBroker{
		subscribers: make(map[broker.Topic][]broker.Handler),
	}
}

func (cmb *CLIMessageBroker) Publish(
	topic broker.Topic,
	message []byte,
) error {
	cmb.mutex.Lock()
	defer cmb.mutex.Unlock()

	if handlers, found := cmb.subscribers[topic]; found {
		for _, handler := range handlers {
			go handler(message)
		}
	}
	fmt.Printf("Published message to topic %s: %s\n", topic, message)
	return nil
}

func (cmb *CLIMessageBroker) Subscribe(
	topic broker.Topic,
	handler broker.Handler,
) error {
	cmb.mutex.Lock()
	defer cmb.mutex.Unlock()

	cmb.subscribers[topic] = append(cmb.subscribers[topic], handler)
	return nil
}

var _ broker.MessageBroker = (*CLIMessageBroker)(nil)
