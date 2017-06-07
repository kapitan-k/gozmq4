package zmq4

// #include <stdlib.h>
// #include "zmq.h"
// #include "socketrecv.h"
import "C"

import (
	"reflect"
	"unsafe"
)

type MsgPart struct {
	msgPtr *C.zmq_msg_t
}

func (self *MsgPart) Free() {
	C.zmq_msg_close(self.msgPtr)
}

type FixedMultipartMsg struct {
	msgPtrs []C.zmq_msg_t
	Datas   [][]byte
}

func (self *FixedMultipartMsg) Free() {
	C.zmq4_msg_close_multi((*C.zmq_msg_t)(unsafe.Pointer(&self.msgPtrs[0])), (C.size_t)(len(self.msgPtrs)))
}

func (self *FixedMultipartMsg) FromZmqMsgs(msgPtrs []C.zmq_msg_t, ptrs []uintptr, szs []uint64) {
	self.msgPtrs = msgPtrs
	l := len(ptrs)
	datas := make([][]byte, l)
	for i, ptr := range ptrs {
		sz := int(szs[i])
		datas[i] = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
			Data: ptr,
			Len:  sz,
			Cap:  sz,
		}))
	}

	self.Datas = datas
}

func (self *FixedMultipartMsg) NewSingleBufferCopyBytes() (data []byte) {
	if self.Datas == nil {
		return nil
	}

	var l int
	for _, dat := range self.Datas {
		l += len(dat)
	}

	data = make([]byte, l)
	l = 0
	for _, dat := range self.Datas {
		copy(data[l:], dat)
		l += len(dat)
	}

	return
}
