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
	errNilChan    = errors.New("Received nil channel")
)

type receiver struct {
	l net.Listener

	msgChan chan *msg
}

func newReceiver(msgChan chan *msg) (*receiver, error) {
	if msgChan == nil {
		return nil, errNilChan
	}

	// Change to 0.0.0.0, hostname, or non-loopback IP to expose outside loopback net.
	l, err := net.Listen("tcp", "127.0.0.1:8300")
	if err != nil {
		return nil, err
	}

	return &receiver{
		l:       l,
		msgChan: msgChan,
	}, nil
}

// API endpoint is simply "/"
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

func (r receiver) addr() string {
	return r.l.Addr().String()
}
