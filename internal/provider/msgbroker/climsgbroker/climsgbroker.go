package climsgbroker

import (
	"fmt"
	"sync"

	"github.com/danielmesquitta/tasks-api/internal/provider/msgbroker"
)

type CLIMessageBroker struct {
	subscribers map[msgbroker.Topic][]msgbroker.Handler
	mutex       sync.Mutex
}

func NewCLIMessageBroker() *CLIMessageBroker {
	return &CLIMessageBroker{
		subscribers: make(map[msgbroker.Topic][]msgbroker.Handler),
	}
}

func (cmb *CLIMessageBroker) Publish(
	topic msgbroker.Topic,
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
	topic msgbroker.Topic,
	handler msgbroker.Handler,
) error {
	cmb.mutex.Lock()
	defer cmb.mutex.Unlock()

	cmb.subscribers[topic] = append(cmb.subscribers[topic], handler)
	return nil
}
