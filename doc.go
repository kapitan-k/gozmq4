/*
A Go interface to ZeroMQ (zmq, 0mq) version 4.
Fork from [http://github.com/pebbe/zmq4]
With some additions to send and receive messages with less go calls
Package conn provides specific connections

This includes partial support for ZeroMQ 4.2 DRAFT. The API pertaining
to this support is subject to change.

For ZeroMQ version 3, see: http://github.com/pebbe/zmq3

For ZeroMQ version 2, see: http://github.com/pebbe/zmq2

http://www.zeromq.org/

See also the wiki: https://github.com/pebbe/zmq4/wiki


*/
package zmq4
