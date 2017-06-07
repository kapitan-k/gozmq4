package t

import (
	. "github.com/kapitan-k/gozmq4"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"testing"
)

func TestZmq(t *testing.T) {
	var socSnd, socRcv *Socket
	var dataRcv []byte
	ctx, err := NewContext()
	require.NoError(t, err)
	defer ctx.Term()

	socSnd, err = ctx.NewSocket(PUB)
	err = socSnd.Bind("inproc://fllflf")
	require.NoError(t, err)
	defer socSnd.Close()

	socRcv, err = ctx.NewSocket(SUB)
	err = socRcv.Connect("inproc://fllflf")
	require.NoError(t, err)
	defer socRcv.Close()
	err = socRcv.SetSubscribe("")
	require.NoError(t, err)

	data := []byte(randString())
	datas := [][]byte{data[:3], data[3:6], data[6:]}

	err = socSnd.SendMultipart(datas)
	require.NoError(t, err)

	dataRcv, err = socRcv.RecvMultipartBytes(4)
	require.NoError(t, err)
	log.Println("dataRcv", string(dataRcv))

	//require.Equal(t, data[:3], dataRcv[:3])
	require.Equal(t, data, dataRcv)
}

/*
func TestZmq(t *testing.T) {
	var socSnd, socRcv *Socket
	var dataRcv []byte
	ctx, err := NewContext()
	require.NoError(t, err)
	defer ctx.Term()

	socSnd, err = ctx.NewSocket(PUB)
	err = socSnd.Bind("inproc://fllflf")
	require.NoError(t, err)
	defer socSnd.Close()

	socRcv, err = ctx.NewSocket(SUB)
	err = socRcv.Connect("inproc://fllflf")
	require.NoError(t, err)
	defer socRcv.Close()
	err = socRcv.SetSubscribe("")
	require.NoError(t, err)

	var data []byte
	data = []byte(randString())

	//datas := [][]byte{data[:3], data[3:6], data[6:]}
	go func(soc *Socket) {
		for {
			data = []byte(randString())
			datas := [][]byte{data[:3], data[3:6], data[6:]}
			err = soc.SendMultipart(datas)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)

		}
	}(socSnd)

	for {
		dataRcv, err = socRcv.RecvMultipartBytes(3)
		require.NoError(t, err)
	}

	//log.Println("dataRcv", string(dataRcv))

	//require.Equal(t, data[:3], dataRcv[:3])
	require.Equal(t, data, dataRcv)
}
*/
func randString() string {
	n := rand.Int63n(256)
	var str string
	if n == 0 {
		n = 1
	}
	for i := int64(0); i < n; i++ {
		str = str + uuid.NewV4().String()
	}

	return str
}
