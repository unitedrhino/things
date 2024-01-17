package sdp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func appendSpace(v []byte) []byte {
	return appendRune(v, fieldsDelimiter)
}

func appendInt(v []byte, i int) []byte {
	if i == 0 {
		return appendRune(v, '0')
	}
	if i > 0 {
		return appendUint(v, i)
	}
	// ALLOCATIONS: suboptimal.
	return append(v, strconv.Itoa(i)...)
}

func appendInt64(v []byte, i int64) []byte {
	// ALLOCATIONS: suboptimal.
	return append(v, strconv.FormatInt(i, 10)...)
}

// AppendUint appends n to dst and returns the extended dst.
func appendUint(dst []byte, n int) []byte {
	if n < 0 {
		panic("BUG: n should be positive")
	}
	var b [20]byte
	buf := b[:]
	i := len(buf)
	var q int
	for n >= 10 {
		i--
		q = n / 10
		buf[i] = '0' + byte(n-q*10)
		n = q
	}
	i--
	buf[i] = '0' + byte(n)
	dst = append(dst, buf[i:]...)
	return dst
}

func appendByte(v []byte, i byte) []byte {
	if i == 0 {
		return appendRune(v, '0')
	}
	return appendUint(v, int(i))
}

func appendJoinStrings(b []byte, v ...string) []byte {
	last := len(v) - 1
	for i, vv := range v {
		b = append(b, vv...)
		if i != last {
			b = appendSpace(b)
		}
	}
	return b
}

func appendIP(b []byte, ip net.IP) []byte {
	switch ipV4 := ip.To4(); ipV4 {
	case nil:
		// ALLOCATIONS: suboptimal.
		return append(b, strings.ToUpper(ip.String())...)
	default:
		return appendIPv4(b, ipV4)
	}
}

func appendIPv4(dst []byte, ip net.IP) []byte {
	dst = appendUint(dst, int(ip[0]))
	for i := 1; i < 4; i++ {
		dst = append(dst, '.')
		dst = appendUint(dst, int(ip[i]))
	}
	return dst
}

// AddRaw appends k=v to Session.
func (s Session) AddRaw(k rune, v string) Session {
	return s.appendString(Type(k), v)
}

// AddLine appends t=v to Session.
func (s Session) AddLine(t Type, v string) Session {
	return s.appendString(t, v)
}

// AddVersion appends Version field to Session.
func (s Session) AddVersion(version int) Session {
	v := make([]byte, 0, 64)
	v = appendInt(v, version)
	return s.append(TypeProtocolVersion, v)
}

// AddPhone appends Phone Address field to Session.
func (s Session) AddPhone(phone string) Session {
	return s.appendString(TypePhone, phone)
}

// AddEmail appends Email Address field to Session.
func (s Session) AddEmail(email string) Session {
	return s.appendString(TypeEmail, email)
}

// AddConnectionData appends Connection Data field to Session
// using ConnectionData struct with sensible defaults.
func (s Session) AddConnectionData(data ConnectionData) Session {
	v := make([]byte, 0, 512)
	v = append(v, data.getNetworkType()...)
	v = appendSpace(v)
	v = append(v, data.getAddressType()...)
	v = appendSpace(v)
	v = data.appendAddress(v)
	return s.append(TypeConnectionData, v)
}

// AddConnectionDataIP appends Connection Data field using only ip address.
func (s Session) AddConnectionDataIP(ip net.IP) Session {
	return s.AddConnectionData(ConnectionData{
		IP: ip,
	})
}

// AddSessionName appends Session Name field to Session.
func (s Session) AddSessionName(name string) Session {
	return s.appendString(TypeSessionName, name)
}
func (s Session) AddSSRC(d string) Session {
	return s.appendString(TypeSSRC, d)
}

// AddSessionInfo appends Session Information field to Session.
func (s Session) AddSessionInfo(info string) Session {
	return s.appendString(TypeSessionInformation, info)
}

// AddURI appends Uniform Resource Identifier field to Session.
func (s Session) AddURI(uri string) Session {
	return s.appendString(TypeURI, uri)
}

// ConnectionData is representation for Connection Data field.
// Only IP field is required. NetworkType and AddressType have
// sensible defaults.
type ConnectionData struct {
	NetworkType string // <nettype>
	AddressType string // <addrtype>
	IP          net.IP // <base multicast address>
	TTL         byte   // <ttl>
	Addresses   byte   // <number of addresses>
}

// Blank determines if ConnectionData is blank value.
func (c ConnectionData) Blank() bool {
	return c.Equal(ConnectionData{})
}

// Equal returns c == b.
func (c ConnectionData) Equal(b ConnectionData) bool {
	if c.NetworkType != b.NetworkType {
		return false
	}
	if c.AddressType != b.AddressType {
		return false
	}
	if !c.IP.Equal(b.IP) {
		return false
	}
	if c.TTL != b.TTL {
		return false
	}
	if c.Addresses != b.Addresses {
		return false
	}
	return true
}

const (
	addrTypeIPv4        = "IP4"
	addrTypeIPv6        = "IP6"
	networkTypeInternet = "IN"
	attributesDelimiter = ':'
)

func (c ConnectionData) getNetworkType() string {
	return getDefault(c.NetworkType, networkTypeInternet)
}

// getAddressType returns Address Type ("addrtype") for ip,
// using addressType as default value if present.
func getAddressType(addr, addressType string) string {
	if addressType != "" {
		return addressType
	}
	for _, s := range addr {
		if s == ':' {
			return addrTypeIPv6
		}
		if s == '.' {
			return addrTypeIPv4
		}
	}
	return getAddressTypeIP(net.ParseIP(addr), addressType)
}

// getAddressType returns Address Type ("addrtype") for ip,
// using addressType as default value if present.
func getAddressTypeIP(ip net.IP, addressType string) string {
	if addressType != "" {
		return addressType
	}
	if ip == nil {
		return addrTypeIPv4
	}
	switch ip.To4() {
	case nil:
		return addrTypeIPv6
	default:
		return addrTypeIPv4
	}
}

func (c ConnectionData) getAddressType() string {
	return getAddressTypeIP(c.IP, c.AddressType)
}

// ConnectionAddress formats <connection-address> sub-field.
func (c ConnectionData) ConnectionAddress() string {
	// <base multicast address>[/<ttl>]/<number of addresses>
	// ALLOCATIONS: suboptimal. Use appendAddress.
	var address = strings.ToUpper(c.IP.String())
	if c.TTL > 0 {
		address += fmt.Sprintf("/%d", c.TTL)
	}
	if c.Addresses > 0 {
		address += fmt.Sprintf("/%d", c.Addresses)
	}
	return address
}

func (c ConnectionData) String() string {
	return fmt.Sprintf("%s %s %s",
		c.getAddressType(), c.getAddressType(), c.ConnectionAddress(),
	)
}

func (c ConnectionData) appendAddress(v []byte) []byte {
	v = appendIP(v, c.IP)
	if c.TTL > 0 {
		v = appendRune(v, '/')
		v = appendByte(v, c.TTL)
	}
	if c.Addresses > 0 {
		v = appendRune(v, '/')
		v = appendByte(v, c.Addresses)
	}
	return v
}

// Origin is field defined in RFC4566 5.2.
// See https://tools.ietf.org/html/rfc4566#section-5.2.
type Origin struct {
	Username       string // <username>
	SessionID      int64  // <sess-id>
	SessionVersion int64  // <sess-version>
	NetworkType    string // <nettype>
	AddressType    string // <addrtype>
	Address        string // <unicast-address>
}

func (o *Origin) getNetworkType() string {
	return getDefault(o.NetworkType, networkTypeInternet)
}

func (o *Origin) getAddressType() string {
	return getAddressType(o.Address, o.AddressType)
}

// Equal returns b == o.
func (o *Origin) Equal(b Origin) bool {
	if o.Username != b.Username {
		return false
	}
	if o.SessionID != b.SessionID {
		return false
	}
	if o.SessionVersion != b.SessionVersion {
		return false
	}
	if o.NetworkType != b.NetworkType {
		return false
	}
	if o.AddressType != b.AddressType {
		return false
	}
	if o.Address != b.Address {
		return false
	}
	return true
}

// AddOrigin appends Origin field to Session.
func (s Session) AddOrigin(o Origin) Session {
	v := make([]byte, 0, 2048)
	v = appendSpace(append(v, o.Username...))
	v = appendSpace(appendInt64(v, o.SessionID))
	v = appendSpace(appendInt64(v, o.SessionVersion))
	v = appendSpace(append(v, o.getNetworkType()...))
	v = appendSpace(append(v, o.getAddressType()...))
	v = append(v, o.Address...)
	return s.append(TypeOrigin, v)
}

const (
	// ntpDelta is seconds from Jan 1, 1900 to Jan 1, 1970.
	ntpDelta = 2208988800
)

// TimeToNTP converts time.Time to NTP timestamp with special case for Zero
// time, that is interpreted as 0 timestamp.
func TimeToNTP(t time.Time) uint64 {
	if t.IsZero() {
		return 0
	}
	return uint64(t.Unix())
	// return uint64(t.Unix()) + ntpDelta
}

// NTPToTime converts NTP timestamp to time.Time with special case for Zero
// time, that is interpreted as 0 timestamp.
func NTPToTime(v uint64) time.Time {
	if v == 0 {
		return time.Time{}
	}
	return time.Unix(int64(v-ntpDelta), 0)
}

func appendUint64(b []byte, v uint64) []byte {
	return strconv.AppendUint(b, v, 10)
}

// AddTiming appends Timing field to Session. Both start and end can be zero.
func (s Session) AddTiming(start, end time.Time) Session {
	v := make([]byte, 0, 256)
	v = appendUint64(v, TimeToNTP(start))
	v = appendSpace(v)
	v = appendUint64(v, TimeToNTP(end))
	return s.append(TypeTiming, v)
}

// AddTimingNTP appends Timing field to Session with NTP timestamps as input.
// It is just wrapper for AddTiming and NTPToTime.
func (s Session) AddTimingNTP(start, end uint64) Session {
	return s.AddTiming(NTPToTime(start), NTPToTime(end))
}

// AddAttribute appends Attribute field to Session in a=<attribute>:<value>"
// form. If len(values) > 1, then "<value>" is "<val1> <val2> ... <valn>",
// and if len(values) == 0, then AddFlag method is used in "a=<flag>" form.
func (s Session) AddAttribute(attribute string, values ...string) Session {
	if len(values) == 0 {
		return s.AddFlag(attribute)
	}
	v := make([]byte, 0, 512)
	v = append(v, attribute...)
	v = appendRune(v, attributesDelimiter)
	v = appendJoinStrings(v, values...)
	return s.append(TypeAttribute, v)
}

// AddFlag appends Attribute field to Session in "a=<flag>" form.
func (s Session) AddFlag(attribute string) Session {
	v := make([]byte, 0, 256)
	v = append(v, attribute...)
	return s.append(TypeAttribute, v)
}

// BandwidthType is <bwtype> sub-field of Bandwidth field.
type BandwidthType string

// Possible values for <bwtype> defined in section 5.8.
const (
	BandwidthConferenceTotal     BandwidthType = "CT"
	BandwidthApplicationSpecific BandwidthType = "AS"
	// defined in RFC 3890
	BandwidthApplicationSpecificTransportIndependent BandwidthType = "TIAS"
)

// AddBandwidth appends Bandwidth field to Session.
func (s Session) AddBandwidth(t BandwidthType, bandwidth int) Session {
	v := make([]byte, 0, 128)
	v = append(v, string(t)...)
	v = appendRune(v, ':')
	v = appendInt(v, bandwidth)
	return s.append(TypeBandwidth, v)
}

type durationUnit struct {
	d time.Duration
	r rune
}

func (u durationUnit) append(v []byte, b time.Duration) []byte {
	v = appendInt(v, int(b/u.d))
	return appendRune(v, u.r)
}

var durationUnits = [...]durationUnit{
	{time.Hour * 24, 'd'},
	{time.Hour, 'h'},
	{time.Minute, 'm'},
}

func appendIntervalCompact(b []byte, d time.Duration) []byte {
	if d == 0 {
		return appendRune(b, '0')
	}
	for _, unit := range durationUnits {
		if d%unit.d == 0 {
			return unit.append(b, d)
		}
	}
	return appendInt(b, int(d.Seconds()))
}

func appendInterval(b []byte, d time.Duration, compact bool) []byte {
	if d == 0 {
		return appendRune(b, '0')
	}
	if compact {
		return appendIntervalCompact(b, d)
	}
	return appendInt(b, int(d.Seconds()))
}

func (s Session) addRepeatTimes(compact bool, interval, duration time.Duration,
	offsets ...time.Duration) Session {
	v := make([]byte, 0, 256)
	v = appendSpace(appendInterval(v, interval, compact))
	v = appendSpace(appendInterval(v, duration, compact))
	for i, offset := range offsets {
		v = appendInterval(v, offset, compact)
		if i != len(offsets)-1 {
			v = appendSpace(v)
		}
	}
	return s.append(TypeRepeatTimes, v)
}

// AddRepeatTimes appends Repeat Times field to Session.
func (s Session) AddRepeatTimes(interval, duration time.Duration,
	offsets ...time.Duration) Session {
	return s.addRepeatTimes(false, interval, duration, offsets...)
}

// AddRepeatTimesCompact appends Repeat Times field to Session using "compact"
// syntax.
func (s Session) AddRepeatTimesCompact(interval, duration time.Duration,
	offsets ...time.Duration) Session {
	return s.addRepeatTimes(true, interval, duration, offsets...)
}

// MediaDescription represents Media Description field value.
type MediaDescription struct {
	Type        string
	Port        int
	PortsNumber int
	Protocol    string
	Formats     []string
}

// Equal returns true if b equals to m.
func (m MediaDescription) Equal(b MediaDescription) bool {
	if m.Type != b.Type {
		return false
	}
	if m.Port != b.Port {
		return false
	}
	if m.PortsNumber != b.PortsNumber {
		return false
	}
	if m.Protocol != b.Protocol {
		return false
	}
	if len(m.Formats) != len(b.Formats) {
		return false
	}
	for i := range m.Formats {
		if m.Formats[i] != b.Formats[i] {
			return false
		}
	}
	return true
}

// AddMediaDescription appends Media Description field to Session.
func (s Session) AddMediaDescription(m MediaDescription) Session {
	v := make([]byte, 0, 512)
	v = appendSpace(append(v, m.Type...))
	v = appendInt(v, m.Port)
	if m.PortsNumber != 0 {
		v = appendRune(v, '/')
		v = appendInt(v, m.PortsNumber)
	}
	v = appendSpace(v)
	v = appendSpace(append(v, m.Protocol...))
	for i := range m.Formats {
		v = append(v, m.Formats[i]...)
		if i != len(m.Formats)-1 {
			v = appendRune(v, fieldsDelimiter)
		}
	}
	return s.append(TypeMediaDescription, v)
}

// AddEncryption appends Encryption and is shorthand for AddEncryptionKey.
func (s Session) AddEncryption(e Encryption) Session {
	return s.AddEncryptionKey(e.Method, e.Key)
}

// AddEncryptionKey appends Encryption Key field with method and key in
// "k=<method>:<encryption key>" format to Session.
func (s Session) AddEncryptionKey(method, key string) Session {
	if key == "" {
		return s.AddEncryptionMethod(method)
	}
	v := make([]byte, 0, 512)
	v = append(v, method...)
	v = appendRune(v, attributesDelimiter)
	v = append(v, key...)
	return s.append(TypeEncryptionKey, v)
}

// AddEncryptionMethod appends Encryption Key field with only method in
// "k=<method>" format to Session.
func (s Session) AddEncryptionMethod(method string) Session {
	return s.appendString(TypeEncryptionKey, method)
}

// TimeZone is representation of <adjustment time> <offset> pair.
type TimeZone struct {
	Start  time.Time
	Offset time.Duration
}

func (t TimeZone) appendInterval(v []byte) []byte {
	return appendIntervalCompact(v, t.Offset)
}

func (t TimeZone) append(v []byte) []byte {
	v = appendUint64(v, TimeToNTP(t.Start))
	v = appendSpace(v)
	return t.appendInterval(v)
}

// AddTimeZones append TimeZones field to Session.
func (s Session) AddTimeZones(zones ...TimeZone) Session {
	v := make([]byte, 0, 512)
	for i, zone := range zones {
		v = zone.append(v)
		if i != len(zones)-1 {
			v = appendSpace(v)
		}
	}
	return s.append(TypeTimeZones, v)
}

func getDefault(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
