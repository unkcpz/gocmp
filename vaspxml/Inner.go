// modified from "mellium.im/xmlstream"
package vaspxml

import (
	"encoding/xml"
	"io"
)

// Inner returns a new TokenReader that returns nil, io.EOF when it consumes the
// end element matching the most recent start element already consumed.
func Inner(r xml.TokenReader) xml.TokenReader {
	count := 1
	return ReaderFunc(func() (xml.Token, error) {
		if count < 1 {
			return nil, io.EOF
		}

		t, err := r.Token()
		if err != nil {
			return nil, err
		}
		switch t.(type) {
		case xml.StartElement:
			count++
		case xml.EndElement:
			count--
			if count < 1 {
				return nil, io.EOF
			}
		}

		return t, err
	})
}

// ReaderFunc type is an adapter to allow the use of ordinary functions as an
// TokenReader. If f is a function with the appropriate signature,
// ReaderFunc(f) is an TokenReader that calls f.
type ReaderFunc func() (xml.Token, error)

// Token calls f.
func (f ReaderFunc) Token() (xml.Token, error) {
	return f()
}
