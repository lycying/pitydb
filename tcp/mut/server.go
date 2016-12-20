package mut

import (
	"net"
	"sync"
)

type Server struct {
	closed       bool
	address      string
	cfg          *Config
	mgr          *ConnMgr
	ln           *net.TCPListener
	wgConn       *sync.WaitGroup
	localAddress *net.TCPAddr
}

func NewServer(address string, cfg *Config) *Server {
	return &Server{
		address: address,
		cfg:     cfg,
		closed:  false,
		mgr:     newConnectionMgr(),
		wgConn:  new(sync.WaitGroup),
	}
}

func (srv *Server) Servo() error {
	err := srv.cfg.validate()
	if err != nil {
		return err
	}

	//resolve
	tcpAddress, err := net.ResolveTCPAddr("tcp", srv.address)
	if err != nil {
		logger.Err(err, "mut# ResolveTCPAddr failed:")
		return err
	}
	srv.localAddress = tcpAddress

	//listen
	ln, err := net.ListenTCP("tcp", srv.localAddress)
	if err != nil {
		logger.Err(err, "mut# ListenTCP failed:")
		return err
	}
	srv.ln = ln

	logger.Info("mut# listen %+v  [ok]", srv.ln.Addr())
	go srv.listenLoop()
	return nil
}

func (srv *Server) listenLoop() {
	retry := newRetry()
	for !srv.closed {
		socket, err := srv.ln.AcceptTCP()
		//retry if err occur
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				retry.retryAfter(srv)
				continue
			}
		}
		retry.reset()

		srv.wgConn.Add(1)

		logger.Info("mut# connection arrived: %+v => %v", socket.RemoteAddr(), socket.LocalAddr())
		conn := NewConnection(socket, srv.cfg)

		conn.server = srv
		srv.mgr.add(conn)

		//only async mode start this loop
		if srv.cfg.AsyncMode {
			go conn.ReadLoop()
			go conn.WriteLoop()
		}
	}
}

func (srv *Server) Close() {
	srv.closed = true
	srv.wgConn.Wait()
}

func (srv *Server) ConnMgr() *ConnMgr {
	return srv.mgr
}
