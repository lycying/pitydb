package network

import "net"

type Server struct {
	channelMgr         *ChannelMgr
	localAddr          *net.TCPAddr
	channelInitializer ChannelInitializer
}

func NewServer(localAddr string, channelInitializer ChannelInitializer) (*Server, error) {
	addr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		localAddr:addr,
		channelMgr:NewChannelMgr(),
		channelInitializer:channelInitializer,
	}

	return server, nil
}

func (this *Server) Bootstrap() {
	listener, err := net.ListenTCP("tcp", this.localAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, _ := listener.AcceptTCP()
		chl := NewChannel(conn)
		chl.socket = conn

		this.channelInitializer.(ChannelInitializer).InitChannel(chl)

		this.channelMgr.OnConnect(chl)
		chl.pipeline.FireConnect(chl)

		go chl.ReadLoop()
	}
}
