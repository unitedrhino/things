package sip

import (
	"bytes"
	"fmt"
	"gitee.com/i-Things/share/utils"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// The whitespace characters recognised by the Augmented Backus-Naur Form syntax
// that SIP uses (RFC 3261 S.25).
const abnfWs = " \t"

// The maximum permissible CSeq number in a SIP message (2**31 - 1).
// C.f. RFC 3261 S. 8.1.1.5.
const maxCseq = 2147483647

// A Parser converts the raw bytes of SIP messages into core.Message objects.
// It allows
type Parser interface {
	// Implements io.Writer. Queues the given bytes to be parsed.
	// If the parser has terminated due to a previous fatal error, it will return n=0 and an appropriate error.
	// Otherwise, it will return n=len(p) and err=nil.
	// Note that err=nil does not indicate that the data provided is valid - simply that the data was successfully queued for parsing.
	Write(p []byte) (n int, err error)
	// Register a custom header parser for a particular header type.
	// This will overwrite any existing registered parser for that header type.
	// If a parser is not available for a header type in a message, the parser will produce a core.GenericHeader struct.
	SetHeaderParser(headerName string, headerParser HeaderParser)

	Stop()

	String() string
	// Reset resets parser state
	Reset()

	ParseHeader(headerText string) (headers []Header, err error)
}

// A HeaderParser is any function that turns raw header data into one or more Header objects.
// The HeaderParser will receive arguments of the form ("max-forwards", "70").
// It should return a slice of headers, which should have length > 1 unless it also returns an error.
type HeaderParser func(headerName string, headerData string) ([]Header, error)

var defaultHeaderParsers = map[string]HeaderParser{
	"to":             parseAddressHeader,
	"t":              parseAddressHeader,
	"from":           parseAddressHeader,
	"f":              parseAddressHeader,
	"contact":        parseAddressHeader,
	"m":              parseAddressHeader,
	"call-id":        parseCallID,
	"cseq":           parseCSeq,
	"via":            parseViaHeader,
	"v":              parseViaHeader,
	"max-forwards":   parseMaxForwards,
	"content-length": parseContentLength,
	"l":              parseContentLength,
	"expires":        parseExpires,
	"user-agent":     parseUserAgent,
	"allow":          parseAllow,
	"content-type":   parseContentType,
	"accept":         parseAccept,
	"c":              parseContentType,
	// "require":        parseRequire,
	"supported":    parseSupported,
	"route":        parseRouteHeader,
	"record-route": parseRecordRouteHeader,
}

// Parse a To, From or Contact header line, producing one or more logical SipHeaders.
func parseAddressHeader(headerName string, headerText string) (
	headers []Header, err error) {
	switch headerName {
	case "to", "from", "contact", "t", "f", "m":
		var displayNames []MaybeString
		var uris []*URI
		var paramSets []Params

		// Perform the actual parsing. The rest of this method is just typeclass bookkeeping.
		displayNames, uris, paramSets, err = ParseAddressValues(headerText)

		if err != nil {
			return
		}
		if len(displayNames) != len(uris) || len(uris) != len(paramSets) {
			// This shouldn't happen unless ParseAddressValues is bugged.
			err = fmt.Errorf("internal parser error: parsed param mismatch. "+
				"%d display names, %d uris and %d param sets "+
				"in %s",
				len(displayNames), len(uris), len(paramSets),
				headerText)
			return
		}

		// Build a slice of headers of the appropriate kind, populating them with the values parsed above.
		// It is assumed that all headers returned by ParseAddressValues are of the same kind,
		// although we do not check for this below.
		for idx := 0; idx < len(displayNames); idx++ {
			var header Header
			if headerName == "to" || headerName == "t" {
				if idx > 0 {
					// Only a single To header is permitted in a SIP message.
					return nil,
						fmt.Errorf("multiple to: headers in message:\n%s: %s",
							headerName, headerText)
				}

				toHeader := ToHeader{
					DisplayName: displayNames[idx],
					Address:     uris[idx],
					Params:      paramSets[idx],
				}
				header = &toHeader
			} else if headerName == "from" || headerName == "f" {
				if idx > 0 {
					// Only a single From header is permitted in a SIP message.
					return nil,
						fmt.Errorf(
							"multiple from: headers in message:\n%s: %s",
							headerName,
							headerText,
						)
				}

				fromHeader := FromHeader{
					DisplayName: displayNames[idx],
					Address:     uris[idx],
					Params:      paramSets[idx],
				}
				header = &fromHeader
			} else if headerName == "contact" || headerName == "m" {
				contactHeader := ContactHeader{
					DisplayName: displayNames[idx],
					Address:     uris[idx],
					Params:      paramSets[idx],
				}
				header = &contactHeader
			}

			headers = append(headers, header)
		}
	}

	return
}

// Parse a string representation of a CSeq header, returning a slice of at most one CSeq.
func parseCSeq(headerName string, headerText string) (
	headers []Header, err error) {
	var cseq CSeq

	parts := SplitByWhitespace(headerText)
	if len(parts) != 2 {
		err = fmt.Errorf(
			"CSeq field should have precisely one whitespace section: '%s'",
			headerText,
		)
		return
	}

	var seqno uint64
	seqno, err = strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return
	}

	if seqno > maxCseq {
		err = fmt.Errorf("invalid CSeq %d: exceeds maximum permitted value "+
			"2**31 - 1", seqno)
		return
	}

	cseq.SeqNo = uint32(seqno)
	cseq.MethodName = RequestMethod(strings.TrimSpace(parts[1]))

	if strings.Contains(string(cseq.MethodName), ";") {
		err = fmt.Errorf("unexpected ';' in CSeq body: %s", headerText)
		return
	}

	headers = []Header{&cseq}

	return
}

// Parse a string representation of a Call-ID header, returning a slice of at most one CallID.
func parseCallID(headerName string, headerText string) (
	headers []Header, err error) {
	headerText = strings.TrimSpace(headerText)
	var callID = CallID(headerText)

	if strings.ContainsAny(string(callID), abnfWs) {
		err = fmt.Errorf("unexpected whitespace in CallID header body '%s'", headerText)
		return
	}
	if strings.Contains(string(callID), ";") {
		err = fmt.Errorf("unexpected semicolon in CallID header body '%s'", headerText)
		return
	}
	if len(string(callID)) == 0 {
		err = fmt.Errorf("empty Call-ID body")
		return
	}

	headers = []Header{&callID}

	return
}

// Parse a string representation of a Via header, returning a slice of at most one ViaHeader.
// Note that although Via headers may contain a comma-separated list, RFC 3261 makes it clear that
// these should not be treated as separate logical Via headers, but as multiple values on a single
// Via header.
func parseViaHeader(headerName string, headerText string) (
	headers []Header, err error) {
	sections := strings.Split(headerText, ",")
	var via = ViaHeader{}
	for _, section := range sections {
		var hop ViaHop
		parts := strings.Split(section, "/")

		if len(parts) < 3 {
			err = fmt.Errorf("not enough protocol parts in via header: '%s'", parts)
			return
		}

		parts[2] = strings.Join(parts[2:], "/")

		// The transport part ends when whitespace is reached, but may also start with
		// whitespace.
		// So the end of the transport part is the first whitespace char following the
		// first non-whitespace char.
		initialSpaces := len(parts[2]) - len(strings.TrimLeft(parts[2], abnfWs))
		sentByIdx := strings.IndexAny(parts[2][initialSpaces:], abnfWs) + initialSpaces + 1
		if sentByIdx == 0 {
			err = fmt.Errorf("expected whitespace after sent-protocol part "+
				"in via header '%s'", section)
			return
		} else if sentByIdx == 1 {
			err = fmt.Errorf("empty transport field in via header '%s'", section)
			return
		}

		hop.ProtocolName = strings.TrimSpace(parts[0])
		hop.ProtocolVersion = strings.TrimSpace(parts[1])
		hop.Transport = strings.TrimSpace(parts[2][:sentByIdx-1])

		if len(hop.ProtocolName) == 0 {
			err = fmt.Errorf("no protocol name provided in via header '%s'", section)
		} else if len(hop.ProtocolVersion) == 0 {
			err = fmt.Errorf("no version provided in via header '%s'", section)
		} else if len(hop.Transport) == 0 {
			err = fmt.Errorf("no transport provided in via header '%s'", section)
		}
		if err != nil {
			return
		}

		viaBody := parts[2][sentByIdx:]

		paramsIdx := strings.Index(viaBody, ";")
		var host string
		var port *Port
		if paramsIdx == -1 {
			// There are no header parameters, so the rest of the Via body is part of the host[:post].
			host, port, err = ParseHostPort(viaBody)
			hop.Host = host
			hop.Port = port
			if err != nil {
				return
			}
			hop.Params = NewParams()
		} else {
			host, port, err = ParseHostPort(viaBody[:paramsIdx])
			if err != nil {
				return
			}
			hop.Host = host
			hop.Port = port

			hop.Params, _, err = ParseParams(viaBody[paramsIdx:],
				';', ';', 0, true, true)
		}
		via = append(via, &hop)
	}

	headers = []Header{via}
	return
}

// Parse a string representation of a Max-Forwards header into a slice of at most one MaxForwards header object.
func parseMaxForwards(headerName string, headerText string) (
	headers []Header, err error) {
	var maxForwards MaxForwards
	var value uint64
	value, err = strconv.ParseUint(strings.TrimSpace(headerText), 10, 32)
	maxForwards = MaxForwards(value)

	headers = []Header{&maxForwards}
	return
}

func parseExpires(headerName string, headerText string) (headers []Header, err error) {
	var expires Expires
	var value uint64
	value, err = strconv.ParseUint(strings.TrimSpace(headerText), 10, 32)
	expires = Expires(value)
	headers = []Header{&expires}

	return
}

func parseUserAgent(headerName string, headerText string) (headers []Header, err error) {
	var userAgent UserAgentHeader
	headerText = strings.TrimSpace(headerText)
	userAgent = UserAgentHeader(headerText)
	headers = []Header{&userAgent}

	return
}

func parseContentType(headerName string, headerText string) (headers []Header, err error) {
	var contentType ContentType
	headerText = strings.TrimSpace(headerText)
	contentType = ContentType(headerText)
	headers = []Header{&contentType}

	return
}

func parseAccept(headerName string, headerText string) (headers []Header, err error) {
	var accept Accept
	headerText = strings.TrimSpace(headerText)
	accept = Accept(headerText)
	headers = []Header{&accept}

	return
}

func parseAllow(headerName string, headerText string) (headers []Header, err error) {
	allow := make(AllowHeader, 0)
	methods := strings.Split(headerText, ",")
	for _, method := range methods {
		allow = append(allow, RequestMethod(strings.TrimSpace(method)))
	}
	headers = []Header{allow}

	return
}

func parseSupported(headerName string, headerText string) (headers []Header, err error) {
	var supported SupportedHeader
	supported.Options = make([]string, 0)
	extensions := strings.Split(headerText, ",")
	for _, ext := range extensions {
		supported.Options = append(supported.Options, strings.TrimSpace(ext))
	}
	headers = []Header{&supported}

	return
}

// Parse a string representation of a Content-Length header into a slice of at most one ContentLength header object.
func parseContentLength(headerName string, headerText string) (
	headers []Header, err error) {
	var contentLength ContentLength
	var value uint64
	value, err = strconv.ParseUint(strings.TrimSpace(headerText), 10, 32)
	contentLength = ContentLength(value)

	headers = []Header{&contentLength}
	return
}

// ParseAddressValues parses a comma-separated list of addresses, returning
// any display names and header params, as well as the SIP URIs themselves.
// ParseAddressValues is aware of < > bracketing and quoting, and will not
// break on commas within these structures.
func ParseAddressValues(addresses string) (
	displayNames []MaybeString,
	uris []*URI,
	headerParams []Params,
	err error,
) {

	prevIdx := 0
	inBrackets := false
	inQuotes := false

	// Append a comma to simplify the parsing code; we split address sections
	// on commas, so use a comma to signify the end of the final address section.
	addresses = addresses + ","

	for idx, char := range addresses {
		if char == '<' && !inQuotes {
			inBrackets = true
		} else if char == '>' && !inQuotes {
			inBrackets = false
		} else if char == '"' {
			inQuotes = !inQuotes
		} else if !inQuotes && !inBrackets && char == ',' {
			var displayName MaybeString
			var uri *URI
			var params Params
			displayName, uri, params, err = ParseAddressValue(addresses[prevIdx:idx])
			if err != nil {
				return
			}
			prevIdx = idx + 1

			displayNames = append(displayNames, displayName)
			uris = append(uris, uri)
			headerParams = append(headerParams, params)
		}
	}

	return
}

// ParseAddressValue parses an address - such as from a From, To, or
// Contact header. It returns:
//   - a MaybeString containing the display name (or not)
//   - a parsed SipUri object
//   - a map containing any header parameters present
//   - the error object
//
// See RFC 3261 section 20.10 for details on parsing an address.
// Note that this method will not accept a comma-separated list of addresses;
// addresses in that form should be handled by ParseAddressValues.
func ParseAddressValue(addressText string) (
	displayName MaybeString,
	uri *URI,
	headerParams Params,
	err error,
) {

	headerParams = NewParams()

	if len(addressText) == 0 {
		err = fmt.Errorf("address-type header has empty body")
		return
	}

	addressTextCopy := addressText
	addressText = strings.TrimSpace(addressText)

	firstAngleBracket := findUnescaped(addressText, '<', quotesDelim)
	displayName = String{}
	if firstAngleBracket > 0 {
		// We have an angle bracket, and it's not the first character.
		// Since we have just trimmed whitespace, this means there must
		// be a display name.
		if addressText[0] == '"' {
			// The display name is within quotations.
			// So it is comprised of all text until the closing quote.
			addressText = addressText[1:]
			nextQuote := strings.Index(addressText, "\"")

			if nextQuote == -1 {
				// Unclosed quotes - parse error.
				err = fmt.Errorf("unclosed quotes in header text: %s",
					addressTextCopy)
				return
			}

			nameField := addressText[:nextQuote]
			displayName = String{Str: nameField}
			addressText = addressText[nextQuote+1:]
		} else {
			// The display name is unquoted, so it is comprised of
			// all text until the opening angle bracket, except surrounding whitespace.
			// According to the ABNF grammar: display-name   =  *(token LWS)/ quoted-string
			// there are certain characters the display name cannot contain unless it's quoted,
			// however we don't check for them here since it doesn't impact parsing.
			// May as well be lenient.
			nameField := addressText[:firstAngleBracket]
			displayName = String{Str: strings.TrimSpace(nameField)}
			addressText = addressText[firstAngleBracket:]
		}
	}

	// Work out where the SIP URI starts and ends.
	addressText = strings.TrimSpace(addressText)
	var endOfURI int
	var startOfParams int
	if addressText[0] != '<' {
		if displayName != nil {
			// The address must be in <angle brackets> if a display name is
			// present, so this is an invalid address line.
			err = fmt.Errorf(
				"invalid character '%c' following display "+
					"name in address line; expected '<': %s",
				addressText[0],
				addressTextCopy,
			)
			return
		}

		endOfURI = strings.Index(addressText, ";")
		if endOfURI == -1 {
			endOfURI = len(addressText)
		}
		startOfParams = endOfURI

	} else {
		addressText = addressText[1:]
		endOfURI = strings.Index(addressText, ">")
		if endOfURI == 0 {
			err = fmt.Errorf("'<' without closing '>' in address %s",
				addressTextCopy)
			return
		}
		startOfParams = endOfURI + 1

	}

	// Now parse the SIP URI.
	uri, err = ParseURI(addressText[:endOfURI])
	if err != nil {
		return
	}

	if startOfParams >= len(addressText) {
		return
	}

	// Finally, parse any header parameters and then return.
	addressText = addressText[startOfParams:]
	headerParams, _, err = ParseParams(addressText, ';', ';', ',', true, true)
	return
}

func parseRouteHeader(headerName string, headerText string) (headers []Header, err error) {
	var routeHeader RouteHeader
	routeHeader.Addresses = []*URI{}
	if _, uris, _, err := ParseAddressValues(headerText); err == nil {
		routeHeader.Addresses = uris
	} else {
		return nil, err
	}
	return []Header{&routeHeader}, nil
}

func parseRecordRouteHeader(headerName string, headerText string) (headers []Header, err error) {
	var routeHeader RecordRouteHeader
	routeHeader.Addresses = []*URI{}
	if _, uris, _, err := ParseAddressValues(headerText); err == nil {
		routeHeader.Addresses = uris
	} else {
		return nil, err
	}
	return []Header{&routeHeader}, nil
}

type parser struct {
	out    chan Message
	in     chan Packet
	isStop bool
}

func newParser() *parser {
	p := &parser{out: make(chan Message), in: make(chan Packet)}
	go p.start()
	return p
}
func (p *parser) stop() {
	p.isStop = true
}

func (p *parser) start() {
	var termErr error
	var msg Message
	var packet Packet

	for !p.isStop {
		termErr = nil
		packet = <-p.in
		startLine, err := packet.nextLine()
		if err != nil {
			logrus.Errorln(err, "parserMessage", "getStartLine", startLine)
			continue
		}
		if isRequest(startLine) {
			method, recipient, sipVersion, err := ParseRequestLine(startLine)
			if err == nil {
				msg = NewRequest("", method, recipient, sipVersion, []Header{}, []byte{})
			} else {
				termErr = utils.NewError(err, "parserMessage", "ParseRequestLine", startLine)
			}
		} else if isResponse(startLine) {
			sipVersion, statusCode, reason, err := ParseStatusLine(startLine)
			if err == nil {
				msg = NewResponse("", sipVersion, statusCode, reason, []Header{}, []byte{})
			} else {
				termErr = utils.NewError(err, "parserMessage", "ParseStatusLine", startLine)
			}
		} else {
			logrus.Errorln("undefind startlint type", startLine)
			continue
		}
		if termErr != nil {
			logrus.Errorln(err)
			continue
		}
		var buffer bytes.Buffer
		headers := make([]Header, 0)

		flushBuffer := func() {
			if buffer.Len() > 0 {
				newHeaders, err := ParseHeader(buffer.String())
				if err == nil {
					headers = append(headers, newHeaders...)
				} else {
					logrus.Warnf("skip header '%s' due to error: %s", buffer, err)
				}
				buffer.Reset()
			}
		}

		for {
			line, err := packet.nextLine()
			if err != nil {
				break
			}
			if len(line) == 0 {
				// We've hit the end of the header section.
				// Parse anything remaining in the buffer, then break out.
				flushBuffer()
				break
			}

			if !strings.Contains(abnfWs, string(line[0])) {
				// This line starts a new header.
				// Parse anything currently in the buffer, then store the new header line in the buffer.
				flushBuffer()
				buffer.WriteString(line)
			} else if buffer.Len() > 0 {
				// This is a continuation line, so just add it to the buffer.
				buffer.WriteString(" ")
				buffer.WriteString(line)
			}
		}
		// Store the headers in the message object.
		for _, header := range headers {
			msg.AppendHeader(header)
		}
		if length, ok := msg.ContentLength(); ok {
			if int(*length) > packet.bodylength {
				packet.bodylength = int(*length)
			}
		}
		body, err := packet.getBody()
		if err != nil {
			logrus.Errorln("readbody error:", err, "want size:", packet.bodylength, "body:", len(body), string(body))
			continue
		}
		if len(body) != 0 {
			msg.SetBody(body, false)
		}
		msg.SetSource(packet.raddr)
		p.out <- msg
	}
}

// func getStartLine(data []byte) (string, error) {
// 	return bufio.NewReader(bytes.NewBuffer(data)).ReadString('\r')
// }

// Calculate the size of a SIP message's body, given the entire contents of the message as a byte array.
func getBodyLength(data []byte) int {

	s := string(data)
	// Body starts with first character following a double-CRLF.
	idx := strings.Index(s, "\r\n\r\n")
	if idx == -1 {
		return -1
	}

	bodyStart := idx + 4
	return len(s) - bodyStart
}

// Heuristic to determine if the given transmission looks like a SIP request.
// It is guaranteed that any RFC3261-compliant request will pass this test,
// but invalid messages may not necessarily be rejected.
func isRequest(startLine string) bool {
	// SIP request lines contain precisely two spaces.
	if strings.Count(startLine, " ") != 2 {
		return false
	}

	// Check that the version string starts with SIP.
	parts := strings.Split(startLine, " ")
	if len(parts) < 3 {
		return false
	} else if len(parts[2]) < 3 {
		return false
	} else {
		return strings.ToUpper(parts[2][:3]) == "SIP"
	}
}

// ParseRequestLine the first line of a SIP request, e.g:
//
//	INVITE bob@example.com SIP/2.0
//	REGISTER jane@telco.com SIP/1.0
func ParseRequestLine(requestLine string) (
	method RequestMethod, recipient *URI, sipVersion string, err error) {
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		err = fmt.Errorf("request line should have 2 spaces: '%s'", requestLine)
		return
	}

	method = RequestMethod(strings.ToUpper(parts[0]))
	recipient, err = ParseURI(parts[1])
	sipVersion = parts[2]
	return
}

// ParseStatusLine the first line of a SIP response, e.g:
//
//	SIP/2.0 200 OK
//	SIP/1.0 403 Forbidden
func ParseStatusLine(statusLine string) (
	sipVersion string, statusCode int, reasonPhrase string, err error) {
	parts := strings.Split(statusLine, " ")
	if len(parts) < 3 {
		err = fmt.Errorf("status line has too few spaces: '%s'", statusLine)
		return
	}

	sipVersion = parts[0]
	statusCodeRaw, err := strconv.ParseUint(parts[1], 10, 16)
	statusCode = int(statusCodeRaw)
	reasonPhrase = strings.Join(parts[2:], " ")

	return
}

// Heuristic to determine if the given transmission looks like a SIP response.
// It is guaranteed that any RFC3261-compliant response will pass this test,
// but invalid messages may not necessarily be rejected.
func isResponse(startLine string) bool {
	// SIP status lines contain at least two spaces.
	if strings.Count(startLine, " ") < 2 {
		return false
	}

	// Check that the version string starts with SIP.
	parts := strings.Split(startLine, " ")
	if len(parts) < 3 {
		return false
	} else if len(parts[0]) < 3 {
		return false
	} else {
		return strings.ToUpper(parts[0][:3]) == "SIP"
	}
}

// ParseURI converts a string representation of a URI into a Uri object.
// If the URI is malformed, or the URI schema is not recognised, an error is returned.
// URIs have the general form of schema:address.
func ParseURI(uriStr string) (uri *URI, err error) {
	if strings.TrimSpace(uriStr) == "*" {
		// Wildcard '*' URI used in the Contact headers of REGISTERs when unregistering.
		err = fmt.Errorf("no Support *")
		return
	}

	colonIdx := strings.Index(uriStr, ":")
	if colonIdx == -1 {
		err = fmt.Errorf("no ':' in URI %s", uriStr)
		return
	}

	switch strings.ToLower(uriStr[:colonIdx]) {
	case "sip", "sips":
		var sipURI URI
		sipURI, err = ParseSipURI(uriStr)
		uri = &sipURI
	default:
		err = fmt.Errorf("unsupported URI schema %s", uriStr[:colonIdx])
	}

	return
}

// ParseSipURI converts a string representation of a SIP or SIPS URI into a SipURI object.
func ParseSipURI(uriStr string) (uri URI, err error) {
	// Store off the original URI in case we need to print it in an error.
	uriStrCopy := uriStr

	// URI should start 'sip' or 'sips'. Check the first 3 chars.
	if strings.ToLower(uriStr[:3]) != "sip" {
		err = fmt.Errorf("invalid SIP uri protocol name in '%s'", uriStrCopy)
		return
	}
	uriStr = uriStr[3:]

	if strings.ToLower(uriStr[0:1]) == "s" {
		// URI started 'sips', so it's encrypted.
		uri.FIsEncrypted = true
		uriStr = uriStr[1:]
	}

	// The 'sip' or 'sips' protocol name should be followed by a ':' character.
	if uriStr[0] != ':' {
		err = fmt.Errorf("no ':' after protocol name in SIP uri '%s'", uriStrCopy)
		return
	}
	uriStr = uriStr[1:]

	// SIP URIs may contain a user-info part, ending in a '@'.
	// This is the only place '@' may occur, so we can use it to check for the
	// existence of a user-info part.
	endOfUserInfoPart := strings.Index(uriStr, "@")
	if endOfUserInfoPart != -1 {
		// A user-info part is present. These take the form:
		//     user [ ":" password ] "@"
		endOfUsernamePart := strings.Index(uriStr, ":")
		if endOfUsernamePart > endOfUserInfoPart {
			endOfUsernamePart = -1
		}

		if endOfUsernamePart == -1 {
			// No password component; the whole of the user-info part before
			// the '@' is a username.
			uri.FUser = String{Str: uriStr[:endOfUserInfoPart]}
		} else {
			uri.FUser = String{Str: uriStr[:endOfUsernamePart]}
			uri.FPassword = String{Str: uriStr[endOfUsernamePart+1 : endOfUserInfoPart]}
		}
		uriStr = uriStr[endOfUserInfoPart+1:]
	}

	// A ';' indicates the beginning of a URI params section, and the end of the URI itself.
	endOfURIPart := strings.Index(uriStr, ";")
	if endOfURIPart == -1 {
		// There are no URI parameters, but there might be header parameters (introduced by '?').
		endOfURIPart = strings.Index(uriStr, "?")
	}
	if endOfURIPart == -1 {
		// There are no parameters at all. The URI ends after the host[:port] part.
		endOfURIPart = len(uriStr)
	}

	uri.FHost, uri.FPort, err = ParseHostPort(uriStr[:endOfURIPart])
	uriStr = uriStr[endOfURIPart:]
	if err != nil {
		return
	} else if len(uriStr) == 0 {
		uri.FUriParams = NewParams()
		uri.FHeaders = NewParams()
		return
	}

	// Now parse any URI parameters.
	// These are key-value pairs separated by ';'.
	// They end at the end of the URI, or at the start of any URI headers
	// which may be present (denoted by an initial '?').
	var uriParams Params
	var n int
	if uriStr[0] == ';' {
		uriParams, n, err = ParseParams(uriStr, ';', ';', '?', true, true)
		if err != nil {
			return
		}
	} else {
		uriParams, n = NewParams(), 0
	}
	uri.FUriParams = uriParams
	uriStr = uriStr[n:]

	// Finally parse any URI headers.
	// These are key-value pairs, starting with a '?' and separated by '&'.
	var headers Params
	headers, n, err = ParseParams(uriStr, '?', '&', 0, true, false)
	if err != nil {
		return
	}
	uri.FHeaders = headers
	uriStr = uriStr[n:]
	if len(uriStr) > 0 {
		err = fmt.Errorf("internal error: parse of SIP uri ended early! '%s'",
			uriStrCopy)
		return // Defensive return
	}

	return
}

// ParseHostPort a text representation of a host[:port] pair.
// The port may or may not be present, so we represent it with a *uint16,
// and return 'nil' if no port was present.
func ParseHostPort(rawText string) (host string, port *Port, err error) {
	colonIdx := strings.Index(rawText, ":")
	if colonIdx == -1 {
		host = rawText
		return
	}

	// Surely there must be a better way..!
	var portRaw64 uint64
	var portRaw16 uint16
	host = rawText[:colonIdx]
	portRaw64, err = strconv.ParseUint(rawText[colonIdx+1:], 10, 16)
	portRaw16 = uint16(portRaw64)
	port = (*Port)(&portRaw16)

	return
}

// ParseParams General utility method for parsing 'key=value' parameters.
// Takes a string (source), ensures that it begins with the 'start' character provided,
// and then parses successive key/value pairs separated with 'sep',
// until either 'end' is reached or there are no characters remaining.
// A map of keys to values will be returned, along with the number of characters consumed.
// Provide 0 for start or end to indicate that there is no starting/ending delimiter.
// If quoteValues is true, values can be enclosed in double-quotes which will be validated by the
// parser and omitted from the returned map.
// If permitSingletons is true, keys with no values are permitted.
// These will result in a nil value in the returned map.
func ParseParams(
	source string,
	start uint8,
	sep uint8,
	end uint8,
	quoteValues bool,
	permitSingletons bool,
) (
	params Params, consumed int, err error) {

	params = NewParams()

	if len(source) == 0 {
		// Key-value section is completely empty; return defaults.
		return
	}

	// Ensure the starting character is correct.
	if start != 0 {
		if source[0] != start {
			err = fmt.Errorf(
				"expected %c at start of key-value section; got %c. section was %s",
				start,
				source[0],
				source,
			)
			return
		}
		consumed++
	}

	// Statefully parse the given string one character at a time.
	var buffer bytes.Buffer
	var key string
	parsingKey := true // false implies we are parsing a value
	inQuotes := false
parseLoop:
	for ; consumed < len(source); consumed++ {
		switch source[consumed] {
		case end:
			if inQuotes {
				// We read an end character, but since we're inside quotations we should
				// treat it as a literal part of the value.
				buffer.WriteString(string(end))
				continue
			}

			break parseLoop

		case sep:
			if inQuotes {
				// We read a separator character, but since we're inside quotations
				// we should treat it as a literal part of the value.
				buffer.WriteString(string(sep))
				continue
			}
			if parsingKey && permitSingletons {
				params.Add(buffer.String(), nil)
			} else if parsingKey {
				err = fmt.Errorf(
					"singleton param '%s' when parsing params which disallow singletons: \"%s\"",
					buffer.String(),
					source,
				)
				return
			} else {
				params.Add(key, String{Str: buffer.String()})
			}
			buffer.Reset()
			parsingKey = true

		case '"':
			if !quoteValues {
				// We hit a quote character, but since quoting is turned off we treat it as a literal.
				buffer.WriteString("\"")
				continue
			}

			if parsingKey {
				// Quotes are never allowed in keys.
				err = fmt.Errorf("unexpected '\"' in parameter key in params \"%s\"", source)
				return
			}

			if !inQuotes && buffer.Len() != 0 {
				// We hit an initial quote midway through a value; that's not allowed.
				err = fmt.Errorf("unexpected '\"' in params \"%s\"", source)
				return
			}

			if inQuotes &&
				consumed != len(source)-1 &&
				source[consumed+1] != sep {
				// We hit an end-quote midway through a value; that's not allowed.
				err = fmt.Errorf("unexpected character %c after quoted param in \"%s\"",
					source[consumed+1], source)

				return
			}

			inQuotes = !inQuotes

		case '=':
			if buffer.Len() == 0 {
				err = fmt.Errorf("key of length 0 in params \"%s\"", source)
				return
			}
			if !parsingKey {
				err = fmt.Errorf("unexpected '=' char in value token: \"%s\"", source)
				return
			}
			key = buffer.String()
			buffer.Reset()
			parsingKey = false

		default:
			if !inQuotes && strings.Contains(abnfWs, string(source[consumed])) {
				// Skip unquoted whitespace.
				continue
			}

			buffer.WriteString(string(source[consumed]))
		}
	}

	// The param string has ended. Check that it ended in a valid place, and then store off the
	// contents of the buffer.
	if inQuotes {
		err = fmt.Errorf("unclosed quotes in parameter string: %s", source)
	} else if parsingKey && permitSingletons {
		params.Add(buffer.String(), nil)
	} else if parsingKey {
		err = fmt.Errorf("singleton param '%s' when parsing params which disallow singletons: \"%s\"",
			buffer.String(), source)
	} else {
		params.Add(key, String{Str: buffer.String()})
	}
	return
}

// ParseHeader ParseHeader
func ParseHeader(headerText string) (headers []Header, err error) {

	headers = make([]Header, 0)

	colonIdx := strings.Index(headerText, ":")
	if colonIdx == -1 {
		err = fmt.Errorf("field name with no value in header: %s", headerText)
		return
	}

	fieldName := strings.TrimSpace(headerText[:colonIdx])
	lowerFieldName := strings.ToLower(fieldName)
	fieldText := strings.TrimSpace(headerText[colonIdx+1:])
	if headerParser, ok := defaultHeaderParsers[lowerFieldName]; ok {
		// We have a registered parser for this header type - use it.
		headers, err = headerParser(lowerFieldName, fieldText)
	} else {
		// We have no registered parser for this header type,
		// so we encapsulate the header data in a GenericHeader struct.

		header := GenericHeader{
			HeaderName: fieldName,
			Contents:   fieldText,
		}
		headers = []Header{&header}
	}
	return
}

// SplitByWhitespace Splits the given string into sections, separated by one or more characters
// from c_ABNF_WS.
func SplitByWhitespace(text string) []string {
	var buffer bytes.Buffer
	var inString = true
	result := make([]string, 0)

	for _, char := range text {
		s := string(char)
		if strings.Contains(abnfWs, s) {
			if inString {
				// First whitespace char following text; flush buffer to the results array.
				result = append(result, buffer.String())
				buffer.Reset()
			}
			inString = false
		} else {
			buffer.WriteString(s)
			inString = true
		}
	}

	if buffer.Len() > 0 {
		result = append(result, buffer.String())
	}

	return result
}

// Define common quote characters needed in parsing.
var quotesDelim = delimiter{'"', '"'}

var anglesDelim = delimiter{'<', '>'}

// A delimiter is any pair of characters used for quoting text (i.e. bulk escaping literals).
type delimiter struct {
	start uint8
	end   uint8
}

func findUnescaped(text string, target uint8, delims ...delimiter) int {
	return findAnyUnescaped(text, string(target), delims...)
}

// Find the first instance of any of the targets in the given text that are not enclosed in any delimiters
// from the list provided.
func findAnyUnescaped(text string, targets string, delims ...delimiter) int {
	escaped := false
	var endEscape uint8

	endChars := make(map[uint8]uint8)
	for _, delim := range delims {
		endChars[delim.start] = delim.end
	}

	for idx := 0; idx < len(text); idx++ {
		if !escaped && strings.Contains(targets, string(text[idx])) {
			return idx
		}

		if escaped {
			escaped = text[idx] != endEscape
			continue
		} else {
			endEscape, escaped = endChars[text[idx]]
		}
	}

	return -1
}

// Uint16PtrEq Check two uint16 pointers for equality as follows:
// - If neither pointer is nil, check equality of the underlying uint16s.
// - If either pointer is nil, return true if and only if they both are.
func Uint16PtrEq(a *uint16, b *uint16) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b
}
