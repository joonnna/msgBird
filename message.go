package bird

import (
	"errors"
	"fmt"
	"math/rand"

	messagebird "github.com/messagebird/go-rest-api"
)

var (
	errNoRecipient     = errors.New("Message contained no recipient")
	errOriginator      = errors.New("Message contained no originator")
	errContent         = errors.New("Message contained content")
	errOriginatorSize  = errors.New("Originator value is too large")
	errUnknownEncoding = errors.New("Message contains unknown encoding")

	maxMsgSize      = 160
	splitSize       = 153
	maxSplits       = 9
	headerLength    = 6
	maxAlphanumeric = 11
)

type msg struct {
	Recipient  string `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
}

type processedMsg struct {
	originator string
	recipients []string
	body       string
	params     *messagebird.MessageParams
}

func alphanumeric(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}

	return true
}

// Validates if the message contains invalid parameters.
func (m msg) validate() error {
	if length := len(m.Originator); length <= 0 {
		return errOriginator
	} else if alphanumeric(m.Originator) && length > maxAlphanumeric {
		return errOriginatorSize
	}

	if len(m.Recipient) <= 0 {
		return errNoRecipient
	}

	if len(m.Message) <= 0 {
		return errContent
	}

	// Only support gsm7.
	// Could be of other encodings (8bit, unicode etc).
	for _, r := range m.Message {
		// If rune value is above max 7bit value(127) it's not gsm7.
		if r > 127 {
			return errUnknownEncoding
		}

	}

	return nil
}

// Processes a msg and returns the processed version.
// If splitting is needed, multiple processed msgs are returned.
func (m msg) process() []*processedMsg {
	var ret []*processedMsg
	var splits int
	var body string

	length := len(m.Message)

	if length <= 0 {
		return nil
	}
	// No need for split, sending raw body.
	if length < maxMsgSize {
		params := &messagebird.MessageParams{
			Type:       "sms",
			DataCoding: "auto",
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

	if (length % splitSize) == 0 {
		splits = int(length / splitSize)
	} else {
		splits = int(length/splitSize) + 1
	}

	for i := 0; i < splits; i++ {
		// UDH header specification
		header := make([]byte, headerLength)
		header[0] = byte(5)
		header[1] = byte(0)
		header[2] = byte(3)
		header[3] = byte(refNum)
		header[4] = byte(splits)
		header[5] = byte(i + 1)

		start := i * splitSize
		end := (i + 1) * splitSize

		if i == (splits - 1) {
			body = string(m.Message[start:])
		} else {
			body = string(m.Message[start:end])
		}

		params := &messagebird.MessageParams{
			TypeDetails: make(map[string]interface{}),
			Type:        "binary",
			DataCoding:  "auto",
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
