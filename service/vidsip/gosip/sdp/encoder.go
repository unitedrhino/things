package sdp

func (s Session) appendAttributes(attrs Attributes) Session {
	for _, v := range attrs {
		if v.Value == blank {
			s = s.AddFlag(v.Key)
		} else {
			s = s.AddAttribute(v.Key, v.Value)
		}
	}
	return s
}

// Append encodes message to Session and returns result.
//
// See RFC 4566 Section 5.
func (m *Message) Append(s Session) Session {
	s = s.AddVersion(m.Version)
	s = s.AddOrigin(m.Origin)
	s = s.AddSessionName(m.Name)
	if len(m.Info) > 0 {
		s = s.AddSessionInfo(m.Info)
	}
	if len(m.URI) > 0 {
		s = s.AddURI(m.URI)
	}
	if len(m.Email) > 0 {
		s = s.AddEmail(m.Email)
	}
	if len(m.Phone) > 0 {
		s = s.AddPhone(m.Phone)
	}
	if !m.Connection.Blank() {
		s = s.AddConnectionData(m.Connection)
	}
	for t, v := range m.Bandwidths {
		s = s.AddBandwidth(t, v)
	}
	// One or more time descriptions ("t=" and "r=" lines)
	for _, t := range m.Timing {
		s = s.AddTiming(t.Start, t.End)
		if len(t.Offsets) > 0 {
			s = s.AddRepeatTimesCompact(t.Repeat, t.Active, t.Offsets...)
		}
	}
	if len(m.TZAdjustments) > 0 {
		s = s.AddTimeZones(m.TZAdjustments...)
	}
	if !m.Encryption.Blank() {
		s = s.AddEncryption(m.Encryption)
	}
	s = s.appendAttributes(m.Attributes)

	for i := range m.Medias {
		s = s.AddMediaDescription(m.Medias[i].Description)
		if len(m.Medias[i].Title) > 0 {
			s = s.AddSessionInfo(m.Medias[i].Title)
		}
		if !m.Medias[i].Connection.Blank() {
			s = s.AddConnectionData(m.Medias[i].Connection)
		}
		for t, v := range m.Medias[i].Bandwidths {
			s = s.AddBandwidth(t, v)
		}
		if !m.Medias[i].Encryption.Blank() {
			s = s.AddEncryption(m.Medias[i].Encryption)
		}
		s = s.appendAttributes(m.Medias[i].Attributes)
	}
	s = s.AddSSRC(m.SSRC)
	return s
}
