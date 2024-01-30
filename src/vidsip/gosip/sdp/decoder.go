package sdp

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
	"unsafe"

	"github.com/pkg/errors"
)

const blank = ""

// Decoder decodes session.
type Decoder struct {
	s       Session
	pos     int
	t       Type
	v       []byte
	l       Line
	section section
	sPos    int
	m       Media
}

// NewDecoder returns Decoder for Session.
func NewDecoder(s Session) Decoder {
	return Decoder{
		s: s,
	}
}

func (d *Decoder) newFieldError(msg string) DecodeError {
	return DecodeError{
		Place:  fmt.Sprintf("%s/%s at line %d", d.section, d.t, d.pos),
		Reason: msg,
	}
}

func (d *Decoder) next() bool {
	if d.pos >= len(d.s) {
		return false
	}
	d.l = d.s[d.pos]
	d.v = d.l.Value
	d.t = d.l.Type
	d.pos++
	return true
}

type section int

const (
	sectionSession section = iota
	sectionTime
	sectionMedia
)

func (s section) String() string {
	switch s {
	case sectionSession:
		return "s"
	case sectionTime:
		return "t"
	case sectionMedia:
		return "m"
	default:
		panic("BUG: section overflow")
	}
}

type ordering []Type

var orderingSession = ordering{
	TypeProtocolVersion,
	TypeOrigin,
	TypeSessionName,
	TypeSessionInformation,
	TypeURI,
	TypeEmail,
	TypePhone,
	TypeConnectionData,
	TypeBandwidth,     // 0 or more
	TypeTimeZones,     // *
	TypeEncryptionKey, // ordering after time start
	TypeAttribute,     // 0 or more
}

const orderingAfterTime = 10

var orderingTime = ordering{
	TypeTiming,
	TypeRepeatTimes,
}

var orderingMedia = ordering{
	TypeMediaDescription,
	TypeSessionInformation, // title
	TypeConnectionData,
	TypeBandwidth,
	TypeEncryptionKey,
	TypeAttribute,
}

var errUnknownType = errors.New("unknown type")

// isKnown returns true if t is defined in RFC 4566.
func isKnown(t Type) bool {
	switch t {
	case TypeProtocolVersion,
		TypeOrigin, TypeSessionName,
		TypeSessionInformation, TypeURI,
		TypeEmail, TypePhone,
		TypeConnectionData, TypeBandwidth,
		TypeTiming, TypeRepeatTimes,
		TypeTimeZones, TypeEncryptionKey,
		TypeAttribute, TypeMediaDescription:
		return true
	default:
		return false
	}
}

// isExpected determines if t is expected on pos in s section and returns nil,
// if it is expected and DecodeError if not.
func isExpected(t Type, s section, pos int) error {
	o := getOrdering(s)
	if len(o) > pos {
		for _, expected := range o[pos:] {
			if expected == t {
				return nil
			}
			if isOptional(expected) {
				continue
			}
		}
	}

	// Checking possible section transitions.
	switch s {
	case sectionSession:
		if pos < orderingAfterTime && isExpected(t, sectionTime, 0) == nil {
			return nil
		}
		if isExpected(t, sectionMedia, 0) == nil {
			return nil
		}
	case sectionTime:
		if isExpected(t, sectionSession, orderingAfterTime) == nil {
			return nil
		}
	case sectionMedia:
		if pos != 0 && isExpected(t, sectionMedia, 0) == nil {
			return nil
		}
	}
	if !isKnown(t) {
		return errUnknownType
	}
	// Attribute is known, but out of order.
	msg := fmt.Sprintf("no matches in ordering array at %s[%d] for %s",
		s, pos, t,
	)
	err := newSectionDecodeError(s, msg)
	return errors.Wrapf(err, "field %s is unexpected", t)
}

func getOrdering(s section) ordering {
	switch s {
	case sectionSession:
		return orderingSession
	case sectionMedia:
		return orderingMedia
	case sectionTime:
		return orderingTime
	default:
		panic("BUG: section overflow")
	}
}

func isOptional(t Type) bool {
	switch t {
	case TypeProtocolVersion, TypeOrigin, TypeSessionName:
		return false
	case TypeTiming:
		return false
	case TypeMediaDescription:
		return false
	default:
		return true
	}
}

func isZeroOrMore(t Type) bool {
	switch t {
	case TypeBandwidth, TypeAttribute:
		return true
	default:
		return false
	}
}

func newSectionDecodeError(s section, m string) DecodeError {
	place := fmt.Sprintf("section %s", s)
	return newDecodeError(place, m)
}

func (d *Decoder) decodeKV() (k, v string, err error) {
	var (
		key     []byte
		value   []byte
		isValue bool
	)
	for _, v := range d.v {
		if v == ':' && !isValue {
			isValue = true
			continue
		}
		if isValue {
			value = append(value, v)
		} else {
			key = append(key, v)
		}
	}

	if isValue && len(value) < 1 {
		msg := fmt.Sprintf("attribute without value")
		err := newSectionDecodeError(d.section, msg)
		return "", "", err
	}

	return string(key), string(value), nil
}

func (d *Decoder) decodeTiming(m *Message) error {
	d.sPos = 0
	d.section = sectionTime
	for d.next() {
		if err := isExpected(d.t, d.section, d.sPos); err != nil {
			if canSkip(err) {
				continue
			}
			return errors.Wrap(err, "decode failed")
		}
		if !isZeroOrMore(d.t) {
			d.sPos++
		}
		switch d.t {
		case TypeTiming, TypeRepeatTimes:
			if err := d.decodeField(m); err != nil {
				return errors.Wrap(err, "decode failed")
			}
		default:
			// possible switch to Media or Session description
			d.pos--
			return nil
		}
	}
	return nil
}

func (d *Decoder) decodeMedia(m *Message) error {
	d.sPos = 0
	d.section = sectionMedia
	d.m = Media{}
	for d.next() {
		if err := isExpected(d.t, d.section, d.sPos); err != nil {
			if canSkip(err) {
				continue
			}
			return errors.Wrap(err, "decode failed")
		}
		if d.t == TypeMediaDescription && d.sPos != 0 {
			d.pos--
			break
		}
		if !isZeroOrMore(d.t) {
			d.sPos++
		}
		if err := d.decodeField(m); err != nil {
			return errors.Wrap(err, "failed to decode field")
		}
	}
	m.Medias = append(m.Medias, d.m)
	return nil
}

func (d *Decoder) decodeVersion(m *Message) error {
	n, err := strconv.Atoi(string(d.v))
	if err != nil {
		return errors.Wrap(err, "failed to parse version")
	}
	m.Version = n
	return nil
}

func addAttribute(a Attributes, k, v string) Attributes {
	return append(a, Attribute{Key: k, Value: v})
}

func (d *Decoder) decodeAttribute(m *Message) error {
	k, v, err := d.decodeKV()
	if err != nil {
		return errors.Wrap(err, "failed to decode attribute")
	}
	switch d.section {
	case sectionMedia:
		d.m.Attributes = addAttribute(d.m.Attributes, k, v)
	default:
		m.Attributes = addAttribute(m.Attributes, k, v)
	}
	return nil
}

func (d *Decoder) decodeSessionName(m *Message) error {
	m.Name = string(d.v)
	return nil
}

func (d *Decoder) decodeSessionInfo(m *Message) error {
	if d.section == sectionMedia {
		d.m.Title = string(d.v)
	} else {
		m.Info = string(d.v)
	}
	return nil
}

func (d *Decoder) decodeEmail(m *Message) error {
	m.Email = string(d.v)
	return nil
}

func (d *Decoder) decodePhone(m *Message) error {
	m.Phone = string(d.v)
	return nil
}

func (d *Decoder) decodeURI(m *Message) error {
	m.URI = string(d.v)
	return nil
}

func (d *Decoder) decodeEncryption(m *Message) error {
	k, v, err := d.decodeKV()
	if err != nil {
		return errors.Wrap(err, "failed to decode encryption")
	}
	e := Encryption{
		Key:    v,
		Method: k,
	}
	switch d.section {
	case sectionMedia:
		d.m.Encryption = e
	default:
		m.Encryption = e
	}
	return nil
}

// ErrFailedToDecodeIP means that decoder failed to parse IP.
var ErrFailedToDecodeIP = errors.New("invalid IP")

func decodeIP(dst net.IP, v []byte) (net.IP, error) {
	// ALLOCATIONS: suboptimal.
	ip := net.ParseIP(string(v))
	if ip == nil {
		return dst, ErrFailedToDecodeIP
	}
	return ip, nil
}

func decodeByte(dst []byte) (byte, error) {
	// ALLOCATIONS: suboptimal.
	n, err := strconv.ParseInt(string(dst), 10, 16)
	if err != nil {
		return 0, err
	}
	return byte(n), err
}

func isIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

//nolint:gocognit // TODO: simplify
func (d *Decoder) decodeConnectionData(m *Message) error {
	// c=<nettype> <addrtype> <connection-address>
	var (
		netType           []byte
		addressType       []byte
		connectionAddress []byte
		subField          int
		err               error
	)
	for _, v := range d.v {
		if v == fieldsDelimiter {
			subField++
			continue
		}
		switch subField {
		case 0:
			netType = append(netType, v)
		case 1:
			addressType = append(addressType, v)
		case 2:
			connectionAddress = append(connectionAddress, v)
		default:
			err = d.newFieldError("unexpected subfield count")
			return errors.Wrap(err, "failed to decode connection data")
		}
	}
	if len(netType) == 0 {
		err = d.newFieldError("nettype is empty")
		return errors.Wrap(err, "failed to decode connection data")
	}
	if len(addressType) == 0 {
		err = d.newFieldError("addrtype is empty")
		return errors.Wrap(err, "failed to decode connection data")
	}
	if len(connectionAddress) == 0 {
		err := d.newFieldError("connection-address is empty")
		return errors.Wrap(err, "failed to decode connection data")
	}
	switch d.section {
	case sectionMedia:
		d.m.Connection.AddressType = string(addressType)
		d.m.Connection.NetworkType = string(netType)
	case sectionSession:
		m.Connection.AddressType = string(addressType)
		m.Connection.NetworkType = string(netType)
	}
	// Decoding address.
	// <base multicast address>[/<ttl>]/<number of addresses>
	var (
		base   []byte
		first  []byte
		second []byte
	)
	subField = 0
	for _, v := range connectionAddress {
		if v == '/' {
			subField++
			continue
		}
		switch subField {
		case 0:
			base = append(base, v)
		case 1:
			first = append(first, v)
		case 2:
			second = append(second, v)
		default:
			err = d.newFieldError("unexpected fourth element in address")
			return errors.Wrap(err, "failed to decode connection data")
		}
	}
	switch d.section {
	case sectionMedia:
		d.m.Connection.IP, err = decodeIP(m.Connection.IP, base)
	case sectionSession:
		m.Connection.IP, err = decodeIP(m.Connection.IP, base)
	}
	if err != nil {
		return errors.Wrap(err, "failed to decode connection data")
	}

	var isV4 bool
	switch d.section {
	case sectionMedia:
		isV4 = isIPv4(d.m.Connection.IP)
	case sectionSession:
		isV4 = isIPv4(m.Connection.IP)
	}
	if len(second) > 0 {
		if !isV4 {
			err := d.newFieldError("unexpected TTL for IPv6")
			return errors.Wrap(err, "failed to decode connection data")
		}
		switch d.section {
		case sectionMedia:
			d.m.Connection.TTL, err = decodeByte(first)
		case sectionSession:
			m.Connection.TTL, err = decodeByte(first)
		}
		if err != nil {
			return errors.Wrap(err, "failed to decode connection data")
		}
		switch d.section {
		case sectionMedia:
			d.m.Connection.Addresses, err = decodeByte(second)
		case sectionSession:
			m.Connection.Addresses, err = decodeByte(second)
		}
		if err != nil {
			return errors.Wrap(err, "failed to decode connection data")
		}
	} else if len(first) > 0 {
		if isV4 {
			switch d.section {
			case sectionMedia:
				d.m.Connection.TTL, err = decodeByte(first)
			case sectionSession:
				m.Connection.TTL, err = decodeByte(first)
			}
		} else {
			switch d.section {
			case sectionMedia:
				d.m.Connection.Addresses, err = decodeByte(second)
			case sectionSession:
				m.Connection.Addresses, err = decodeByte(second)
			}
		}
		if err != nil {
			msg := fmt.Sprintf("bad connection data <%s> at <%s>",
				b2s(first), b2s(connectionAddress),
			)
			return errors.Wrap(err, msg)
		}
	}
	return nil
}

func (d *Decoder) decodeBandwidth(m *Message) error {
	k, v, err := d.decodeKV()
	if err != nil {
		return errors.Wrap(err, "failed to decode bandwidth")
	}
	if v == "" {
		msg := "no value specified"
		err := newSectionDecodeError(d.section, msg)
		return errors.Wrap(err, "failed to decode bandwidth")
	}
	var (
		t BandwidthType
		n int
	)
	switch bandWidthType := BandwidthType(k); bandWidthType {
	case BandwidthApplicationSpecific, BandwidthConferenceTotal, BandwidthApplicationSpecificTransportIndependent:
		t = bandWidthType
	default:
		msg := fmt.Sprintf("bad bandwidth type %s", k)
		err = newSectionDecodeError(d.section, msg)
		return errors.Wrap(err, "failed to decode bandwidth")
	}
	if n, err = strconv.Atoi(v); err != nil {
		return errors.Wrap(err, "failed to convert decode bandwidth")
	}
	if d.section == sectionMedia {
		if d.m.Bandwidths == nil {
			d.m.Bandwidths = make(Bandwidths)
		}
		d.m.Bandwidths[t] = n
	} else {
		if m.Bandwidths == nil {
			m.Bandwidths = make(Bandwidths)
		}
		m.Bandwidths[t] = n
	}
	return nil
}

func parseNTP(v []byte) (uint64, error) {
	return strconv.ParseUint(string(v), 10, 64)
}

func (d *Decoder) decodeTimingField(m *Message) error {
	var (
		startV, endV []byte
		isEndV       bool
		err          error
	)
	for _, v := range d.v {
		if v == fieldsDelimiter {
			if isEndV {
				msg := "unexpected second space in timing"
				err = newSectionDecodeError(d.section, msg)
				return errors.Wrap(err, "failed to decode timing")
			}
			isEndV = true
			continue
		}
		if isEndV {
			endV = append(endV, v)
		} else {
			startV = append(startV, v)
		}
	}
	var (
		ntpStart, ntpEnd uint64
	)
	if ntpStart, err = parseNTP(startV); err != nil {
		return errors.Wrap(err, "failed to parse start time")
	}
	if ntpEnd, err = parseNTP(endV); err != nil {
		return errors.Wrap(err, "failed to parse end time")
	}
	t := Timing{}
	t.Start = NTPToTime(ntpStart)
	t.End = NTPToTime(ntpEnd)
	m.Timing = append(m.Timing, t)
	return nil
}

func decodeString(v []byte, s *string) {
	*s = b2s(v)
}

func decodeInt(v []byte, i *int) error {
	var err error
	*i, err = strconv.Atoi(b2s(v))
	return err
}

func decodeInt64(v []byte, i *int64) error {
	var err error
	*i, err = strconv.ParseInt(b2s(v), 10, 64)
	return err
}

func (d *Decoder) subfields() ([][]byte, error) {
	n := bytes.Count(d.v, []byte{fieldsDelimiter})
	result := make([][]byte, n+1)
	subField := 0
	hitSpace := true
	for _, v := range d.v {
		if v == fieldsDelimiter {
			if hitSpace {
				msg := "unexpected second space in subfields"
				return nil, newSectionDecodeError(d.section, msg)
			}
			subField++
			hitSpace = true
			continue
		} else {
			hitSpace = false
		}

		result[subField] = append(result[subField], v)
	}

	return result[:subField+1], nil
}

func (d *Decoder) decodeOrigin(m *Message) error {
	// o=0<username> 1<sess-id> 2<sess-version> 3<nettype> 4<addrtype>
	// 5<unicast-address>
	// ALLOCATIONS: suboptimal
	// CPU: suboptimal
	var (
		err error
	)
	p, err := d.subfields()
	if err != nil {
		return errors.Wrap(err, "failed to decode origin")
	}
	if len(p) != 6 {
		msg := fmt.Sprintf("unexpected subfields count %d != %d", len(p), 6)
		err = newSectionDecodeError(d.section, msg)
		return errors.Wrap(err, "failed to decode origin")
	}
	o := m.Origin
	decodeString(p[0], &o.Username)
	if err = decodeInt64(p[1], &o.SessionID); err != nil {
		return errors.Wrap(err, "failed to decode sess-id")
	}
	if err = decodeInt64(p[2], &o.SessionVersion); err != nil {
		return errors.Wrap(err, "failed to decode sess-version")
	}
	decodeString(p[3], &o.NetworkType)
	decodeString(p[4], &o.AddressType)
	decodeString(p[5], &o.Address)
	m.Origin = o
	return nil
}

func decodeInterval(b []byte, v *time.Duration) error {
	if len(b) == 1 && b[0] == '0' {
		*v = 0
		return nil
	}
	var (
		unit            time.Duration
		noUnitSpecified bool
		val             int
	)
	switch b[len(b)-1] {
	case 'd':
		unit = time.Hour * 24
	case 'h':
		unit = time.Hour
	case 'm':
		unit = time.Minute
	case 's':
		unit = time.Second
	default:
		unit = time.Second
		noUnitSpecified = true
	}
	if !noUnitSpecified {
		if len(b) < 2 {
			err := io.ErrUnexpectedEOF
			return errors.Wrap(err, "unit without value is invalid duration")
		}
		b = b[:len(b)-1]
	}
	if err := decodeInt(b, &val); err != nil {
		return errors.Wrap(err, "unable to decode value")
	}
	*v = time.Duration(val) * unit
	return nil
}

func (d *Decoder) decodeRepeatTimes(m *Message) error {
	// r=0<repeat interval> 1<active duration> 2<offsets from start-time>
	var err error
	if len(m.Timing) < 1 {
		msg := fmt.Sprintf("repeat without timing")
		err = newSectionDecodeError(d.section, msg)
		return errors.Wrap(err, "failed to decode repeat")
	}

	p, err := d.subfields()
	if err != nil {
		return errors.Wrap(err, "failed to decode repeat")
	}
	if len(p) < 3 {
		msg := fmt.Sprintf("unexpected subfields count %d < 3", len(p))
		err = newSectionDecodeError(d.section, msg)
		return errors.Wrap(err, "failed to decode repeat")
	}
	t := m.Timing[len(m.Timing)-1]
	if err = decodeInterval(p[0], &t.Repeat); err != nil {
		return errors.Wrap(err, "failed to decode repeat interval")
	}
	if err = decodeInterval(p[1], &t.Active); err != nil {
		return errors.Wrap(err, "failed to decode active duration")
	}
	var dd time.Duration
	for i, pp := range p[2:] {
		if err = decodeInterval(pp, &dd); err != nil {
			return errors.Wrapf(err, "failed to decode offset %d", i)
		}
		t.Offsets = append(t.Offsets, dd)
	}
	return nil
}

func (d *Decoder) decodeTimeZoneAdjustments(m *Message) error {
	// z=<adjustment time> <offset> <adjustment time> <offset> ....
	p, err := d.subfields()
	if err != nil {
		return errors.Wrap(err, "failed to decode tz-adjustments")
	}
	var (
		adjustment TimeZone
		t          uint64
	)
	if len(p)%2 != 0 {
		msg := fmt.Sprintf("unexpected subfields count %d", len(p))
		err = newSectionDecodeError(d.section, msg)
		return errors.Wrap(err, "failed to decode tz-adjustments")
	}
	for i := 0; i < len(p); i += 2 {
		if t, err = parseNTP(p[i]); err != nil {
			return errors.Wrap(err, "failed to decode adjustment start")
		}
		adjustment.Start = NTPToTime(t)
		if err = decodeInterval(p[i+1], &adjustment.Offset); err != nil {
			return errors.Wrap(err, "failed to decode offset")
		}
		m.TZAdjustments = append(m.TZAdjustments, adjustment)
	}
	return nil
}

func (d *Decoder) decodeMediaDescription(_ *Message) error {
	// m=0<media> 1<port> 2<proto> 3<fmt> ...
	var (
		desc MediaDescription
		err  error
	)
	p, err := d.subfields()
	if err != nil {
		return errors.Wrap(err, "failed to decode media description")
	}
	if len(p) < 3 {
		msg := fmt.Sprintf("unexpected subfields count %d < 3", len(p))
		err = newSectionDecodeError(d.section, msg)
		return errors.Wrap(err, "failed to decode media description")
	}
	decodeString(p[0], &desc.Type)
	// port: port/ports_number
	pp := bytes.Split(p[1], []byte{'/'})
	if err = decodeInt(pp[0], &desc.Port); err != nil {
		return errors.Wrap(err, "failed to decode port")
	}
	if len(pp) > 1 {
		if err = decodeInt(pp[1], &desc.PortsNumber); err != nil {
			return errors.Wrap(err, "failed to decode ports number")
		}
	}
	decodeString(p[2], &desc.Protocol)
	for _, rawFormat := range p[3:] {
		desc.Formats = append(desc.Formats, string(rawFormat))
	}
	d.m.Description = desc
	return nil
}

func (d *Decoder) decodeField(m *Message) error {
	switch d.t {
	case TypeProtocolVersion:
		return d.decodeVersion(m)
	case TypeAttribute:
		return d.decodeAttribute(m)
	case TypeSessionName:
		return d.decodeSessionName(m)
	case TypeSessionInformation:
		return d.decodeSessionInfo(m)
	case TypeEmail:
		return d.decodeEmail(m)
	case TypePhone:
		return d.decodePhone(m)
	case TypeURI:
		return d.decodeURI(m)
	case TypeEncryptionKey:
		return d.decodeEncryption(m)
	case TypeBandwidth:
		return d.decodeBandwidth(m)
	case TypeTiming:
		return d.decodeTimingField(m)
	case TypeConnectionData:
		return d.decodeConnectionData(m)
	case TypeOrigin:
		return d.decodeOrigin(m)
	case TypeRepeatTimes:
		return d.decodeRepeatTimes(m)
	case TypeTimeZones:
		return d.decodeTimeZoneAdjustments(m)
	case TypeMediaDescription:
		return d.decodeMediaDescription(m)
	default:
		// d.t is explicitly checked before calling decodeField,
		// so this code must be unreachable.
		panic("BUG: unexpected filed type in decodeField")
	}
}

func canSkip(err error) bool {
	return errors.Cause(err) == errUnknownType
}

func (d *Decoder) decodeSession(m *Message) error {
	d.sPos = 0
	d.section = sectionSession
	for d.next() {
		if err := isExpected(d.t, d.section, d.sPos); err != nil {
			if canSkip(err) {
				continue
			}
			return errors.Wrap(err, "decode failed")
		}
		if !isZeroOrMore(d.t) {
			d.sPos++
		}
		switch d.t {
		case TypeTiming:
			d.pos--
			oldPosition := d.sPos
			if err := d.decodeTiming(m); err != nil {
				return errors.Wrap(err, "failed to decode timing")
			}
			d.sPos = oldPosition
			d.section = sectionSession
		case TypeMediaDescription:
			d.pos--
			oldPosition := d.sPos
			if err := d.decodeMedia(m); err != nil {
				return errors.Wrap(err, "failed to decode media")
			}
			d.sPos = oldPosition
			d.section = sectionSession
		default:
			if err := d.decodeField(m); err != nil {
				return errors.Wrap(err, "failed to decode field")
			}
		}
	}

	if m.Origin.Address == "" {
		msg := fmt.Sprintf("origin address not set")
		err := newSectionDecodeError(sectionSession, msg)
		return errors.Wrap(err, "failed to decode message")
	}
	if m.Name == "" {
		msg := fmt.Sprintf("session name not set")
		err := newSectionDecodeError(sectionSession, msg)
		return errors.Wrap(err, "failed to decode message")
	}

	return nil
}

// Decode message from session.
func (d *Decoder) Decode(m *Message) error {
	return d.decodeSession(m)
}

// b2s converts byte slice to a string without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b)) // #nosec
}
