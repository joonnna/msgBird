package something

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
)

var (
	logLocation     = "var/log/something"
	errMaxRequest   = "Receiving too many requests, throttling"
	errNoMsg        = "Body contains no message"
	maxMsgSize      = 160
	maxSplitMsgSize = 153
	maxSplits       = 9
	headerLength    = 6
)

type Receiver struct {
	l net.Listener

	log *log.Logger

	msgChan chan *msg
	gsm     *converter
}

func NewReceiver() (*Receiver, error) {
	logFile, err := os.Create(logLocation)
	if err != nil {
		return nil, err
	}

	localLog := log.New(logFile, "", log.Lshortfile)

	l, err := net.Listen("tcp", ":")
	if err != nil {
		localLog.Println(err)
		return nil, err
	}

	return &Receiver{
		l:       l,
		log:     localLog,
		msgChan: make(chan *msg, 1000),
		gsm:     newConverter(),
	}
}

func (r *Receiver) handler(w ResponseWriter, r *Request) {
	defer r.Body.Close()

	var msgs []*msgs

	m := &msg{}

	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		r.log.Println(err)
		http.Error(w, errNoMsg.Error(), http.StatusBadRequest)
		return
	}

	if err := m.validate(); err != nil {
		r.log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g, err := r.gsm.convert(m.Message)
	if err != nil {
		r.log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(g) > maxMsgSize {
		msgs = splitMsg(m, g)
	} else {
		msgs = append(msgs, m)
	}

	for _, m := range msgs {
		select {
		case r.msgChan <- m:
		default:
			http.Error(w, errMaxRequest.Error(), http.StatusServiceUnavailable)
		}
	}
}

func (r *Receiver) Start() {
	go http.ListenAndServe(l, r.handler)
}

func (r *Receiver) Stop() {
	r.l.Close()
}

func splitMessage(m *msg, gsm []byte) []*msg {
	var splits int

	// 1 byte, max 255
	refNum := rand.Int() % 256

	len := len(m.Message)

	if len%maxSplitMsgSize == 0 {
		splits = int(len / maxSplitMsgSize)
	} else {
		splits = int(len/maxSplitMsgSize) + 1
	}

	ret := make([]*msg, splits)

	for i := 0; i < splits; i++ {
		header := make([]byte, headerLength)
		header[0] = byte(5)
		header[1] = byte(0)
		header[2] = byte(3)
		header[3] = byte(refNum)
		header[4] = byte(splits)
		header[5] = byte(i + 1)

		start := i * maxSplitMsgSize
		end := (i + 1) * maxSplitMsgSize

		newMsg := &msg{
			Recipient:  m.Recipient,
			Originator: m.Originator,
			Message:    m.Message[stard:end],
			isSplit:    true,
			header:     fmt.Sprintf("%x", header),
		}

		ret[i] = newMsg
	}

	return ret
}
