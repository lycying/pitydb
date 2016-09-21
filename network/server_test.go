package network

import (
	"testing"
	"net"
	"time"
	"log"
)

type  i_am_a_simple_handler struct {
}

func (this i_am_a_simple_handler) OnRead(chl *Channel, packet Packet) ExchangePacket {
	msg := "got " + string(packet.([]byte))
	chl.Write(msg)
	return msg
}
func (this i_am_a_simple_handler) OnWrite(chl *Channel, packet Packet) ExchangePacket {
	return []byte("ret:I love you!")
}

func (this i_am_a_simple_handler) OnClose(chl *Channel) {
	log.Println("closing channel:", chl.GetRemote())
}

type myInitializer struct {
}

func (myInitializer *myInitializer) InitChannel(channel *Channel) error {
	pipeline := NewPipeline()
	pipeline.AddFirst("first", i_am_a_simple_handler{})
	channel.SetPipeline(pipeline)
	return nil
}

func TestServer_Accept(t *testing.T) {
	s, _ := NewServer("127.0.0.1:14000", &myInitializer{})

	go s.Bootstrap()

	time.Sleep(10 * time.Millisecond)

	//conn, err := net.Dial("tcp", "localhost:14000")
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:14000")
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		t.Fatal("Failed to connect to test server")
	}
	conn.Write([]byte("Test message\n"))
	reads := make([]byte, 4096)
	num, _ := conn.Read(reads)

	t.Log(string(reads[:num]))
	time.Sleep(10 * time.Millisecond)
	t.Log(s.channelMgr.Info())

	conn.Close()
	time.Sleep(10 * time.Millisecond)
	t.Log(s.channelMgr.Info())
}
