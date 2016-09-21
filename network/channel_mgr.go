package network

import (
	"bytes"
)

type ChannelMgr struct {
	channels map[string]*Channel
}

func NewChannelMgr() *ChannelMgr {
	return &ChannelMgr{
		channels:make(map[string]*Channel),
	}
}

func (this *ChannelMgr) Info() string {
	buff := new(bytes.Buffer)
	for item,_ := range this.channels {
		buff.WriteString(item)
		buff.WriteString("|")
	}

	return buff.String()
}

func (this *ChannelMgr) OnConnect(channel *Channel) {
	this.channels[channel.socket.RemoteAddr().String()] = channel
}
func (this *ChannelMgr) OnClose(channel *Channel) {
	delete(this.channels, channel.socket.RemoteAddr().String())
}
