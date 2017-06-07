package zmq4

/*

#include <zmq.h>
#include <stdlib.h>
#include <string.h>
#include "socketsend.h"

*/
import "C"

import (
	"unsafe"
)

func SendWithRoutingID(soc *Socket, data []byte, routingId uint32) (int, error) {
	return SendUnsafeDataWithRoutingID(soc, unsafe.Pointer(&data[0]), len(data), routingId)
}

func SendWithGroup(soc *Socket, data []byte, group string) (int, error) {
	return SendUnsafeDataWithGroup(soc, unsafe.Pointer(&data[0]), len(data), group)
}

func SendUnsafe(soc *Socket, ptr unsafe.Pointer, sz int, flags Flag) (int, error) {
	size, err := C.zmq_send(soc.Soc(), ptr, C.size_t(sz), C.int(flags))
	if size < 0 {
		return int(size), errget(err)
	}
	return int(size), nil
}

func SendUnsafeData(soc *Socket, ptr unsafe.Pointer, sz int) (int, error) {
	size, err := C.zmq4_send_data(soc.Soc(), ptr, C.size_t(sz))
	if size < 0 {
		return int(size), errget(err)
	}
	return int(size), nil
}

func SendUnsafeDataWithRoutingID(soc *Socket, ptr unsafe.Pointer, sz int, routingId uint32) (int, error) {
	size, err := C.zmq4_send_data_with_routing_id(soc.soc, ptr, C.size_t(sz), C.uint32_t(routingId))
	if size < 0 {
		return int(size), errget(err)
	}
	return int(size), nil
}

func SendUnsafeDataWithGroup(soc *Socket, ptr unsafe.Pointer, sz int, group string) (int, error) {
	bgroup := []byte(group)
	size, err := C.zmq4_send_data_with_group(soc.Soc(), ptr, C.size_t(sz), unsafe.Pointer(&bgroup[0]), C.size_t(len(bgroup)))
	if size < 0 {
		return int(size), errget(err)
	}
	return int(size), nil
}

func SendMultipart(soc *Socket, datas [][]byte) error {
	l := len(datas)

	szs := make([]C.size_t, l)
	ptrs := make([]uintptr, l)
	for i, data := range datas {
		szs[i] = C.size_t(len(data))
		ptrs[i] = uintptr(unsafe.Pointer(&data[0]))
	}

	size, err := C.zmq4_send_data_multipart(soc.soc, (C.size_t)(l), (*unsafe.Pointer)(unsafe.Pointer(&ptrs[0])), (*C.size_t)(unsafe.Pointer(&szs[0])))
	if size < 0 {
		return errget(err)
	}
	return nil
}
