package bird

import (
	"errors"
	"fmt"
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

	return &sender{
		ticker:      time.NewTicker(time.Second * 1),
		receiveChan: msgCh,
		sendChan:    make(chan *processedMsg, cap(msgCh)*2),
		exitChan:    make(chan bool),
		client:      messagebird.New(key),
	}, nil
}

func (s *sender) start() {
	go s.receiveLoop()
	go s.sendLoop()
}

func (s *sender) stop() {
	close(s.exitChan)
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

		case <-s.exitChan:
			return
		}
	}
}

func (s *sender) sendLoop() {
	for {
		select {
		case <-s.ticker.C:
			m := <-s.sendChan

			s.send(m)
		case <-s.exitChan:
			return
		}
	}
}

func (s *sender) send(m *processedMsg) {
	_, err := s.client.NewMessage(m.originator, m.recipients, m.body, m.params)
	if err != nil {
		fmt.Println(err)
	}
}
