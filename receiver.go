package bird

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

var (
	errMaxRequest = errors.New("Receiving too many requests, throttling")
	errNoMsg      = errors.New("Bodr contains no message")
)

type receiver struct {
	l net.Listener

	msgChan chan *msg
	gsm     *converter
}

func newReceiver(msgChan chan *msg) (*Receiver, error) {
	l, err := net.Listen("tcp", ":")
	if err != nil {
		return nil, err
	}

	return &Receiver{
		l:       l,
		msgChan: msgChan,
		gsm:     newConverter(),
	}, nil
}

func (r *receiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var msgs []*msg

	m := &msg{}

	err := json.NewDecoder(req.Body).Decode(m)
	if err != nil {
		http.Error(w, errNoMsg.Error(), http.StatusBadRequest)
		return
	}

	if err := m.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(m.Message) > maxMsgSize {
		m.gsm, err = r.gsm.convert(m.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	select {
	case r.msgChan <- m:
	default:
		http.Error(w, errMaxRequest.Error(), http.StatusServiceUnavailable)
	}
}

func (r *receiver) start() {
	go http.Serve(r.l, r)
}

func (r *receiver) stop() {
	r.l.Close()
}
