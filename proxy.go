package bird

type Proxy struct {
	r *receiver
	s *sender
}

func NewProxy(key string) (*Proxy, error) {
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
		r: rec,
		s: send,
	}, nil
}

func (p *Proxy) Start() {
	r.start()
	s.start()
}

func (p *Proxy) Stop() {
	r.stop()
	s.stop()
}
