package bird

type Proxy struct {
	rec  *receiver
	send *sender
}

func NewProxy(key string) (*Proxy, error) {
	// msgCh represents pending messages(not processed yet).
	// Its size determines how many pending messages are queued
	// before dropping incoming messages and returning errors.
	msgCh := make(chan *msg, 1000)

	rec, err := newReceiver(msgCh)
	if err != nil {
		return nil, err
	}

	send, err := newSender(msgCh, key)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		rec:  rec,
		send: send,
	}, nil
}

func (p *Proxy) Start() {
	p.rec.start()
	p.send.start()
}

func (p *Proxy) Stop() {
	p.rec.stop()
	p.send.stop()
}

func (p Proxy) addr() string {
	return p.rec.addr()
}
