package sdp

import (
	"strings"
	"time"
)

// Attribute represents key:value pair.
type Attribute struct {
	Key   string
	Value string
}

// Attributes is list of k:v.
type Attributes []Attribute

// Value returns value of first attribute.
func (a Attributes) Value(attribute string) string {
	for _, v := range a {
		if v.Key == attribute {
			return v.Value
		}
	}
	return blank
}

// Values returns list of values associated to attribute.
func (a Attributes) Values(attribute string) []string {
	var values []string
	for _, v := range a {
		if v.Key == attribute {
			values = append(values, v.Value)
		}
	}
	return values
}

// Flag returns true if set.
func (a Attributes) Flag(flag string) bool {
	for _, v := range a {
		if v.Key == flag {
			return true
		}
	}
	return false
}

// Message is top level abstraction.
type Message struct {
	Version       int
	Origin        Origin
	Name          string
	Info          string
	Email         string
	Phone         string
	URI           string
	Connection    ConnectionData
	Attributes    Attributes
	Medias        Medias
	Encryption    Encryption
	Bandwidths    map[BandwidthType]int
	BandwidthType BandwidthType
	Timing        []Timing
	TZAdjustments []TimeZone
	SSRC          string
}

// Timing wraps "repeat times" and "timing" information.
type Timing struct {
	Start   time.Time
	End     time.Time
	Repeat  time.Duration
	Active  time.Duration
	Offsets []time.Duration
}

// Start returns start of session.
func (m *Message) Start() time.Time {
	if len(m.Timing) == 0 {
		return time.Time{}
	}
	return m.Timing[0].Start
}

// End returns end of session.
func (m *Message) End() time.Time {
	if len(m.Timing) == 0 {
		return time.Time{}
	}
	return m.Timing[0].End
}

// Flag returns true if set.
func (m *Message) Flag(flag string) bool {
	if len(m.Attributes) > 0 {
		return m.Attributes.Flag(flag)
	}
	return false
}

// Attribute returns string v.
func (m *Message) Attribute(attribute string) string {
	if len(m.Attributes) > 0 {
		return m.Attributes.Value(attribute)
	}
	return blank
}

// AddAttribute appends new k-v pair to attribute list.
func (m *Message) AddAttribute(k, v string) {
	m.Attributes = addAttribute(m.Attributes, k, v)
}

// AddFlag appends new flag to attribute list.
func (m *Message) AddFlag(f string) {
	m.AddAttribute(f, blank)
}

// Medias is list of Media.
type Medias []Media

// Encryption wraps encryption Key and Method.
type Encryption struct {
	Method string
	Key    string
}

// Blank determines whether Encryption is blank value.
func (e Encryption) Blank() bool {
	return e.Equal(Encryption{})
}

// Equal returns e == b.
func (e Encryption) Equal(b Encryption) bool {
	return e == b
}

// Bandwidths is map of BandwidthsType and int (bytes per second).
type Bandwidths map[BandwidthType]int

// Media is media description and attributes.
type Media struct {
	Title       string
	Description MediaDescription
	Connection  ConnectionData
	Attributes  Attributes
	Encryption  Encryption
	Bandwidths  Bandwidths
}

// PayloadFormat returns payload format from a=rtpmap.
// See RFC 4566 Section 6.
func (m *Media) PayloadFormat(payloadType string) string {
	for _, v := range m.Attributes.Values("rtpmap") {
		if strings.HasPrefix(v, payloadType) {
			return strings.TrimSpace(
				strings.TrimPrefix(v, payloadType),
			)
		}
	}
	return ""
}

// AddAttribute appends new k-v pair to attribute list.
func (m *Media) AddAttribute(k string, values ...string) {
	v := strings.Join(values, " ")
	m.Attributes = addAttribute(m.Attributes, k, v)
}

// AddFlag appends new flag to attribute list.
func (m *Media) AddFlag(f string) {
	m.AddAttribute(f, blank)
}

// Flag returns true if set.
func (m *Media) Flag(f string) bool {
	return m.Attributes.Flag(f)
}

// Attribute returns string v.
func (m *Media) Attribute(k string) string {
	return m.Attributes.Value(k)
}
