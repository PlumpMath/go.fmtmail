package fmtmail

import (
	"bytes"
	"io"
	"net/mail"
	"sort"
	"strings"
	"testing"
)

func checkReadBack(Header mail.Header, body string, t *testing.T) {
	buf := &bytes.Buffer{}
	err := WriteMessage(buf, &mail.Message{
		Header: Header,
		Body:   strings.NewReader(body),
	})
	if err != nil {
		t.Fatalf("Error from WriteMessage in checkReadBack: %q", err)
	}
	text := buf.String()
	t.Logf("Written message text was: %q", text)
	msg2, err := mail.ReadMessage(buf)
	if err != nil {
		t.Logf("Error from ReadMessage in checkReadBack: %q", err)
		t.FailNow()
	}

	checkSameHeaders(Header, msg2.Header, t)
	checkSameBody(strings.NewReader(body), msg2.Body, t)
}

func checkSameBody(body1, body2 io.Reader, t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	io.Copy(buf1, body1)
	io.Copy(buf2, body2)

	text1 := buf1.String()
	text2 := buf2.String()

	if text1 != text2 {
		t.Fatalf("Message bodies differ: %q vs %q", text1, text2)
	}
}

func checkSameHeaders(hdr1, hdr2 mail.Header, t *testing.T) {
	for k, vs := range hdr1 {
		if len(hdr2[k]) != len(vs) {
			t.Fatalf("Message header %q differs: %v vs %v", k, vs, hdr2[k])
		}
		sort.Strings(hdr2[k])
		sort.Strings(vs)
		for i := range vs {
			if hdr2[k][i] != vs[i] {
				t.Fatalf("Message header %q differs: %v vs %v", k, vs, hdr2[k])
			}
		}
	}
}

var testMessages = []struct {
	Header mail.Header
	Body   string
}{
	{
		Header: mail.Header{
			"To":      []string{"Alice <alice@example.com>"},
			"Subject": []string{"Hi"},
		},
		Body: "Hey there!",
	},
	{
		Header: mail.Header{
			"To": []string{
				"Bob <bob@example.net>",
				"Alice <alice@example.com>",
			},
			"From":    []string{"Mallory <evil@example.org>"},
			"Subject": []string{"MWHAAHA!"},
		},
		Body: "I will destroy you!",
	},
}

func TestReadBack(t *testing.T) {
	for i, v := range testMessages {
		t.Logf("TestReadback: Testing message %d...", i)
		checkReadBack(v.Header, v.Body, t)
	}
}