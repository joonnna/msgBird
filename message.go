package bird

import (
	"errors"
	"fmt"
	"math/rand"

	messagebird "github.com/messagebird/go-rest-api"
)

var (
	errNoRecipient = errors.New("Message contained no recipient")
	errOriginator  = errors.New("Message contained no originator")
	errContent     = errors.New("Message contained content")

	maxMsgSize      = 160
	maxSplitMsgSize = 153
	maxSplits       = 9
	headerLength    = 6
)

type msg struct {
	Recipient  string `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
	gsm        []byte
}

type processedMsg struct {
	originator string
	recipients []string
	body       string
	params     *messagebird.MessageParams
}

func (m msg) validate() error {
	if m.Recipient == "" {
		return errNoRecipient
	}

	if m.Originator == "" {
		return errOriginator
	}

	if len(m.Message) <= 0 {
		return errContent
	}

	return nil
}

func (m msg) process() []*processedMsg {
	var ret []*processedMsg
	var splits int
	var body string

	if len(m.Message) < maxMsgSize {
		params := &messagebird.MessageParams{
			Type: "sms",
		}

		newMsg := &processedMsg{
			originator: m.Originator,
			recipients: []string{m.Recipient},
			body:       m.Message,
			params:     params,
		}
		return append(ret, newMsg)
	}

	// 1 byte, max 255
	refNum := rand.Int() % 256

	len := len(m.Message)

	if len%maxSplitMsgSize == 0 {
		splits = int(len / maxSplitMsgSize)
	} else {
		splits = int(len/maxSplitMsgSize) + 1
	}

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

		if i == (splits - 1) {
			body = m.Message[start:]
		} else {
			body = m.Message[start:end]
		}

		params := &messagebird.MessageParams{
			TypeDetails: make(map[string]interface{}),
			Type:        "binary",
			DataCoding:  "plain",
		}
		params.TypeDetails["udh"] = fmt.Sprintf("%x", header)

		newMsg := &processedMsg{
			recipients: []string{m.Recipient},
			originator: m.Originator,
			body:       fmt.Sprintf("%x", body),
			params:     params,
		}

		ret = append(ret, newMsg)
	}

	return ret
}
