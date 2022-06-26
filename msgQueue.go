package pop_shark

type msgQueue interface {
	map[string]chan interface{}
	Send(topic string, msg interface{})
	Receive(topic string) interface{}
}

type mq struct {
	queueMap   map[string]chan interface{}
	bufferSize int
}

func NewMsgQueue() *mq {
	return &mq{
		queueMap:   make(map[string]chan interface{}),
		bufferSize: 25,
	}
}
func (m mq) Receive(topic string) interface{} {
	if _, ok := m.queueMap[topic]; !ok {
		m.queueMap[topic] = make(chan interface{}, m.bufferSize)
	}

	return <-m.queueMap[topic]

}
func (m mq) Send(topic string, msg interface{}) {
	if _, ok := m.queueMap[topic]; !ok {
		m.queueMap[topic] = make(chan interface{}, m.bufferSize)
	}

	m.queueMap[topic] <- msg
}
