package linebased

import (
	"github.com/lycying/pitydb/tcp/mut"
)

type Packet struct {
	mut.Packet
	Data string
}

func NewPacket(str string) *Packet {
	return &Packet{
		Data: str,
	}
}

type Codec struct {
}

func NewCodec() *Codec {
	return &Codec{}
}
func (codec *Codec) ReadPacket(c *mut.Conn) (mut.Packet, error) {
	packet := &Packet{}

	buf, err := c.ReadBytes('\n')
	if err != nil {
		return packet, err
	}
	packet.Data = string(buf)

	return packet, nil
}

func (codec *Codec) MessageToBytes(c *mut.Conn, p mut.Packet) ([]byte, error) {
	packet := p.(*Packet)
	return []byte(packet.Data), nil
}
