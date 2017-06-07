package conn

import (
	zmq "github.com/kapitan-k/gozmq4"
)

func SendSinglePartBase(soc *zmq.Socket, data []byte) (err error) {
	_, err = soc.SendBytes(data, 0)
	return
}

func SendSinglePartBaseWithGroup(soc *zmq.Socket, data []byte, group string) (err error) {
	_, err = zmq.SendWithGroup(soc, data, group)
	return
}

func SendWithRoutingIDBase(soc *zmq.Socket, data []byte, routingId uint32) (err error) {
	_, err = zmq.SendWithRoutingID(soc, data, routingId)
	return
}

func SendPubSinglePartBase(soc *zmq.Socket, data []byte, pub []byte) (err error) {
	if len(data) > 0 {
		_, err = soc.SendBytes(pub, zmq.SNDMORE)
		if err != nil {
			return
		}
		_, err = soc.SendBytes(data, 0)
		return
	}

	_, err = soc.SendBytes(data, 0)
	return
}
