package pair

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	GrantTTLSeconds = 300
	grantType       = "s01_pair_grant"
)

var (
	ErrInvalidProductSecret = errors.New("invalid_product_secret")
	ErrInvalidMAC           = errors.New("invalid_mac")
	ErrInvalidGrantToken    = errors.New("invalid_grant_token")
	ErrGrantTokenExpired    = errors.New("grant_token_expired")
	ErrGrantTokenMismatch   = errors.New("grant_token_mismatch")
	ErrPairAckInvalid       = errors.New("pair_ack_invalid")
	ErrPairAckAuthInvalid   = errors.New("pair_ack_auth_invalid")
	ErrBindEpochTooOld      = errors.New("bind_epoch_too_old")
)

type GrantInput struct {
	ProductID         string
	MAC               string
	DeviceName        string
	UserID            string
	ObservedBindEpoch int64
	Now               int64
	Nonce             []byte
	MK                []byte
	SigningKey        []byte
}

type GrantResponse struct {
	ProductID  string
	MAC        string
	DeviceName string
	GrantToken string
	Nonce      string
	AuthTag    string
	TTLSec     int64
}

type VerifyGrantInput struct {
	Token      string
	SigningKey []byte
	Now        int64
	ProductID  string
	MAC        string
	DeviceName string
	UserID     string
}

type GrantPayload struct {
	Version           int64  `json:"ver"`
	Type              string `json:"type"`
	ProductID         string `json:"product_id"`
	MAC               string `json:"mac"`
	DeviceName        string `json:"device_name"`
	UserID            string `json:"user_id"`
	ObservedBindEpoch int64  `json:"observed_bind_epoch"`
	IssuedAt          int64  `json:"issued_at"`
	ExpiresAt         int64  `json:"exp"`
	NonceHex          string `json:"nonce_hex"`
	PairKeyHex        string `json:"pair_key_hex"`
	AuthTagHex        string `json:"auth_tag_hex"`
}

type PairAck struct {
	MAC       string
	BindEpoch int64
	Msg16     []byte
	Tag       []byte
}

func DecodeProductMK(secret string) ([]byte, error) {
	secret = strings.TrimSpace(secret)
	if len(secret) != 32 || !isHex(secret) {
		return nil, ErrInvalidProductSecret
	}
	mk, err := hex.DecodeString(secret)
	if err != nil || len(mk) != aes.BlockSize {
		return nil, ErrInvalidProductSecret
	}
	return mk, nil
}

func GrantSigningKey(productID, productSecretHex string) []byte {
	sum := sha256.Sum256([]byte("ykhl:s01-route-b:grant:v1:" + productID + ":" + strings.ToLower(strings.TrimSpace(productSecretHex))))
	return sum[:]
}

func NormalizeMAC(mac string) (string, []byte, error) {
	var b strings.Builder
	for _, r := range mac {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= 'a' && r <= 'f':
			b.WriteRune(r - 'a' + 'A')
		case r >= 'A' && r <= 'F':
			b.WriteRune(r)
		case r == ':' || r == '-' || r == '.' || r == ' ' || r == '\t' || r == '\n' || r == '\r':
			continue
		default:
			return "", nil, ErrInvalidMAC
		}
	}
	normalized := b.String()
	if len(normalized) != 12 || !isHex(normalized) {
		return "", nil, ErrInvalidMAC
	}
	raw, err := hex.DecodeString(normalized)
	if err != nil || len(raw) != 6 {
		return "", nil, ErrInvalidMAC
	}
	return normalized, raw, nil
}

func BuildGrant(in GrantInput) (*GrantResponse, error) {
	mac, macRaw, err := NormalizeMAC(in.MAC)
	if err != nil {
		return nil, err
	}
	deviceName := in.DeviceName
	if deviceName == "" {
		deviceName = mac
	}
	if len(in.MK) != aes.BlockSize {
		return nil, ErrInvalidProductSecret
	}
	if len(in.SigningKey) == 0 {
		return nil, ErrInvalidProductSecret
	}
	nonce := in.Nonce
	if len(nonce) == 0 {
		nonce = make([]byte, 8)
		if _, err := rand.Read(nonce); err != nil {
			return nil, err
		}
	}
	if len(nonce) != 8 {
		return nil, fmt.Errorf("invalid_nonce")
	}
	now := in.Now
	if now == 0 {
		now = time.Now().Unix()
	}

	frame := make([]byte, 0, 16)
	frame = append(frame, 0xA1, 0x10)
	frame = append(frame, macRaw...)
	frame = append(frame, nonce...)

	dak, err := cmac128(in.MK, macRaw)
	if err != nil {
		return nil, err
	}
	authTag, err := cmac64(dak, frame)
	if err != nil {
		return nil, err
	}
	pairKey, err := cmac128(dak, append([]byte("PK"), nonce...))
	if err != nil {
		return nil, err
	}

	payload := GrantPayload{
		Version:           2,
		Type:              grantType,
		ProductID:         in.ProductID,
		MAC:               mac,
		DeviceName:        deviceName,
		UserID:            in.UserID,
		ObservedBindEpoch: in.ObservedBindEpoch,
		IssuedAt:          now,
		ExpiresAt:         now + GrantTTLSeconds,
		NonceHex:          strings.ToUpper(hex.EncodeToString(nonce)),
		PairKeyHex:        strings.ToUpper(hex.EncodeToString(pairKey)),
		AuthTagHex:        strings.ToUpper(hex.EncodeToString(authTag)),
	}
	token, err := signPayload(payload, in.SigningKey)
	if err != nil {
		return nil, err
	}
	return &GrantResponse{
		ProductID:  in.ProductID,
		MAC:        mac,
		DeviceName: deviceName,
		GrantToken: token,
		Nonce:      payload.NonceHex,
		AuthTag:    payload.AuthTagHex,
		TTLSec:     GrantTTLSeconds,
	}, nil
}

func VerifyGrant(in VerifyGrantInput) (*GrantPayload, error) {
	if len(in.SigningKey) == 0 {
		return nil, ErrInvalidGrantToken
	}
	payload, err := verifyToken(in.Token, in.SigningKey)
	if err != nil {
		return nil, err
	}
	now := in.Now
	if now == 0 {
		now = time.Now().Unix()
	}
	if payload.Version != 2 || payload.Type != grantType {
		return nil, ErrInvalidGrantToken
	}
	if now > payload.ExpiresAt {
		return nil, ErrGrantTokenExpired
	}
	mac, _, err := NormalizeMAC(in.MAC)
	if err != nil {
		return nil, err
	}
	deviceName := in.DeviceName
	if deviceName == "" {
		deviceName = mac
	}
	if payload.ProductID != in.ProductID || payload.MAC != mac || payload.DeviceName != deviceName || payload.UserID != in.UserID {
		return nil, ErrGrantTokenMismatch
	}
	return payload, nil
}

func VerifyPairAck(payloadHex, expectedMAC, pairKeyHex string, observedBindEpoch int64) (*PairAck, error) {
	payload, err := hex.DecodeString(strings.TrimSpace(payloadHex))
	if err != nil || len(payload) != 24 {
		return nil, ErrPairAckInvalid
	}
	if payload[0] != 0xA1 || payload[1] != 0x81 {
		return nil, ErrPairAckInvalid
	}
	mac, _, err := NormalizeMAC(expectedMAC)
	if err != nil {
		return nil, err
	}
	ackMAC := strings.ToUpper(hex.EncodeToString(payload[2:8]))
	if ackMAC != mac {
		return nil, ErrGrantTokenMismatch
	}
	bindEpoch := int64(payload[14]) | int64(payload[15])<<8
	if bindEpoch < observedBindEpoch {
		return nil, ErrBindEpochTooOld
	}
	pairKey, err := hex.DecodeString(strings.TrimSpace(pairKeyHex))
	if err != nil || len(pairKey) != aes.BlockSize {
		return nil, ErrInvalidGrantToken
	}
	expectedTag, err := cmac64(pairKey, payload[:16])
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(expectedTag, payload[16:24]) {
		return nil, ErrPairAckAuthInvalid
	}
	return &PairAck{
		MAC:       ackMAC,
		BindEpoch: bindEpoch,
		Msg16:     append([]byte(nil), payload[:16]...),
		Tag:       append([]byte(nil), payload[16:24]...),
	}, nil
}

func signPayload(payload GrantPayload, key []byte) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadB64 := b64url(body)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(payloadB64))
	return payloadB64 + "." + b64url(mac.Sum(nil)), nil
}

func verifyToken(token string, key []byte) (*GrantPayload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidGrantToken
	}
	sig, err := b64urlDecode(parts[1])
	if err != nil {
		return nil, ErrInvalidGrantToken
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(parts[0]))
	if !hmac.Equal(mac.Sum(nil), sig) {
		return nil, ErrInvalidGrantToken
	}
	body, err := b64urlDecode(parts[0])
	if err != nil {
		return nil, ErrInvalidGrantToken
	}
	var payload GrantPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, ErrInvalidGrantToken
	}
	return &payload, nil
}

func cmac64(key, msg []byte) ([]byte, error) {
	out, err := cmac128(key, msg)
	if err != nil {
		return nil, err
	}
	return out[:8], nil
}

func cmac128(key, msg []byte) ([]byte, error) {
	if len(key) != aes.BlockSize {
		return nil, fmt.Errorf("cmac_invalid_key")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	k1, k2 := cmacSubkeys(block)
	n := len(msg) / aes.BlockSize
	rem := len(msg) % aes.BlockSize
	if len(msg) == 0 {
		n = 1
	}
	var last []byte
	if rem == 0 && len(msg) != 0 {
		last = xorBlock(msg[(n-1)*aes.BlockSize:n*aes.BlockSize], k1)
	} else {
		start := n * aes.BlockSize
		if len(msg) == 0 {
			start = 0
		}
		padded := make([]byte, aes.BlockSize)
		copy(padded, msg[start:])
		padded[len(msg[start:])] = 0x80
		last = xorBlock(padded, k2)
		if len(msg) != 0 {
			n++
		}
	}
	x := make([]byte, aes.BlockSize)
	for i := 0; i < n-1; i++ {
		y := xorBlock(x, msg[i*aes.BlockSize:(i+1)*aes.BlockSize])
		block.Encrypt(x, y)
	}
	y := xorBlock(x, last)
	out := make([]byte, aes.BlockSize)
	block.Encrypt(out, y)
	return out, nil
}

func cmacSubkeys(block cipherBlock) ([]byte, []byte) {
	zero := make([]byte, aes.BlockSize)
	l := make([]byte, aes.BlockSize)
	block.Encrypt(l, zero)
	k1 := leftShiftOneBit(l)
	if l[0]&0x80 != 0 {
		k1[15] ^= 0x87
	}
	k2 := leftShiftOneBit(k1)
	if k1[0]&0x80 != 0 {
		k2[15] ^= 0x87
	}
	return k1, k2
}

type cipherBlock interface {
	Encrypt(dst, src []byte)
}

func leftShiftOneBit(in []byte) []byte {
	out := make([]byte, len(in))
	var carry byte
	for i := len(in) - 1; i >= 0; i-- {
		nextCarry := (in[i] & 0x80) >> 7
		out[i] = (in[i] << 1) | carry
		carry = nextCarry
	}
	return out
}

func xorBlock(a, b []byte) []byte {
	out := make([]byte, aes.BlockSize)
	for i := 0; i < aes.BlockSize; i++ {
		out[i] = a[i] ^ b[i]
	}
	return out
}

func b64url(raw []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(raw), "=")
}

func b64urlDecode(value string) ([]byte, error) {
	if rem := len(value) % 4; rem != 0 {
		value += strings.Repeat("=", 4-rem)
	}
	return base64.URLEncoding.DecodeString(value)
}

func isHex(value string) bool {
	if value == "" {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
