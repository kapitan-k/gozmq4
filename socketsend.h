
#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdlib.h>
#include <zmq.h>


int zmq4_send_data(void *soc, const void* data, size_t sz_data);
int zmq4_send_data_multipart(void *soc, size_t num_msgs, const void** datas, size_t* szs);
int zmq4_send_data_with_routing_id(void *soc, const void* data, size_t sz_data, uint32_t routing_id);
int zmq4_send_data_with_group(void *soc, const void* data, size_t sz_data, const void* group, size_t sz_group);


#ifdef __cplusplus
}  /* end extern "C" */
#endif