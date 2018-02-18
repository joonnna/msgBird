package something

import "time"

type pusher struct {
	ticker         *time.Ticker
	msgChan        chan *msg
	processMsgChan chan *msg
	exitChan       chan bool
}

type processedMsg struct {
	originator string
	recipients []string
	body       string
	msgParams  *messagebird.MessageParams
}

func newPusher(msgCh chan *msg) *pusher {
	return &pusher{
		ticker:         time.NewTicker(time.Second * 1),
		msgCHan:        msgCh,
		processMsgChan: make(chan *processsedMsg, cap(msgCh)*2),
		exitChan:       make(chan bool),
	}
}

func (p *pusher) receiveMsgs() {
	for {
		select {
		case m := <-p.msgCHan:
			if m == nil {
				continue
			}

			m.processMsgChan <- processMsg(m)
		case <-p.exitChan:
			return
		}
		time.Sleep(time.Second * 1)
	}
}

func (p *pusher) rateLimiter() {
	for {
		select {
		case <-p.ticker.C:
			m := <-m.processMsgChan
			sendMsg(m)

		case <-p.exitChan:
			return
		}

	}
}

func (p *pusher) stop() {
	close(p.exitChan)
}

func processMsg(m *msg) *processedMsg {
	ret := &processedMsg{
		originator: m.Originator,
	}

	body := processBody(m.Message)

	/*
		params := &messagebird.MessageParams{
			//type:
			datacoding:
			typeDetails:
			mclass: 1,
		}
	*/
}

func processBody(body string) string {

}

func sendMsg(m *msg) {

}
