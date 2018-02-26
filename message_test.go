package bird

import (
	"testing"
)

func TestValidate(t *testing.T) {
	m1 := &msg{Recipient: "test", Originator: "test2"}

	if err := m1.validate(); err == nil {
		t.Error("Should return error with non-existing message body")
	}

	m2 := &msg{Recipient: "test", Message: "test message"}
	if err := m2.validate(); err == nil {
		t.Error("Should return error with non-existing originator")
	}

	m3 := &msg{Originator: "test", Message: "test message"}
	if err := m3.validate(); err == nil {
		t.Error("Should return error with non-existing recipient")
	}

	m4 := &msg{Recipient: "test", Originator: "tooLongOriginator", Message: "test message"}
	if err := m4.validate(); err == nil {
		t.Error("Should return error with originator of excessive size")
	}

	m5 := &msg{Recipient: "test", Originator: "******************", Message: "test message"}
	if err := m5.validate(); err != nil {
		t.Error("Non alpanumeric originator should be allowed to exceed size limit")
	}

}

func TestProcess(t *testing.T) {
	m1 := &msg{Recipient: "test", Originator: "test2"}

	if msgs := m1.process(); msgs != nil {
		t.Error("Should return nil with non-existing message body")
	}

	m2 := &msg{Recipient: "test", Originator: "test2", Message: "test message"}

	if msgs := m2.process(); msgs == nil {
		t.Error("processing a valid message should return non-nil")
	} else if nummsgs := len(msgs); nummsgs != 1 {
		t.Error("messages below max size should not be split")
	}

	var largeBody string
	for i := 0; i <= maxMsgSize; i++ {
		largeBody += "c"
	}

	if length := len(largeBody); length < maxMsgSize {
		t.Errorf("Failed to generate body above max size, got: %d", length)
	}

	m3 := &msg{Recipient: "test", Originator: "test2", Message: largeBody}

	if msgs := m3.process(); msgs == nil {
		t.Error("processing a valid message should return non-nil")
	} else if numMsgs := len(msgs); numMsgs <= 1 {
		t.Error("messages above max size should be split")
	}

}
