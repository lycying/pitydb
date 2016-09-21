package network

type ChannelInitializer interface {
	InitChannel(channel *Channel) error
}