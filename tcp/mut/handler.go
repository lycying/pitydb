package mut

type Handler interface {
	// OnConnect : you may want add some UserData here
	OnConnect(c *Conn)
	// OnMessage only available while async mode
	OnMessage(c *Conn, p Packet)

	OnClose(c *Conn)
	OnError(c *Conn, err error)
}

type HandlerSkeleton struct {
	Handler
	_OnConnect func(c *Conn)
	_OnMessage func(c *Conn, p Packet)
	_OnClose   func(c *Conn)
	_OnError   func(c *Conn, err error)
}

func NewHandlerSkeleton() *HandlerSkeleton {
	h := &HandlerSkeleton{}
	h._OnConnect = func(c *Conn) {
		logger.Debug("mut# callback on connect")
	}
	h._OnMessage = func(c *Conn, p Packet) {
		logger.Debug("mut# callback on msg %v", p)
	}
	h._OnClose = func(c *Conn) {
		logger.Debug("mut# callback on close")
	}
	h._OnError = func(c *Conn, err error) {
		logger.Err(err, "mut# callback on error")
		c.Close()
	}
	return h
}
func (h *HandlerSkeleton) OnConnect(c *Conn) {
	h._OnConnect(c)
}
func (h *HandlerSkeleton) OnClose(c *Conn) {
	h._OnClose(c)
}
func (h *HandlerSkeleton) OnError(c *Conn, err error) {
	h._OnError(c, err)
}
func (h *HandlerSkeleton) OnMessage(c *Conn, p Packet) {
	h._OnMessage(c, p)
}

func (h *HandlerSkeleton) SetOnConnect(f func(c *Conn)) {
	h._OnConnect = f
}
func (h *HandlerSkeleton) SetOnMessage(f func(c *Conn, p Packet)) {
	h._OnMessage = f
}
func (h *HandlerSkeleton) SetOnClose(f func(c *Conn)) {
	h._OnClose = f
}
func (h *HandlerSkeleton) SetOnError(f func(c *Conn, err error)) {
	h._OnError = f
}
