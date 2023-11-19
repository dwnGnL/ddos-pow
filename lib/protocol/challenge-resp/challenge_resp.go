package challenge_resp

import (
	"fmt"
	"strings"
)

const (
	REQUEST_CHALLENGE  = "request-challenge"  // request challenge from server
	RESPONSE_CHALLENGE = "response-challenge" // send challenge to client
	REQUEST_RESOURCE   = "request-resource"   // request resource from server
	RESPONSE_RESOURCE  = "response_resource"  // send resource to client
	FAIL               = "fail"               // in case solution from client is incorrect, send appropriate response
)

// Message for tcp exchange data struct
type Message struct {
	Header  string
	Payload string
}

// Stringify for tcp exchange data format
func (m *Message) Stringify() string {
	return fmt.Sprintf("%s|%s", m.Header, m.Payload)
}

// ParseMessage to parse tcp message into Message
func ParseMessage(str string) (*Message, error) {
	str = strings.TrimSpace(str)

	parts := strings.Split(str, "|")
	if len(parts) < 1 || len(parts) > 2 {
		return nil, fmt.Errorf("message doesn't match protocol")
	}

	msg := Message{
		Header: parts[0],
	}
	if len(parts) == 2 {
		msg.Payload = parts[1]
	}
	return &msg, nil
}
