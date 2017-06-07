package conn

import (
	. "github.com/kapitan-k/goutilities/data"
	. "github.com/kapitan-k/gozmq4"
	"github.com/kapitan-k/messaging"
	"log"
)

// async server for use with ReqClient
// it uses the binaryprotocol of github.com/kapitan-k/messaging/protocol
// with the header saved in msg.Header
// when sending data, msg.Header is ignored and only msg.Data is considered
// msg.Data must be prefixed with the Header
// this server uses only a single message part
// You can call Recv handle the request and Send
// the VarData returned in the Recv(msg) msg.VarData must not be overwritten
type RepServer struct {
	ConnectedBaseConn
	fnOnError func(err error)
}

func RepServerNew() (self *RepServer) {
	self = &RepServer{}
	return
}

func (self *RepServer) Init(ctx *Context, socType SocketType, fnOnError func(err error)) (err error) {
	var soc *Socket
	if socType != SERVER {
		log.Fatalln("RepServer Init: socType should be CLIENT")
	}

	soc, err = ctx.NewSocket(socType)
	if err != nil {
		return
	}

	self.soc = soc
	self.fnOnError = fnOnError

	return
}

// the VarData returned msg.VarData must not be overwritten
func (self *RepServer) Recv(msg *messaging.Msg) (err error) {
	cmsg, data, routingID, err := RecvWithRoutingID(self.soc)
	defer CZmqMsgClose(&cmsg)
	if err != nil {
		if self.fnOnError != nil {
			self.fnOnError(err)
			return
		}
	}

	msg.Data = CopyBuf(data)
	msg.VarData = routingID

	return
}

func (self *RepServer) RecvFn(msg *messaging.Msg, fn messaging.FnOnMsg) (err error) {
	err = self.Recv(msg)
	fn(msg)
	return
}

func (self *RepServer) Send(msg *messaging.Msg) (err error) {
	return SendWithRoutingIDBase(self.soc, msg.Data, msg.VarData.(uint32))
}

func (self *RepServer) Bind(path string) (err error) {
	return self.soc.Bind(path)
}

func (self *RepServer) Close() (err error) {
	return self.soc.Close()
}
