package sip

import (
	"bytes"
	"fmt"
	"net"
)

// MessageID MessageID
type MessageID string

// RequestMethod This is syntactic sugar around the string type, so make sure to use
// the Equals method rather than built-in equality, or you'll fall foul of case differences.
// If you're defining your own Method, uppercase is preferred but not compulsory.
type RequestMethod string

// It's nicer to avoid using raw strings to represent methods, so the following standard
// method names are defined here as constants for convenience.
const (
	INVITE   RequestMethod = "INVITE"
	ACK      RequestMethod = "ACK"
	CANCEL   RequestMethod = "CANCEL"
	BYE      RequestMethod = "BYE"
	REGISTER RequestMethod = "REGISTER"
	OPTIONS  RequestMethod = "OPTIONS"
	// SUBSCRIBE RequestMethod = "SUBSCRIBE"
	// NOTIFY  RequestMethod = "NOTIFY"
	// REFER   RequestMethod = "REFER"
	INFO    RequestMethod = "INFO"
	MESSAGE RequestMethod = "MESSAGE"
)

// Message introduces common SIP message RFC 3261 - 7.
type Message interface {
	MessageID() MessageID

	Clone() Message
	// Start line returns message start line.
	StartLine() string
	// String returns string representation of SIP message in RFC 3261 form.
	String() string
	// SipVersion returns SIP protocol version.
	SipVersion() string
	// SetSipVersion sets SIP protocol version.
	SetSipVersion(version string)

	// Headers returns all message headers.
	Headers() []Header
	// GetHeaders returns slice of headers of the given type.
	GetHeaders(name string) []Header
	// AppendHeader appends header to message.
	AppendHeader(header Header)
	// PrependHeader prepends header to message.
	RemoveHeader(name string)

	// Body returns message body.
	Body() []byte
	// SetBody sets message body.
	SetBody(body []byte, setContentLength bool)

	/* Helper getters for common headers */
	// CallID returns 'Call-ID' header.
	CallID() (*CallID, bool)
	// Via returns the top 'Via' header field.
	Via() (ViaHeader, bool)
	// ViaHop returns the first segment of the top 'Via' header.
	ViaHop() (*ViaHop, bool)
	// From returns 'From' header field.
	From() (*FromHeader, bool)
	// To returns 'To' header field.
	To() (*ToHeader, bool)
	// CSeq returns 'CSeq' header field.
	CSeq() (*CSeq, bool)
	ContentLength() (*ContentLength, bool)
	ContentType() (*ContentType, bool)
	Contact() (*ContactHeader, bool)

	Transport() string
	Source() net.Addr
	SetSource(src net.Addr)
	Destination() net.Addr
	SetDestination(dest net.Addr)

	IsCancel() bool
	IsAck() bool
}

type message struct {
	// message headers
	*headers
	messID       MessageID
	sipVersion   string
	body         []byte
	source, dest net.Addr
	startLine    func() string
}

// MessageID MessageID
func (msg *message) MessageID() MessageID {
	return msg.messID
}

// StartLine StartLine
func (msg *message) StartLine() string {
	return msg.startLine()
}

func (msg *message) String() string {
	var buffer bytes.Buffer

	// write message start line
	buffer.WriteString(msg.StartLine() + "\r\n")
	// Write the headers.
	buffer.WriteString(msg.headers.String())
	// message body
	buffer.WriteString("\r\n")
	buffer.Write(msg.Body())

	return buffer.String()
}

// SipVersion SipVersion
func (msg *message) SipVersion() string {
	return msg.sipVersion
}

// SetSipVersion SetSipVersion
func (msg *message) SetSipVersion(version string) {
	msg.sipVersion = version
}

// Body Body
func (msg *message) Body() []byte {
	return msg.body
}

// SetBody sets message body, calculates it length and add 'Content-Length' header.
func (msg *message) SetBody(body []byte, setContentLength bool) {
	msg.body = body
	if setContentLength {
		hdrs := msg.GetHeaders("Content-Length")
		if len(hdrs) == 0 {
			length := ContentLength(len(body))
			msg.AppendHeader(&length)
		} else {
			length := ContentLength(len(body))
			hdrs[0] = &length
		}
	}
}

// Transport  Transport
func (msg *message) Transport() string {
	if viaHop, ok := msg.ViaHop(); ok {
		return viaHop.Transport
	}
	return DefaultProtocol
}

// Source Source
func (msg *message) Source() net.Addr {
	return msg.source
}

// SetSource SetSource
func (msg *message) SetSource(src net.Addr) {
	msg.source = src
}

// Destination Destination
func (msg *message) Destination() net.Addr {
	return msg.dest
}

// SetDestination SetDestination
func (msg *message) SetDestination(dest net.Addr) {
	msg.dest = dest
}

// URI  A SIP or SIPS URI, including all params and URI header params.
// noinspection GoNameStartsWithPackageName
type URI struct {
	// True if and only if the URI is a SIPS URI.
	FIsEncrypted bool

	// The user part of the URI: the 'joe' in sip:joe@bloggs.com
	// This is a pointer, so that URIs without a user part can have 'nil'.
	FUser MaybeString

	// The password field of the URI. This is represented in the URI as joe:hunter2@bloggs.com.
	// Note that if a URI has a password field, it *must* have a user field as well.
	// This is a pointer, so that URIs without a password field can have 'nil'.
	// Note that RFC 3261 strongly recommends against the use of password fields in SIP URIs,
	// as they are fundamentally insecure.
	FPassword MaybeString

	// The host part of the URI. This can be a domain, or a string representation of an IP address.
	FHost string

	// The port part of the URI. This is optional, and so is represented here as a pointer type.
	FPort *Port

	// Any parameters associated with the URI.
	// These are used to provide information about requests that may be constructed from the URI.
	// (For more details, see RFC 3261 section 19.1.1).
	// These appear as a semicolon-separated list of key=value pairs following the host[:port] part.
	FUriParams Params

	// Any headers to be included on requests constructed from this URI.
	// These appear as a '&'-separated list at the end of the URI, introduced by '?'.
	// Although the values of the map are MaybeStrings, they will never be NoString in practice as the parser
	// guarantees to not return blank values for header elements in SIP URIs.
	// You should not set the values of headers to NoString.
	FHeaders Params
}

// User User
func (uri *URI) User() MaybeString {
	return uri.FUser
}

// Host Host
func (uri *URI) Host() string {
	return uri.FHost
}

// SetHost SetHost
func (uri *URI) SetHost(host string) {
	uri.FHost = host
}

// Generates the string representation of a SipUri struct.
func (uri *URI) String() string {
	var buffer bytes.Buffer

	// Compulsory protocol identifier.
	if uri.FIsEncrypted {
		buffer.WriteString("sips")
		buffer.WriteString(":")
	} else {
		buffer.WriteString("sip")
		buffer.WriteString(":")
	}

	// Optional userinfo part.
	if user, ok := uri.FUser.(String); ok && user.String() != "" {
		buffer.WriteString(uri.FUser.String())
		if pass, ok := uri.FPassword.(String); ok && pass.String() != "" {
			buffer.WriteString(":")
			buffer.WriteString(pass.String())
		}
		buffer.WriteString("@")
	}

	// Compulsory hostname.
	buffer.WriteString(uri.FHost)

	// Optional port number.
	if uri.FPort != nil {
		buffer.WriteString(fmt.Sprintf(":%d", *uri.FPort))
	}

	if (uri.FUriParams != nil) && uri.FUriParams.Length() > 0 {
		buffer.WriteString(";")
		buffer.WriteString(uri.FUriParams.ToString(';'))
	}

	if (uri.FHeaders != nil) && uri.FHeaders.Length() > 0 {
		buffer.WriteString("?")
		buffer.WriteString(uri.FHeaders.ToString('&'))
	}

	return buffer.String()
}

// Clone the Sip URI.
func (uri *URI) Clone() *URI {
	var newURI *URI
	if uri == nil {
		return newURI
	}

	newURI = &URI{
		FIsEncrypted: uri.FIsEncrypted,
		FUser:        uri.FUser,
		FPassword:    uri.FPassword,
		FHost:        uri.FHost,
		FUriParams:   cloneWithNil(uri.FUriParams),
		FHeaders:     cloneWithNil(uri.FHeaders),
	}
	if uri.FPort != nil {
		newURI.FPort = uri.FPort.Clone()
	}
	return newURI
}

// Equals Determine if the SIP URI is equal to the specified URI according to the rules laid down in RFC 3261 s. 19.1.4.
// TODO: The Equals method is not currently RFC-compliant; fix this!
func (uri *URI) Equals(val interface{}) bool {
	otherPtr, ok := val.(*URI)
	if !ok {
		return false
	}

	if uri == otherPtr {
		return true
	}
	if uri == nil && otherPtr != nil || uri != nil && otherPtr == nil {
		return false
	}

	other := *otherPtr
	result := uri.FIsEncrypted == other.FIsEncrypted &&
		uri.FUser == other.FUser &&
		uri.FPassword == other.FPassword &&
		uri.FHost == other.FHost &&
		Uint16PtrEq((*uint16)(uri.FPort), (*uint16)(other.FPort))

	if !result {
		return false
	}

	if uri.FUriParams != otherPtr.FUriParams {
		if uri.FUriParams == nil {
			result = result && otherPtr.FUriParams != nil
		} else {
			result = result && uri.FUriParams.Equals(otherPtr.FUriParams)
		}
	}

	if uri.FHeaders != otherPtr.FHeaders {
		if uri.FHeaders == nil {
			result = result && otherPtr.FHeaders != nil
		} else {
			result = result && uri.FHeaders.Equals(otherPtr.FHeaders)
		}
	}

	return result
}

func cloneWithNil(params Params) Params {
	if params == nil {
		return NewParams()
	}
	return params.Clone()
}
