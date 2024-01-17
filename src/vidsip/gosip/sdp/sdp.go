package sdp

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// DecodeError wraps Reason of error and occurrence Place.
type DecodeError struct {
	Reason string
	Place  string
}

func (e DecodeError) Error() string {
	return fmt.Sprintf("DecodeError in %s: %s", e.Place, e.Reason)
}

func newDecodeError(place, reason string) DecodeError {
	return DecodeError{
		Reason: reason,
		Place:  place,
	}
}

const (
	lineDelimiter   = '='
	fieldsDelimiter = ' '
	newLine         = '\n'
)

// Line of SDP session.
//
// Form
//
//	<type>=<value>
//
// Where <type> MUST be exactly one case-significant character and
// <value> is structured text whose format depends on <type>.
type Line struct {
	Type  Type
	Value []byte
}

// Equal returns true if l == b.
func (l Line) Equal(b Line) bool {
	if l.Type != b.Type {
		return false
	}
	return bytes.Equal(l.Value, b.Value)
}

// Decode parses b into l and returns error if any.
//
// Decode does not reuse b, so it is safe to corrupt it.
func (l *Line) Decode(b []byte) error {
	delimiter := bytes.IndexRune(b, lineDelimiter)
	if delimiter == -1 {
		reason := `delimiter "=" not found`
		err := newDecodeError("line", reason)
		return errors.Wrap(err, "failed to decode")
	}
	if len(b) <= (delimiter + 1) {
		reason := fmt.Sprintf(
			"len(b) %d < (%d + 1), no value found after delimiter",
			len(b), delimiter,
		)
		err := newDecodeError("line", reason)
		return errors.Wrap(err, "failed to decode")
	}
	r, _ := utf8.DecodeRune(b[:delimiter])
	l.Type = Type(r)
	l.Value = append(l.Value, b[delimiter+1:]...)
	return nil
}

func (l Line) String() string {
	return fmt.Sprintf("%s: %s", l.Type, string(l.Value))
}

func appendCLRF(b []byte) []byte {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, '\r')
	b = append(b, buf[:n]...)
	n = utf8.EncodeRune(buf, '\n')
	b = append(b, buf[:n]...)
	return b
}

func appendRune(b []byte, r rune) []byte {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, r)
	b = append(b, buf[:n]...)
	return b
}

// AppendTo appends Line encoded value to b.
func (l Line) AppendTo(b []byte) []byte {
	b = l.Type.appendTo(b)
	b = appendRune(b, lineDelimiter)
	return append(b, l.Value...)
}

// Type of SDP Line is exactly one case-significant character.
type Type rune

func (t Type) appendTo(b []byte) []byte {
	return appendRune(b, rune(t))
}

var typeToStr = map[Type]string{
	TypeAttribute:          "attribute",
	TypePhone:              "phone",
	TypeEmail:              "email",
	TypeConnectionData:     "connection data",
	TypeURI:                "uri",
	TypeSessionName:        "session name",
	TypeOrigin:             "origin",
	TypeProtocolVersion:    "version",
	TypeTiming:             "timing",
	TypeBandwidth:          "bandwidth",
	TypeSessionInformation: "session info",
	TypeRepeatTimes:        "repeat times",
	TypeTimeZones:          "time zones",
	TypeEncryptionKey:      "encryption keys",
	TypeMediaDescription:   "media description",
	TypeSSRC:               "ssrc",
}

func (t Type) String() string {
	s, ok := typeToStr[t]
	if ok {
		return s
	}
	// Falling back to raw value.
	return string(rune(t))
}

// Attribute types as described in RFC 4566.
const (
	TypeProtocolVersion    Type = 'v'
	TypeOrigin             Type = 'o'
	TypeSessionName        Type = 's'
	TypeSessionInformation Type = 'i'
	TypeURI                Type = 'u'
	TypeEmail              Type = 'e'
	TypePhone              Type = 'p'
	TypeConnectionData     Type = 'c'
	TypeBandwidth          Type = 'b'
	TypeTiming             Type = 't'
	TypeRepeatTimes        Type = 'r'
	TypeTimeZones          Type = 'z'
	TypeEncryptionKey      Type = 'k'
	TypeAttribute          Type = 'a'
	TypeMediaDescription   Type = 'm'
	TypeSSRC               Type = 'y'
)

// Session is set of Lines.
type Session []Line

func (s Session) reset() Session {
	return s[:0]
}

// AppendTo appends all session lines to b and returns b.
func (s Session) AppendTo(b []byte) []byte {
	for _, l := range s {
		b = l.AppendTo(b)
		b = appendCLRF(b)
	}
	return b
}

// Equal returns true if b == s.
func (s Session) Equal(b Session) bool {
	if len(s) != len(b) {
		return false
	}
	for i := range s {
		if !s[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func (s Session) getLine(t Type) Line {
	line := Line{
		Type: t,
	}
	// trying to reuse some memory
	l := len(s)
	if cap(s) > l+1 {
		line.Value = s[:l+1][l].Value[:0]
	}
	return line
}

func (s Session) append(t Type, v []byte) Session {
	line := s.getLine(t)
	line.Value = append(line.Value, v...)
	return append(s, line)
}

func (s Session) appendString(t Type, v string) Session {
	line := s.getLine(t)
	line.Value = append(line.Value, v...)
	return append(s, line)
}

// sliceScanner is custom in-memory scanner for slice
// that will scan all non-whitespace lines.
type sliceScanner struct {
	pos  int
	end  int
	v    []byte
	line []byte
}

func newScanner(v []byte) sliceScanner {
	return sliceScanner{v: v}
}

func (s sliceScanner) Line() []byte {
	return s.line
}

func (s *sliceScanner) Scan() bool {
	// CPU: suboptimal.
	for {
		s.pos = s.end
		if s.pos >= len(s.v) {
			// EOF
			s.line = s.line[:0]
			s.v = s.v[:0]
			return false
		}
		newLinePos := bytes.IndexRune(s.v[s.pos:], newLine)
		s.end = s.pos + newLinePos + 1
		if newLinePos < 0 {
			// next line symbol not found
			s.end = len(s.v)
		}
		s.line = bytes.TrimSpace(s.v[s.pos:s.end])
		if len(s.line) == 0 {
			continue
		}
		return true
	}
}

// DecodeSession decodes Session from b, returning error if any. Blank
// lines and leading/trialing whitespace are ignored.
//
// If s is passed, it will be reused with its lines.
// It is safe to mutate b.
func DecodeSession(b []byte, s Session) (Session, error) {
	var (
		line Line
		err  error
	)
	scanner := newScanner(b)
	for scanner.Scan() {
		// trying to reuse some memory
		l := len(s)
		if cap(s) > l+1 {
			// picking element from s that is not in
			// slice bounds, but in underlying array
			// and reusing it byte slice
			line.Value = s[:l+1][l].Value[:0]
		}
		if err = line.Decode(scanner.Line()); err != nil {
			break
		}
		s = append(s, line)
		line.Value = nil // not corrupting.
	}
	return s, err
}

// Decode decodes b as SDP message, returning error if any.
func Decode(b []byte) (*Message, error) {
	s, err := DecodeSession(b, nil)
	if err != nil {
		return nil, err
	}
	m := new(Message)
	d := NewDecoder(s)
	if err := d.Decode(m); err != nil {
		return nil, err
	}
	return m, nil
}
