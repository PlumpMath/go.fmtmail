package fmtmail

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"time"
)

const timeFormat = time.RFC822

// An email message
type Message mail.Message

// Set the message's Date header to the indicated time.
func (m *Message) SetDate(t time.Time) {
	m.Header["Date"] = append(m.Header["Date"], t.Format(timeFormat))
}

// Write m to w. Can return an error because of:
//
// 1. An error reported by w
// 2. An error encoding the message, e.g. unexpected non-ascii characters.
func (m *Message) WriteTo(w io.Writer) (n int64, err error) {
	var nThis64 int64
	var nThis int
	for k, v := range m.Header {
		for i := range v {
			nThis64, err = writeHeader(w, k, v[i])
			n += nThis64
			if err != nil {
				return
			}
		}
	}
	nThis, err = w.Write([]byte("\r\n"))
	n += int64(nThis)
	if err != nil {
		return
	}
	nThis64, err = io.Copy(w, m.Body)
	n += nThis64
	return
}

// Write a single header key: value pair.
func writeHeader(w io.Writer, k string, v string) (n int64, err error) {
	nThis, err := w.Write([]byte(k))
	n += int64(nThis)
	if err != nil {
		return
	}
	nThis, err = w.Write([]byte{':'})
	n += int64(nThis)
	if err != nil {
		return
	}
	nThis64, err := writeHeaderValue(len(k)+1, w, v)
	n += nThis64
	return
}

// Write the value part of a header. `cols` is the starting offset from the
// beginning of the line, e.g. if the header is "To", then cols will be 3
// (two characters plus the colon).
func writeHeaderValue(cols int, w io.Writer, v string) (n int64, err error) {
	r := strings.NewReader(v)
	buf := &bytes.Buffer{}
	ch, _, err := r.ReadRune()
	for err == nil {
		if ch == ':' || ch > '~' || (ch != '\t' && ch < ' ') {
			return n, fmt.Errorf(
				"Illegal character in header value: %q",
				ch)
		}
		if cols >= 78 {
			buf.WriteString("\r\n ")
			cols = 1
		}
		buf.Write([]byte{byte(ch)})
		cols += 1
		ch, _, err = r.ReadRune()
	}
	if err != io.EOF {
		return
	}
	buf.Write([]byte("\r\n"))
	buf.String()
	return buf.WriteTo(w)
}
