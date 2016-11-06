package mut

// Packet
type Packet interface {
}

type Codec interface {
	ReadPacket(c *Conn) (Packet, error)
	MessageToBytes(c *Conn, p Packet) ([]byte, error)
}
