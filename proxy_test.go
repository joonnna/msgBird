package bird

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	p, err := NewProxy("testkey")
	if err != nil {
		t.Fatal(err)
	}

	// By replacing the ticker channel, the sender will no longer receive
	// outbound messages in the sendLoop.
	// Now we can receive all outbound messages on the senders send channel.
	p.send.ticker.C = make(chan time.Time)
	p.Start()

	m := &msg{Recipient: "test", Originator: "test2", Message: "This is a test message."}

	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(m)
	if err != nil {
		t.Fatal(err)
	}

	addr := fmt.Sprintf("http://%s/", p.addr())

	resp, err := http.Post(addr, "application/json", &buf)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Received http error status code")
	}

	out := <-p.send.sendChan

	if out.body != m.Message {
		t.Error("Outbound message has different message body")
	}

	if amount := len(out.recipients); amount <= 0 {
		t.Error("Outbound message has no recipients")
	} else if out.recipients[0] != m.Recipient {
		t.Error("Outbound message has different recipient")
	}

	if out.originator != m.Originator {
		t.Error("Outbound message has different originator")
	}

	p.Stop()
}
