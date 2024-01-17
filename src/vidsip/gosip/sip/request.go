package sip

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"net"
	//"github.com/gofrs/uuid"
)

// Request Request
type Request struct {
	message
	method    RequestMethod
	recipient *URI
}

// NewRequest NewRequest
func NewRequest(
	messID MessageID,
	method RequestMethod,
	recipient *URI,
	sipVersion string,
	hdrs []Header,
	body []byte,
) *Request {
	req := new(Request)
	if messID == "" {
		tmpUUID, _ := uuid.GenerateUUID()
		req.messID = MessageID(tmpUUID)
	} else {
		req.messID = messID
	}
	req.SetSipVersion(sipVersion)
	req.startLine = req.StartLine
	req.headers = newHeaders(hdrs)
	req.SetMethod(method)
	req.SetRecipient(recipient)

	if len(body) != 0 {
		req.SetBody(body, true)
	}
	return req
}

// NewRequestFromResponse NewRequestFromResponse
func NewRequestFromResponse(method RequestMethod, inviteResponse *Response) *Request {
	contact, _ := inviteResponse.Contact()
	ackRequest := NewRequest(
		inviteResponse.MessageID(),
		method,
		contact.Address,
		inviteResponse.SipVersion(),
		[]Header{},
		[]byte{},
	)

	CopyHeaders("Via", inviteResponse, ackRequest)
	viaHop, _ := ackRequest.ViaHop()
	// update branch, 2xx ACK is separate Tx
	viaHop.Params.Add("branch", String{Str: GenerateBranch()})

	if len(inviteResponse.GetHeaders("Route")) > 0 {
		CopyHeaders("Route", inviteResponse, ackRequest)
	} else {
		for _, h := range inviteResponse.GetHeaders("Record-Route") {
			uris := make([]*URI, 0)
			for _, u := range h.(*RecordRouteHeader).Addresses {
				uris = append(uris, u.Clone())
			}
			ackRequest.AppendHeader(&RouteHeader{
				Addresses: uris,
			})
		}
	}

	CopyHeaders("From", inviteResponse, ackRequest)
	CopyHeaders("To", inviteResponse, ackRequest)
	CopyHeaders("Call-ID", inviteResponse, ackRequest)
	cseq, _ := inviteResponse.CSeq()
	cseq.MethodName = method
	cseq.SeqNo++
	ackRequest.AppendHeader(cseq)
	ackRequest.SetSource(inviteResponse.Destination())
	ackRequest.SetDestination(inviteResponse.Source())
	return ackRequest
}

// StartLine returns Request Line - RFC 2361 7.1.
func (req *Request) StartLine() string {
	var buffer bytes.Buffer

	// Every SIP request starts with a Request Line - RFC 2361 7.1.
	buffer.WriteString(
		fmt.Sprintf(
			"%s %s %s",
			string(req.method),
			req.Recipient(),
			req.SipVersion(),
		),
	)

	return buffer.String()
}

// Method Method
func (req *Request) Method() RequestMethod {
	return req.method
}

// SetMethod SetMethod
func (req *Request) SetMethod(method RequestMethod) {
	req.method = method
}

// Recipient Recipient
func (req *Request) Recipient() *URI {
	return req.recipient
}

// SetRecipient SetRecipient
func (req *Request) SetRecipient(recipient *URI) {
	req.recipient = recipient
}

// IsInvite IsInvite
func (req *Request) IsInvite() bool {
	return req.Method() == INVITE
}

// IsAck IsAck
func (req *Request) IsAck() bool {
	return req.Method() == ACK
}

// IsCancel IsCancel
func (req *Request) IsCancel() bool {
	return req.Method() == CANCEL
}

// Source Source
func (req *Request) Source() net.Addr {
	return req.source
}

// SetSource SetSource
func (req *Request) SetSource(src net.Addr) {
	req.source = src
}

// Destination Destination
func (req *Request) Destination() net.Addr {
	return req.dest
}

// SetDestination SetDestination
func (req *Request) SetDestination(dest net.Addr) {
	req.dest = dest
}

// Clone Clone
func (req *Request) Clone() Message {
	return NewRequest(
		"",
		req.Method(),
		req.Recipient().Clone(),
		req.SipVersion(),
		req.headers.CloneHeaders(),
		req.Body(),
	)
}
