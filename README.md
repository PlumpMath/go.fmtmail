Go library for outputting email messages.

Package `fmtmail` builds on top of the standard library's `net/mail`, by
adding a single function:

    func WriteMessage(w io.Writer, *mail.Message) error

...Which outputs the message to `w`.

The basic functionality already works, but there are still some details
to finish up:

* Handle outputting "structured" fields; we can't just split everything
  on character boundaries.
* Go over RFC 5322 and make sure we're hitting all of the edge cases.
  Right now we're probably missing some important stuff.

Released under a simple permissive license, see `COPYING`.
