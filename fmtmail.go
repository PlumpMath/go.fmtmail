package fmtmail

import (
	"bytes"
	"fmt"
	"io"
	"net/mail"
	"strings"
)

// Write msg to w. An error will be returned if:
//
// 1. The writer returns an error during the course of writing.
// 2. The message contains characters that cannot be encoded
//    (mainly non-ascii in the headers). Determining what escaping
//    the standards allow and implementing it is still TODO
func WriteMessage(w io.Writer, msg *mail.Message) (err error) {
	for k, v := range msg.Header {
		for i := range v {
			err = writeHeader(w, k, v[i])
			if err != nil {
				return err
			}
		}
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	_, err = io.Copy(w, msg.Body)
	return err
}

// Write a single header key: value pair.
func writeHeader(w io.Writer, k string, v string) (err error) {
	if _, err = w.Write([]byte(k)); err != nil {
		return err
	}
	if _, err = w.Write([]byte{':'}); err != nil {
		return err
	}
	return writeHeaderValue(len(k)+1, w, v)
}

// Write the value part of a header. `cols` is the starting offset from the
// beginning of the line, e.g. if the header is "To", then cols will be 3
// (two characters plus the colon).
func writeHeaderValue(cols int, w io.Writer, v string) (err error) {
	r := strings.NewReader(v)
	buf := &bytes.Buffer{}
	ch, _, err := r.ReadRune()
	for err == nil {
		if ch == ':' || ch > '~' || (ch != '\t' && ch < ' ') {
			return fmt.Errorf(
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
		return err
	}
	buf.Write([]byte("\r\n"))
	buf.String()
	_, err = buf.WriteTo(w)
	return err
}
