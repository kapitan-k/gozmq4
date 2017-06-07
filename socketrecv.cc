#include "socketrecv.h"


#include <stdlib.h>
#include <string.h>
#include <vector>
#include "zmq4.h"
#include <zmq.h>



int zmq4_recv_data_part(void* soc, int* pmore, zmq_msg_t* msg, void** msg_data, size_t* msg_len) {
	int more;
	size_t more_size = sizeof(more);
	int rc = zmq_msg_init(msg);
	if (rc != 0) {
		return rc;
	}
	
	rc = zmq_recvmsg(soc, msg, 0);
	if (rc < 0) {
		return rc;
	}
	*msg_len = rc;
	
	rc = zmq_getsockopt(soc, ZMQ_RCVMORE, &more, &more_size);
	if (rc != 0) {
		return rc;
	}

	*pmore = more;
	*msg_data = zmq_msg_data(msg);

	return 0;
}

int zmq4_recv_with_routingid(void* soc, zmq_msg_t* msg, void** msg_data, size_t* msg_len, uint32_t *prouting_id) {
	int more;
	size_t more_size = sizeof(more);
	int rc = zmq_msg_init(msg);
	if (rc != 0) {
		return rc;
	}
	
	rc = zmq_recvmsg(soc, msg, 0);
	if (rc < 0) {
		return rc;
	}
	*msg_len = rc;
	
	*msg_data = zmq_msg_data(msg);

	*prouting_id = zmq_msg_routing_id(msg);

	return 0;
}


int zmq4_recv_data_multipart(void* soc, size_t limit, int *pmore, size_t* num_msgs, zmq_msg_t* msgs, void** datas, size_t *sizes) {
	int more;
	size_t cnt = 0;
	size_t more_size = sizeof(more);
	do {
		size_t sz = 0;
		zmq_msg_t* part = &msgs[cnt];
		int rc = zmq_msg_init(part);
		if (rc != 0) {
			return rc;
		}
		
		rc = zmq_recvmsg(soc, part, 0);
		if (rc < 0) {
			return rc;
		}
		sz = rc;
		
		rc = zmq_getsockopt(soc, ZMQ_RCVMORE, &more, &more_size);
		if (rc != 0) {
			return rc;
		}

		datas[cnt] = zmq_msg_data(part);
		sizes[cnt] = sz;

		++cnt;
	} while (more && cnt < limit);

	*pmore = more;
	*num_msgs = cnt;

	return 0;
}


int zmq4_recv_data_multipart_unknown(void *soc, size_t limit, int *pmore, size_t* num_msgs, zmq_msg_t** msgs, void** datas, size_t *sizes) {
	std::vector<zmq_msg_t> msgs_vec;
	int more;
	size_t more_size = sizeof(more);
	do {
		
		zmq_msg_t part;
		int rc = zmq_msg_init(&part);
		if (rc != 0) {
			return rc;
		}
		
		rc = zmq_recvmsg(soc, &part, 0);
		if (rc < 0) {
			return rc;
		}
		
		rc = zmq_getsockopt(soc, ZMQ_RCVMORE, &more, &more_size);
		if (rc != 0) {
			return -1;
		}
		msgs_vec.push_back(part);
	} while (more && msgs_vec.size() < limit);


	*pmore = more;
	*num_msgs = msgs_vec.size();
	*msgs = msgs_vec.data();

	std::vector<void*> datas_vec(msgs_vec.size());
	std::vector<size_t> sizes_vec(msgs_vec.size());

	for (unsigned i = 0; i < msgs_vec.size(); ++i) {
		zmq_msg_t *msg = &msgs_vec[i];
		datas_vec[i] = zmq_msg_data(msg);
		sizes_vec[i] = zmq_msg_size(msg);
	}

	return 0;
}


void zmq4_msg_close_multi(zmq_msg_t* msgs, size_t num_msgs) {
	for (size_t i = 0; i < num_msgs; ++i) {
		zmq_msg_close(&msgs[i]);
	}
}