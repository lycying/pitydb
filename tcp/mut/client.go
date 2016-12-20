package mut

import (
	"errors"
	"io/ioutil"
	"net"
)

type Client struct {
	serverAddr string
	cfg        *Config
	conn       *Conn
	connected  bool
	closed     bool
}

func NewClient(serverAddr string, cfg *Config) *Client {
	return &Client{
		serverAddr: serverAddr,
		cfg:        cfg,
		connected:  false,
		closed:     false,
	}
}
func (client *Client) DialAsync() error {
	retry := newRetry()
	for !client.closed {
		tcpAddress, err := net.ResolveTCPAddr("tcp", client.serverAddr)

		if err != nil {
			logger.Error("mut# ResolveTCPAddr failed: %v", client.serverAddr)
		}
		socket, err := net.DialTCP("tcp", nil, tcpAddress)
		if err != nil {
			logger.Error("mut# DialTCP failed: %v", client.serverAddr)
		}

		if err != nil {
			retry.retryAfter(client)
			continue
		}

		client.connected = true
		client.closed = false
		//reconnect
		if client.conn != nil {
			client.conn.socket = socket
			client.conn.resetState()
		} else {
			conn := NewConnection(socket, client.cfg)
			client.conn = conn
			conn.client = client
		}

		client.Go()
		return nil
	}

	return errors.New("timeout or be closed ")
}
func (client *Client) DialSync() error {
	tcpAddress, err := net.ResolveTCPAddr("tcp", client.serverAddr)

	if err != nil {
		logger.Err(err, "mut# ResolveTCPAddr failed: ")
		return err
	}

	socket, err := net.DialTCP("tcp", nil, tcpAddress)

	if err != nil {
		logger.Err(err, "mut# DialTCP failed:")
		return err
	}

	client.closed = false
	client.connected = true
	//reconnect
	if client.conn != nil {
		client.conn.socket = socket
		client.conn.resetState()
	} else {
		conn := NewConnection(socket, client.cfg)
		client.conn = conn
		conn.client = client
	}
	return nil
}
func (client *Client) ReadAll() ([]byte, error) {
	result, err := ioutil.ReadAll(client.conn.br)
	if err != nil {
		logger.Err(err, "%v", err.Error())
		return nil, nil
	}
	return result, nil
}

func (client *Client) Go() {
	go client.conn.ReadLoop()
	go client.conn.WriteLoop()
}

func (client *Client) Close() {
	if client.closed {
		return
	}
	client.closed = true
	client.connected = false
	if nil != client.conn {
		client.conn.Close()
	}
}

func (client *Client) ReConnect() {
	if nil != client.conn {
		client.conn.Close()
	}
	client.closed = false
	client.connected = false
	client.DialAsync()
}
func (client *Client) Conn() *Conn {
	return client.conn
}
func (client *Client) IsConnected() bool {
	return client.connected
}
