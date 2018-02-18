package something

var (
	errNoRecipient = errrors.New("Message contained no recipient")
	errOriginator  = errrors.New("Message contained no originator")
	errContent     = errrors.New("Message contained content")
)

type msg struct {
	Recipient  int    `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
	isSplit    bool
	header     string
}

func (m msg) validate() error {
	if m.Recipient == "" {
		return errNoRecipient
	}

	if m.Originator == "" {
		return errOriginator
	}

	if len(m.Messge) <= 0 {
		return errContent
	}

	return nil
}

func checkGsm(body string) bool {
}
