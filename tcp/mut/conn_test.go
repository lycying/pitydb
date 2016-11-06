package mut_test

import (
	"bytes"
	"github.com/lycying/pitydb/log"
	"github.com/lycying/pitydb/tcp/mut"
	"testing"
)

var logger, _ = log.New(log.DEBUG, "")

func TestConnection_ReadAll(t *testing.T) {
	var _innerPacket string
	cfg := mut.DefaultConfig()

	func() {
		client := mut.NewClient("amazon.com:80", cfg)
		err := client.DialSync()
		if err != nil {
			logger.Err(err, "error dial")
			return
		}
		buf := new(bytes.Buffer)
		buf.WriteString("GET / HTTP/1.0\r\n\r\n\r\n")
		client.Conn().WriteRaw(buf.Bytes())
		client.Conn().Flush()
		_innerPacket = string(func() []byte {
			b, _ := client.ReadAll()
			return b
		}())
	}()

	logger.Info("%v", _innerPacket)
}
