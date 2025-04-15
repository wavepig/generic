package mq

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMQ(t *testing.T) {
	indexMQ := NewInMemoryMQ[string]()

	for i := 0; i < 10; i++ {
		indexMQ.Publish(Message[string]{Content: fmt.Sprintf("%d", i)})
	}

	for i := 0; i < 10; i++ {
		indexMQ.Consume(func(msg Message[string]) {
			indent, err := json.MarshalIndent(msg, "", "  ")
			if err != nil {
				return
			}
			t.Log(string(indent))
		})
	}
}

func TestMQ2(t *testing.T) {
	indexMQ := make(chan Message[string], 10)

	for i := 0; i < 10; i++ {
		indexMQ <- Message[string]{Content: fmt.Sprintf("%d", i)}
	}

	for i := 0; i < 10; i++ {
		msg := <-indexMQ
		indent, err := json.MarshalIndent(msg, "", "  ")
		if err != nil {
			return
		}
		t.Log(string(indent))
	}
}
