
#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>
#include <zmq.h>


int zmq4_recv_data_part(void* soc, int* pmore, zmq_msg_t* msg, void** msg_data, size_t* msg_len);

int zmq4_recv_with_routingid(void* soc, zmq_msg_t* msg, void** msg_data, size_t* msg_len, uint32_t *prouting_id);

// often the number of message parts which should be received is known
int zmq4_recv_data_multipart(void* soc, size_t limit, int *pmore, size_t* num_msgs, zmq_msg_t* msgs, void** datas, size_t *sizes);

// receives as many messages as given with ZMQ_RCVMORE flag
// stores these messages in msgs
int zmq4_recv_data_multipart_unknown(void* soc, size_t limit, int *pmore, size_t* num_msgs, zmq_msg_t** msgs, void** datas, size_t *sizes);





void zmq4_msg_close_multi(zmq_msg_t* msgs, size_t num_msgs);


#ifdef __cplusplus
}  /* end extern "C" */
#endif