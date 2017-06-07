package zmq4

import (
	. "github.com/kapitan-k/goutilities/errors"
	. "github.com/kapitan-k/goutilities/junk"
	messaging "github.com/kapitan-k/messaging"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

// partly overlapping with czmq
// taken from old project

func ZmqProt(path string) (prot string) {
	return path[0:6]
}

const (
	zmq_tcp_prot    = "tcp://"
	zmq_epgm_prot   = "epgm:/"
	zmq_inproc_prot = "inproc"
)

func ZmqCheckPath(path string) {
	prot := ZmqProt(path)
	if prot != zmq_tcp_prot && prot != zmq_epgm_prot && prot != zmq_inproc_prot {
		panic("sorry, no other protocols implemented")
	}
}

// if path starts with tcp://
// if path in format tcp://127.0.0.1:port tries to bind to this port
// if path in format tcp://127.0.0.1:* tries to bind to arbitrary port starting at port
func ZmqBindPortV4(soc *Socket, path string, port uint16, idService string) (rpath string, rport uint16, err error) {
	prot := ZmqProt(path)
	if prot == zmq_tcp_prot && path[len(path)-1] == '*' {
		path_orig := path[:len(path)-2]
		path, port, err = messaging.BindArbitraryPortV4(soc, path_orig, port)
		if err != nil {
			return
		}
	} else {
		if path[len(path)-1] == ':' {
			path = path[:len(path)-1] + ":" + strconv.Itoa(int(port))
		} else {
			if port != 0 {
				path = path + ":" + strconv.Itoa(int(port))
			}
		}

	}

	if idService == "" {
		soc.SetIdentity(path)
	} else {
		soc.SetIdentity(idService)
	}

	return path, port, soc.Bind(path)

}

// use golang duration
func ZmqEmptySockets(timeout int64, socs ...*Socket) error {
	if timeout < 500 {
		panic("set timeout sufficiently")
	}
	poller := Poller{}
	for _, soc := range socs {
		poller.Add(soc, POLLIN)
	}
	tim := TimeNowUTCMs()
	timEnd := tim + timeout
	var polled []Polled
	var err error
	for {
		polled, err = poller.Poll(0)
		if err == nil && len(polled) > 0 {
			for _, polli := range polled {
				soc := polli.Socket
				soc.Recv(0)
				emptySoc(soc)
			}
		}
		if err != nil || polled == nil || len(polled) == 0 {
			return nil
		}
		tim = TimeNowUTCMs()
		if tim > timEnd {
			break
		}
	}
	if polled != nil {
		return ErrFail
	}
	return nil
}

func emptySoc(soc *Socket) {
	for {
		if b, err := soc.GetRcvmore(); !b || err != nil {
			return
		} else {
			soc.RecvBytes(0)
		}
	}

}

func RandomInproc(serviceID string) string {
	return ToValidInproc("inproc://" + serviceID + strings.ToLower(uuid.NewV4().String()))
}

func ToValidInproc(addr string) string {
	return strings.Replace(addr, "-", "", -1)
}
