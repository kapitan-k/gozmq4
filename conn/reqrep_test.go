package conn

import (
	. "github.com/kapitan-k/goutilities/errors"
	"github.com/kapitan-k/messaging"
	. "github.com/kapitan-k/messaging/protocol/flatbinary"
	"time"
	//"github.com/satori/go.uuid"
	. "github.com/kapitan-k/gotestutil"
	. "github.com/kapitan-k/gozmq4"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"testing"
)

func TestReqRep(t *testing.T) {
	path := "inproc://blablaaaaa" + string(RandomData(256))
	rand.Seed(time.Now().Unix())
	//var dataRcv []byte
	ctx, err := NewContext()
	require.NoError(t, err)
	defer ctx.Term()

	reqID := RandomMaxUint64()
	reqID2 := RandomMaxUint64()

	msg := &messaging.Msg{
		Data: RandomDataWithMin(HeaderByteSz, 4096),
	}

	msg2 := &messaging.Msg{
		Data: RandomDataWithMin(HeaderByteSz, 4096),
	}

	// set the request ids in header
	{
		hed := Header{}
		hed.RequestID = reqID
		HeaderToBuf(&hed, msg.Data)
		hed.RequestID = reqID2
		HeaderToBuf(&hed, msg2.Data)
	}

	req := messaging.DefaultRequestNew(msg)
	req2 := messaging.DefaultRequestNew(msg2)

	fnErrClient := func(cerr error) {
		// only prints because error is received at close of ctx
		log.Println("error on client", cerr)
	}

	fnErrServer := func(cerr error) {
		log.Fatalln("error on server", cerr)
	}

	fnRequestIDExtractor := func(data []byte) uint64 {
		return HeaderByBuf(data).RequestID
	}

	repServer := RepServerNew()
	err = repServer.Init(ctx, SERVER, fnErrServer)
	require.NoError(t, err)
	defer repServer.Close()
	err = repServer.Bind(path)
	require.NoError(t, err)

	reqClient := ReqClientNew()
	err = reqClient.Init(ctx, CLIENT, fnRequestIDExtractor, fnErrClient)
	require.NoError(t, err)
	err = reqClient.Start()
	require.NoError(t, err)
	defer reqClient.Close()
	err = reqClient.Connect(path)
	require.NoError(t, err)

	{
		err = reqClient.Request(reqID, msg, req, time.Second*2)
		require.NoError(t, err)

		err = reqClient.Request(reqID2, msg, req2, time.Second*2)
		require.NoError(t, err)

		// test simple receive

		msgRecv := &messaging.Msg{}
		err = repServer.Recv(msgRecv)
		require.NoError(t, err)
		require.Equal(t, msg.Data, msgRecv.Data)
		err = repServer.Send(msgRecv)
		require.NoError(t, err)

		resMsg := req.AwaitResult()
		require.Equal(t, msg.Data, resMsg.Data)
	}

	rTo := req2.AwaitResult()
	require.Equal(t, ErrTimeout, rTo.Err)

	require.Equal(t, 0, len(reqClient.mrh.Map()))

}
