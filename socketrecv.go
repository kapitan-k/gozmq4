package zmq4

// #include <stdlib.h>
// #include "zmq.h"
// #include "socketrecv.h"
import "C"

import (
	"github.com/kapitan-k/cgohelpers"
	"unsafe"
)

func RecvDataPart(soc *Socket) (hasMore bool, msg C.zmq_msg_t, data []byte, err error) {
	var rc, more C.int
	var msgLen uint64

	var msg_data unsafe.Pointer
	rc, err = C.zmq4_recv_data_part(soc.soc, &more, &msg, &msg_data, (*C.size_t)(&msgLen))
	if rc < 0 {
		err = errget(err)
		return
	}

	err = nil
	l := int(msgLen)
	cgohelpers.SetBytesSliceHeader(&data, uintptr(msg_data), l, l)
	return
}

func RecvWithRoutingID(soc *Socket) (msg C.zmq_msg_t, data []byte, routingID uint32, err error) {
	var rc C.int
	var msgLen uint64

	var msg_data unsafe.Pointer
	rc, err = C.zmq4_recv_with_routingid(soc.soc, &msg, &msg_data, (*C.size_t)(&msgLen), (*C.uint32_t)(unsafe.Pointer(&routingID)))
	if rc < 0 {
		err = errget(err)
		return
	}
	err = nil

	l := int(msgLen)
	cgohelpers.SetBytesSliceHeader(&data, uintptr(msg_data), l, l)
	return
}

func RecvDataPartWithRoutingID(soc *Socket) (hasMore bool, msg C.zmq_msg_t, data []byte, err error) {
	var rc, more C.int
	var msgLen uint64

	var msg_data unsafe.Pointer
	rc, err = C.zmq4_recv_data_part(soc.soc, &more, &msg, &msg_data, (*C.size_t)(&msgLen))
	if rc < 0 {
		err = errget(err)
		return
	}
	err = nil

	l := int(msgLen)
	cgohelpers.SetBytesSliceHeader(&data, uintptr(msg_data), l, l)
	return
}

func RecvMultipart(soc *Socket, limit uint64) (msg FixedMultipartMsg, err error) {
	var rc, more C.int

	msgPtrs := make([]C.zmq_msg_t, limit)
	pdatas := make([]uintptr, limit)
	szs := make([]uint64, limit)

	//ptr := (unsafe.Pointer(&pdatas[0]))
	rc, err = C.zmq4_recv_data_multipart(soc.soc, (C.size_t)(limit), &more,
		(*C.size_t)(&limit),
		(*C.zmq_msg_t)(unsafe.Pointer(&msgPtrs[0])),
		(*unsafe.Pointer)(unsafe.Pointer(&pdatas[0])),
		(*C.size_t)(unsafe.Pointer(&szs[0])),
	)
	if rc < 0 {
		err = errget(err)
		return
	}
	err = nil

	msg.FromZmqMsgs(msgPtrs, pdatas, szs)

	return msg, nil
}
