package conn

import (
	zmq "github.com/kapitan-k/gozmq4"
	"github.com/kapitan-k/messaging"
)

type Conn interface {
	messaging.Conn
	Soc() (soc *zmq.Socket)
}

type ConnZmqMsgSender interface {
	Conn
}

type ConnZmqConnect interface {
	Conn
	ConnectTo(path string) (err error)
}

type ConnZmqBind interface {
	Conn
	Bind(path string) (err error)
}

type BaseConn struct {
	soc *zmq.Socket
}

func BaseConnInit(self *BaseConn, soc *zmq.Socket) {
	self.soc = soc
}

func (self *BaseConn) Soc() (soc *zmq.Socket) {
	return self.soc
}

func (self *BaseConn) UnderlyingConn() interface{} {
	return self.Soc()
}

func (self *BaseConn) Close() (err error) {
	return self.soc.Close()
}

func (self *BaseConn) Send(msg *messaging.Msg) (err error) {
	panic("not implemented")
}

func (self *BaseConn) Recv(msg *messaging.Msg) (err error) {
	panic("not implemented")
}

func (self *BaseConn) RecvFn(msg *messaging.Msg, fn messaging.FnOnMsg) (err error) {
	panic("not implemented")
}

type ConnectedBaseConn struct {
	BaseConn
}

func (self *ConnectedBaseConn) Connect(endpoints ...string) (err error) {
	for _, ep := range endpoints {
		err = self.soc.Connect(ep)
		if err != nil {
			return
		}
	}

	return
}
