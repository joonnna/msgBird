package bird

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

var (
	errMaxRequest = errors.New("Receiving too many requests, throttling")
	errNoMsg      = errors.New("Body contains no message")
)

type receiver struct {
	l net.Listener

	msgChan chan *msg
	gsm     *converter
}

func newReceiver(msgChan chan *msg) (*receiver, error) {
	l, err := net.Listen("tcp", "127.0.0.1:8300")
	if err != nil {
		return nil, err
	}

	return &receiver{
		l:       l,
		msgChan: msgChan,
		gsm:     newConverter(),
	}, nil
}

func (r *receiver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

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
