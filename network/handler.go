package network

type ConnectEvent interface {
	OnConnect(channel *Channel)
}
type ReadEvent interface {
	OnRead(channel *Channel, packet Packet) ExchangePacket
}
type WriteEvent interface {
	OnWrite(channel *Channel, packet Packet) ExchangePacket
}
type CloseEvent interface {
	OnClose(channel *Channel)
}
type ReadWriteEvent interface {
	ReadEvent
	WriteEvent
}
type ConnectCloseEvent interface {
	ConnectEvent
	CloseEvent
}
type Handler interface {
	ConnectEvent
	ReadEvent
	WriteEvent
	CloseEvent
}