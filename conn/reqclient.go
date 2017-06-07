package conn

import (
	. "github.com/kapitan-k/goutilities/data"
	. "github.com/kapitan-k/goutilities/errors"
	. "github.com/kapitan-k/gozmq4"
	"github.com/kapitan-k/messaging"
	"log"
	"sync"
	"time"
)

// is an async request client
// it uses the binaryprotocol of github.com/kapitan-k/messaging/protocol
// with the header saved in msg.Header
// when sending data, msg.Header is ignored and only msg.Data is considered
// msg.Data must be prefixed with the Header
type ReqClient struct {
	ConnectedBaseConn
	mrh                  messaging.MapRequestHolder
	fnRequestIDExtractor messaging.FnRequestIDExtractor
	fnOnError            func(err error)
	lock                 sync.Mutex
}

func ReqClientNew() (self *ReqClient) {
	self = &ReqClient{}
	self.mrh = messaging.MapRequestHolderCreate(nil)
	return
}

// Inits the client
// Afterwards Connect and Start should be called
func (self *ReqClient) Init(ctx *Context, socType SocketType, fnRequestIDExtractor messaging.FnRequestIDExtractor, fnOnError func(err error)) (err error) {
	var soc *Socket
	if socType != CLIENT {
		log.Fatalln("ReqClient Init: socType should be CLIENT")
	}

	soc, err = ctx.NewSocket(socType)
	if err != nil {
		return
	}

	self.soc = soc
	self.fnRequestIDExtractor = fnRequestIDExtractor
	self.fnOnError = fnOnError

	return
}

// Runs a receive goroutine
func (self *ReqClient) Start() (nerr error) {
	go func() {
		soc := self.soc
		for {
			msg := &messaging.Msg{}
			_, cmsg, data, err := RecvDataPart(soc)
			defer CZmqMsgClose(&cmsg)
			if err != nil {
				if self.fnOnError != nil {
					self.fnOnError(err)
					return
				}
			}

			reqID := self.fnRequestIDExtractor(data)

			self.lock.Lock()
			res := self.mrh.RemoveByID(reqID)
			self.lock.Unlock()
			if res != nil {
				msg.Data = CopyBuf(data)
				res.Resolve(msg)
			}
		}
	}()
	return
}

// returns ErrFail to all requests still open
// and closes the underlying socket
func (self *ReqClient) Close() (err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	m := self.mrh.Map()
	for id, r := range m {
		r.ResolveError(ErrFail)
		delete(m, id)
	}

	return self.soc.Close()
}

// in theory, if you dont hold the lock when sending, the loop might already gather the reply
// before the lock is acquired, therefore send is called while holding the lock
func (self *ReqClient) Request(requestID uint64, msg *messaging.Msg, resolver messaging.Resolver, timeout time.Duration) (err error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	err = self.send(0, msg)
	if err != nil {
		return
	}

	self.mrh.Add(requestID, resolver)

	// timeout per request
	time.AfterFunc(timeout, func() {
		self.lock.Lock()
		res := self.mrh.RemoveByID(requestID)
		self.lock.Unlock()
		if res != nil {
			res.ResolveError(ErrTimeout)
		}

	})

	return
}

func (self *ReqClient) send(requestID uint64, msg *messaging.Msg) (err error) {
	return SendWithRoutingIDBase(self.soc, msg.Data, uint32(requestID))
}
