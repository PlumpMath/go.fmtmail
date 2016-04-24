package fmtmail

import (
	"bytes"
	"net/mail"
	"strings"
	"testing"
	"time"
)

func TestSetDate(t *testing.T) {
	now := time.Now()
	msgIn := &Message{
		Header: map[string][]string{
			"To":      []string{"Alice <alice@example.com>"},
			"From":    []string{"Bob <bob@example.net>"},
			"Subject": []string{"SetDate method"},
		},
		Body: strings.NewReader(""),
	}
	msgIn.SetDate(now)
	buf := &bytes.Buffer{}
	_, err := msgIn.WriteTo(buf)
	if err != nil {
		t.Fatalf("TestSetDate: Error writing message: %v\n", err)
	}
	msgOut, err := mail.ReadMessage(buf)
	if err != nil {
		t.Fatalf("TestSetDate: Error reading message: %v\n", err)
	}
	outDate, err := msgOut.Header.Date()
	if err != nil {
		t.Fatalf("TestSetDate: Error getting date: %v\n", err)
	}
	if !outDate.Equal(now) {
		t.Fatalf("TestSetDate: Dates were unequal; %v vs %v", outDate, now)
	}
}
