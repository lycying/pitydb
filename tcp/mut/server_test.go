package mut_test

import (
	"fmt"
	"github.com/lycying/pitydb/tcp/mut"
	"github.com/lycying/pitydb/tcp/mut/codec/typelen"
	"sync"
	"testing"
	"time"
)

func getStrPacket(str string) *typelen.Packet {
	p := &typelen.Packet{}
	p.Type = 1
	p.Data = []byte(str)
	p.Len = uint32(len(p.Data))
	return p
}

func _onMsg(conn *mut.Conn, p mut.Packet) {
	x := p.(*typelen.Packet)
	str := string(x.Data)
	side := "Client"
	if conn.IsServerConn() {
		side = "Server"
	}
	logger.Debug("%v: hasRead:=%d,hasWrite:=%d type=%d,len=%d %v",
		side,
		conn.GetReadBytes(),
		conn.GetWriteBytes(),
		x.Type, x.Len, str)

	if str == "close" {
		conn.Close()
	}
	if conn.IsServerConn() {
		conn.WriteAsync(getStrPacket("fuck you !"))
		conn.Server().ConnMgr().Iterate(func(id uint64, c *mut.Conn) bool {
			logger.Debug("MGR:[%v],[%v=>%v]", id, c.Socket().LocalAddr(), c.Socket().RemoteAddr())
			return true
		})
	}

}

func TestServer(t *testing.T) {
	cbk := mut.NewHandlerSkeleton()
	cbk.SetOnMessage(_onMsg)
	cfg := mut.DefaultConfig()
	cfg.SetCallback(cbk)
	cfg.SetCodec(typelen.NewCodec())

	//the server
	server := mut.NewServer(":14000", cfg)
	err := server.Servo()
	if err != nil {
		logger.Err(err, "servo error")
		return
	}

	var grp *sync.WaitGroup = new(sync.WaitGroup)
	//the client
	//make ten
	for i := 0; i < 10; i++ {
		grp.Add(1)
		client := mut.NewClient("localhost:14000", cfg)

		go func() {
			err := client.DialAsync()
			if err != nil {
				logger.Error("dail error")
				return
			}

			for {
				if client.IsConnected() {
					toWrite := getStrPacket(fmt.Sprintf("follow me . p.. %d", i*100))

					client.Conn().WriteAsync(toWrite)

					client.Conn().WriteAsync(getStrPacket("close"))
					time.Sleep(time.Second / 2)
					break
				}
			}

			grp.Done()
		}()

	}

	grp.Wait()

}
