package messageAdapter

type Event struct {
    Command      string
	ReplyChannel ChannelID
	Body         string
}


type ChannelID string

type QueueID string

type Queue interface {
    create() QueueID
	Receive() (Message, bool)
	Send(Message)
}

type Message struct {
    Command      string
	ReplyChannel ChannelID
	Body         string
}

type Channel interface {
	Receive() (Message, bool)
	Send(Message)
}

type ChannelMap map[ChannelID]Channel
