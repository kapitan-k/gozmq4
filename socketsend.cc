#include "socketsend.h"


#include <stdlib.h>
#include <string.h>
#include <zmq.h>


int zmq4_send_data(void *soc, const void* data, size_t sz_data) {
	zmq_msg_t msg0;
    zmq_msg_init_size(&msg0, sz_data);
    memcpy(zmq_msg_data(&msg0), data, sz_data);
    return zmq_msg_send(&msg0, soc, 0);
}


int zmq4_send_data_multipart(void *soc, size_t num_msgs, const void** datas, size_t* szs) {
	if (num_msgs == 1) {
		return zmq_send(soc, datas[0], szs[0], 0);
	}

	for (size_t i = 0; i < num_msgs -1; ++i) {
		int rc = zmq_send(soc, datas[i], szs[i], ZMQ_SNDMORE);
		if (rc <= 0) {
			return rc;
		}
	}

	 return zmq_send(soc, datas[num_msgs-1], szs[num_msgs-1], 0);
}


int zmq4_send_data_with_routing_id(void *soc, const void* data, size_t sz_data, uint32_t routing_id) {
	zmq_msg_t msg0;
    zmq_msg_init_size(&msg0, sz_data);
    memcpy(zmq_msg_data(&msg0), data, sz_data);
	zmq_msg_set_routing_id(&msg0, routing_id);
	return zmq_msg_send(&msg0, soc, 0);
}


int zmq4_send_data_with_group(void *soc, const void* data, size_t sz_data, const void* group, size_t sz_group) {
	zmq_msg_t msg0;
    zmq_msg_init_size(&msg0, sz_data);
    memcpy(zmq_msg_data(&msg0), data, sz_data);

    char* cgroup = (char*)malloc(sz_group+1);
    memset(cgroup, 0, sz_group+1);
    memcpy(cgroup, group, sz_group);
	int rc = zmq_msg_set_group(&msg0, cgroup);
	free(cgroup);

	if (rc <= 0) {
		return rc;
	}

	return zmq_msg_send(&msg0, soc, 0);
} 
