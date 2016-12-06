package mut

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"sync/atomic"
	"time"
)

type Conn struct {
	writeChan  chan []byte
	connID     uint64 //set via connection mgr
	closed     bool
	writeBytes uint64
	readBytes  uint64
	socket     *net.TCPConn
	br         *bufio.Reader
	bw         *bufio.Writer
	cfg        *Config
	server     *Server //only need by server connection
	client     *Client //only need by client connection
	userData   interface{}
}

func NewConnection(socket *net.TCPConn, cfg *Config) *Conn {
	c := &Conn{
		cfg:       cfg,
		socket:    socket,
		closed:    false,
		writeChan: make(chan []byte, cfg.PendingWriteNum),

		client: nil,
		server: nil,

		readBytes:  0,
		writeBytes: 0,
	}

	c.resetState()

	return c
}

// ResetState reset the closed state and the reader and the write buffer
// It be fired when Connect or Reconnect
func (c *Conn) resetState() {
	c.closed = false
	c.br = bufio.NewReaderSize(c.socket, c.cfg.ReadBufferSize)
	c.bw = bufio.NewWriterSize(c.socket, c.cfg.WriteBufferSize)
	if c.cfg.handler != nil {
		c.cfg.handler.OnConnect(c)
	}
	c.socket.SetKeepAlive(c.cfg.TcpKeepAlive)
	c.socket.SetLinger(c.cfg.TcpLinger)
	c.socket.SetNoDelay(c.cfg.TcpNoDelay)
}

// ResetCount reset the writeBytes and readBytes
// It should never be called unless you want recount it
// after reconnect .
// And it should be the Client side
func (c *Conn) ResetCount() {
	atomic.StoreUint64(&c.writeBytes, 0)
	atomic.StoreUint64(&c.readBytes, 0)
}

// WriteASync write the Packet async
func (c *Conn) WriteAsync(p Packet) error {
	if c.closed {
		return errors.New("conn has closed")
	}
	buf, err := c.cfg.codec.MessageToBytes(c, p)
	if err != nil {
		return err
	}
	c.writeChan <- buf
	return nil
}

// WriteSync write the Packet no delay
func (c *Conn) WriteSync(p Packet) error {
	if c.closed {
		return errors.New("conn has closed")
	}
	buf, err := c.cfg.codec.MessageToBytes(c, p)
	if err != nil {
		return err
	}
	c.WriteRaw(buf)
	c.Flush()
	return nil
}

// WriteRaw read bytes from the socket
// And it remember the number of the bytes
func (c *Conn) ReadRaw(buf []byte) (int, error) {
	n, err := io.ReadFull(c.br, buf)
	if n > 0 {
		c.incReadBytes(uint64(n))
	}
	return n, err
}

func (c *Conn) ReadBytes(delim byte) (line []byte, err error) {
	return c.br.ReadBytes(delim)
}

func (c *Conn) incReadBytes(delta uint64) uint64 {
	for {
		current := c.readBytes
		next := current + delta
		if atomic.CompareAndSwapUint64(&c.readBytes, current, next) {
			return current
		}
	}
}
func (c *Conn) incWriteBytes(delta uint64) uint64 {
	for {
		current := c.writeBytes
		next := current + delta
		if atomic.CompareAndSwapUint64(&c.writeBytes, current, next) {
			return current
		}
	}
}

// WriteRaw write bytes to the socket
// And it remember the number of the bytes
func (c *Conn) WriteRaw(buf []byte) {
	n := uint64(len(buf))
	if n > 0 {
		c.incWriteBytes(n)
		c.bw.Write(buf)
	}
}

func (c *Conn) Flush() {
	c.bw.Flush()
}

func (c *Conn) WriteAndGet(p Packet) (Packet, error) {
	if c.closed {
		return nil, errors.New("conn has closed")
	}
	if c.cfg.AsyncMode {
		return nil, errors.New("only sync mode can run this method")
	}
	err := c.WriteSync(p)
	if err != nil {
		return nil, err
	}

	packet, err := c.cfg.codec.ReadPacket(c)

	if err != nil {
		logger.Err(err, "mut# connection break")
		c.cfg.handler.OnError(c, err)
		return err, nil
	}
	return packet, nil
}

func (c *Conn) WriteAndGetWithTimeout(p Packet, timeout time.Duration) (Packet, error) {
	f := NewFuture(func() (interface{}, error) {
		return c.WriteAndGet(p)
	})
	v, err := f.GetWithTimeout(timeout)
	return v, err
}

// ReadLoop enter a loop , It decodes bytes use the config's Codec
// If any error occurred , an Error event will be emitted
func (c *Conn) ReadLoop() {
	if !c.cfg.AsyncMode {
		panic("only async mode can start this method")
	}
	for !c.closed {
		packet, err := c.cfg.codec.ReadPacket(c)
		if err != nil {
			logger.Err(err, "mut# connection break")
			c.cfg.handler.OnError(c, err)
		}
		c.cfg.handler.OnMessage(c, packet)
	}
}

// WriteLoop enter a loop , It fetch data from the write channel
// If any error occurred , an Error event will be emitted
func (c *Conn) WriteLoop() {
	if !c.cfg.AsyncMode {
		panic("only async mode can start this method")
	}
	for !c.closed {
		b := <-c.writeChan
		size := len(c.writeChan)

		//write batch
		if size > 0 {
			if size > c.cfg.BatchWriteNum {
				size = c.cfg.BatchWriteNum
			}

			buf := new(bytes.Buffer)
			buf.Write(b)
			for i := 0; i < size; i++ {
				b := <-c.writeChan
				buf.Write(b)
			}
			c.WriteRaw(buf.Bytes())
		} else {
			c.WriteRaw(b)
		}
		c.Flush()
	}
}

// Close close this connect and emmit a OnClose Event
// Remove this connection from SessionManager if it belong to the Server
func (c *Conn) Close() error {
	if c.closed {
		return nil
	}

	logger.Info("mut# close connection: %v=>%v", c.socket.LocalAddr(), c.socket.RemoteAddr())

	c.closed = true

	c.cfg.handler.OnClose(c)
	if c.IsServerConn() {
		c.server.mgr.remove(c)
	}
	close(c.writeChan)
	return c.socket.Close()
}

// GetWriteBytes return the number of the read bytes
// The number keep growing even if the connection has reconnect
func (c *Conn) GetReadBytes() uint64 {
	return atomic.LoadUint64(&c.readBytes)
}

// GetWriteBytes return the number of the written bytes
// The number keep growing even if the connection has reconnect
func (c *Conn) GetWriteBytes() uint64 {
	return atomic.LoadUint64(&c.writeBytes)
}

// Server return the server if have
func (c *Conn) Server() *Server {
	return c.server
}

// Client return the client  if have
func (c *Conn) Client() *Client {
	return c.client
}

// IsServerConn jude if the connection is a server connection
func (c *Conn) IsServerConn() bool {
	return c.server != nil
}

func (c *Conn) IsClosed() bool {
	return c.closed
}

// IsClientConn jude if the connection is a client connection
func (c *Conn) IsClientConn() bool {
	return c.client != nil
}

// Socket represents the raw net.TCPConn
func (c *Conn) Socket() *net.TCPConn {
	return c.socket
}

func (c *Conn) GetUserData() interface{} {
	return c.userData
}

func (c *Conn) SetUserData(userData interface{}) {
	c.userData = userData
}
