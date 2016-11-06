package typelen

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/lycying/pitydb/tcp/mut"
)

const MaxPacketBytes = 32 * 1024

type Packet struct {
	mut.Packet

	Len  uint32
	Type uint32
	Data []byte
}

func NewPacket() *Packet {
	return &Packet{}
}

type Codec struct {
}

func NewCodec() *Codec {
	return &Codec{}
}
func (codec *Codec) ReadPacket(c *mut.Conn) (mut.Packet, error) {
	packet := &Packet{}

	//read len
	bLen := make([]byte, 4)
	_, err := c.ReadRaw(bLen)
	if err != nil {
		return packet, err
	}
	packet.Len = binary.BigEndian.Uint32(bLen)

	if int(packet.Len) > MaxPacketBytes {
		return packet, errors.New(fmt.Sprintf("bytes too long len=%d", packet.Len))
	}

	//read type
	bType := make([]byte, 4)
	_, err0 := c.ReadRaw(bType)
	if err0 != nil {
		return packet, err0
	}
	packet.Type = binary.BigEndian.Uint32(bType)

	//read data
	bData := make([]byte, int(packet.Len))
	_, err1 := c.ReadRaw(bData)
	if err1 != nil {
		return packet, err1
	}
	packet.Data = bData

	return packet, nil
}

func (codec *Codec) MessageToBytes(c *mut.Conn, p mut.Packet) ([]byte, error) {
	packet := p.(*Packet)
	buf := new(bytes.Buffer)

	bLen := make([]byte, 4)
	binary.BigEndian.PutUint32(bLen, packet.Len)

	bType := make([]byte, 4)
	binary.BigEndian.PutUint32(bType, packet.Type)

	buf.Write(bLen)
	buf.Write(bType)
	buf.Write(packet.Data)

	return buf.Bytes(), nil
}
