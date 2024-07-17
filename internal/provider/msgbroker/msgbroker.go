package msgbroker

type Topic string

const (
	TopicTaskFinished Topic = "task.finished"
)

type Handler func(message []byte)

type MessageBroker interface {
	Publish(topic Topic, message []byte) error
	Subscribe(topic Topic, handler Handler) error
}
