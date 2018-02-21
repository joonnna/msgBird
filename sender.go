package bird

import (
	"errors"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

var (
	errNilChan = errors.New("Received channel is nil.")
	errNoKey   = errors.New("No access key received.")
)

type sender struct {
	ticker      *time.Ticker
	receiveChan chan *msg
	sendChan    chan *processedMsg
	exitChan    chan bool
	client      *messagebird.Client
}

func newSender(msgCh chan *msg, key string) (*sender, error) {
	if msgCh == nil {
		return nil, errNilChan
	}

	if key == "" {
		return nil, errNoKey
	}

	return &pusher{
		ticker:      time.NewTicker(time.Second * 1),
		receiveChan: msgCh,
		sendChan:    make(chan *processedMsg, cap(msgCh)*2),
		exitChan:    make(chan bool),
		client:      messagebird.NewClient(key),
	}
}

func (s *sender) start() {
	go s.receiveLoop()
	go s.sendLoop()
}

func (s *sender) stop() {
	close(p.exitChan)
}

func (s *sender) receiveLoop() {
	for {
		select {
		case m := <-s.receiveChan:
			if m == nil {
				continue
			}

			processed := m.process()

			for _, p := range processed {
				s.sendChan <- p
			}

		case <-p.exitChan:
			return
		}
	}
}

func (s *sender) sendLoop() {
	for {
		select {
		case <-p.ticker.C:
			m := <-p.sendChan

			s.send(m)
		case <-p.exitChan:
			return
		}
	}
}

func (s *sender) send(m *processedMsg) {
	_, err := s.client.NewMessage(m.originator, m.recipients, m.body, m.params)
	if err != nil {

	}
}
