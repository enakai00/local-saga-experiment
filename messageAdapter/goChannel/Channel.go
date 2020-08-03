package goChannel

import (
	"github.com/enakai00/local-saga-experiment/messageAdapter"
)

type GoChannel chan messageAdapter.Message // Instance of Channel

func (ch GoChannel) Send(mes messageAdapter.Message) {
	ch <- mes
}

func (ch GoChannel) Receive() (messageAdapter.Message, bool) {
	mes, ok := <-ch
	return mes, ok
}

func PrepareChannels(channelIDs []messageAdapter.ChannelID) messageAdapter.ChannelMap {
    channels := make(messageAdapter.ChannelMap)
    for _, id := range channelIDs {
        channel := GoChannel(make(chan messageAdapter.Message, 10))
        channels[id] = channel
    }
    return channels
}
