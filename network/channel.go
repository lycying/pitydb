package network

import (
	"net"
	"log"
	"bufio"
)

type Channel struct {
	//server     *Server
	pipeline   *Pipeline
	writeBytes int
	readBytes  int
	socket     *net.TCPConn
	buffIn     []byte

	bufr       *bufio.Reader
	bufw       *bufio.Writer
}

func NewChannel(conn *net.TCPConn) *Channel {
	chl := &Channel{
		//server:server,
		socket:conn,
		bufr:bufio.NewReader(conn),
		bufw:bufio.NewWriter(conn),
		writeBytes:0,
		readBytes:0,
		buffIn:make([]byte, 8096),
	}

	return chl
}
//func (this *Channel) Bind(tcpAddr *net.TCPAddr) error {
//	receiveSocket, err := net.ListenTCP("tcp", tcpAddr)
//	if nil != err {
//		return err
//	}
//	this.server.listener = receiveSocket
//	return nil
//}

func (this *Channel) Connect(tcpAddr *net.TCPAddr) error {
	conn, err := net.DialTCP("tcp", this.GetLocal(), tcpAddr)
	if err != nil {
		return err
	}
	this.socket = conn
	this.pipeline.FireConnect(this)
	return nil
}
func (this *Channel) Ready() {}
func (this *Channel) Reconnect() {}
func (this *Channel) ReadLoop() {
	defer this.Close()

	for {
		msgLen, err := this.bufr.Read(this.buffIn)
		if err != nil {
			log.Println("error reading ", err)
			break
		}

		//log.Println("data read:", string(this.buffIn))

		this.pipeline.FireRead(this, this.buffIn[:msgLen])
		this.readBytes += msgLen
	}
}
func (this *Channel) Write(data interface{}) (int, error) {
	b := this.pipeline.FireWrite(this, data)
	this.bufw.Write(b.([]byte))
	this.bufw.Flush()
	return 0, nil
}
func (this *Channel) Flush() error {
	return this.bufw.Flush()
}
func (this *Channel) Close() {
	log.Println("closing")
	this.socket.Close()
	this.pipeline.FireClose(this)
}
func (this *Channel) CloseForce() {
	this.Close()
}
func (this *Channel) Type() {}
func (this *Channel) IsWriteable() bool {
	return this.bufw.Available() > 0
}
func (this *Channel) IsReadable() bool {
	return this.bufr.Buffered() > 0
}

func (this *Channel) GetLocal() *net.TCPAddr {
	return this.socket.LocalAddr().(*net.TCPAddr)
}
func (this *Channel) GetRemote() *net.TCPAddr {
	return this.socket.RemoteAddr().(*net.TCPAddr)
}
func (this *Channel) GetReadBytes() int {
	return this.readBytes
}
func (this *Channel) GetWriteBytes() int {
	return this.writeBytes
}
func (this *Channel) IsOpen() bool {
	//TODO
	return true
}

func (this *Channel) Pipeline() *Pipeline {
	return this.pipeline
}

func (this *Channel) SetPipeline(pipeline *Pipeline) {
	this.pipeline = pipeline
}
