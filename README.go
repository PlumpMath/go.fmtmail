// Package `fmtmail` extends the standard library's `net/mail` to add support
// for outputting mail, not just parsing it.
//
// The basic functionality already works, but there are still some details
// to finish up:
//
// * Handle outputting "structured" fields; we can't just split everything
//   on character boundaries.
// * Go over RFC 5322 and make sure we're hitting all of the edge cases.
//   Right now we're probably missing some important stuff.
// * Investigate what we need to do to accomodate MIME.
//
// Released under a simple permissive license, see `COPYING`.
package fmtmail
