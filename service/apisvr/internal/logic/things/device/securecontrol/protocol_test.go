package securecontrol

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"
)

func TestDeriveAppControlKeyAESCMAC64(t *testing.T) {
	key, err := DeriveAppControlKey(GrantKeyInput{
		ProductID:    "P001",
		DeviceName:   "D001",
		DeviceSecret: "AQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHyA=",
		KeyVersion:   1,
		AuthAlg:      AuthAlgAESCMAC64,
	})
	if err != nil {
		t.Fatalf("DeriveAppControlKey returned error: %v", err)
	}
	if got, want := strings.ToUpper(hex.EncodeToString(key)), "DD3B7CE4694CE9081627106ABB3ECDE2"; got != want {
		t.Fatalf("K_app_control = %s, want %s", got, want)
	}
}

func TestDeriveAppControlKeyRejectsWeakInputs(t *testing.T) {
	for _, tc := range []GrantKeyInput{
		{ProductID: "P001", DeviceName: "D001", DeviceSecret: "AQID", KeyVersion: 1, AuthAlg: AuthAlgAESCMAC64},
		{ProductID: "P001", DeviceName: "D001", DeviceSecret: "not-base64", KeyVersion: 1, AuthAlg: AuthAlgAESCMAC64},
		{ProductID: "P001", DeviceName: "D001", DeviceSecret: "AQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHyA=", KeyVersion: 1, AuthAlg: "hmac-sha1-64"},
	} {
		if _, err := DeriveAppControlKey(tc); err == nil {
			t.Fatalf("DeriveAppControlKey(%+v) succeeded, want error", tc)
		}
	}
}

func TestBuildKeyGrantPlainResponse(t *testing.T) {
	resp, err := BuildKeyGrant(GrantInput{
		ProductID:    "P001",
		DeviceName:   "D001",
		DeviceSecret: "AQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHyA=",
		KeyVersion:   1,
		AuthAlg:      AuthAlgAESCMAC64,
	})
	if err != nil {
		t.Fatalf("BuildKeyGrant returned error: %v", err)
	}
	if resp.KeyEncoding != KeyEncodingPlainBase64 {
		t.Fatalf("KeyEncoding = %s, want %s", resp.KeyEncoding, KeyEncodingPlainBase64)
	}
	if resp.AppControlKey == "" || resp.WrappedKey != "" || resp.ServerWrapPub != "" || resp.WrapNonce != "" {
		t.Fatalf("plain grant response has wrong key fields: %+v", resp)
	}
	raw, err := base64.StdEncoding.DecodeString(resp.AppControlKey)
	if err != nil {
		t.Fatalf("AppControlKey is not base64: %v", err)
	}
	if got, want := strings.ToUpper(hex.EncodeToString(raw)), "DD3B7CE4694CE9081627106ABB3ECDE2"; got != want {
		t.Fatalf("AppControlKey = %s, want %s", got, want)
	}
	if strings.Contains(resp.AppControlKey, "AQIDBAUG") {
		t.Fatalf("response leaked deviceSecret-shaped content: %+v", resp)
	}
}

func TestBuildKeyGrantWrappedResponseCanBeOpenedByAppKey(t *testing.T) {
	curve := ecdh.X25519()
	appPriv, err := curve.NewPrivateKey(bytes.Repeat([]byte{0x11}, 32))
	if err != nil {
		t.Fatalf("app private key: %v", err)
	}
	serverPriv, err := curve.NewPrivateKey(bytes.Repeat([]byte{0x22}, 32))
	if err != nil {
		t.Fatalf("server private key: %v", err)
	}
	resp, err := BuildKeyGrant(GrantInput{
		ProductID:      "P001",
		DeviceName:     "D001",
		DeviceSecret:   "AQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHyA=",
		KeyVersion:     1,
		AuthAlg:        AuthAlgAESCMAC64,
		AppKeyWrapPub:  base64.StdEncoding.EncodeToString(appPriv.PublicKey().Bytes()),
		ServerWrapPriv: serverPriv,
		WrapNonce:      bytes.Repeat([]byte{0x33}, 12),
	})
	if err != nil {
		t.Fatalf("BuildKeyGrant returned error: %v", err)
	}
	if resp.KeyEncoding != KeyEncodingWrapped {
		t.Fatalf("KeyEncoding = %s, want %s", resp.KeyEncoding, KeyEncodingWrapped)
	}
	if resp.AppControlKey != "" || resp.WrappedKey == "" || resp.ServerWrapPub == "" || resp.WrapNonce == "" {
		t.Fatalf("wrapped grant response has wrong key fields: %+v", resp)
	}
	opened, err := OpenWrappedKeyForTest(*resp, appPriv)
	if err != nil {
		t.Fatalf("OpenWrappedKeyForTest returned error: %v", err)
	}
	if got, want := strings.ToUpper(hex.EncodeToString(opened)), "DD3B7CE4694CE9081627106ABB3ECDE2"; got != want {
		t.Fatalf("opened wrapped key = %s, want %s", got, want)
	}
}

func OpenWrappedKeyForTest(resp GrantResponse, appPriv *ecdh.PrivateKey) ([]byte, error) {
	serverPubBytes, err := base64.StdEncoding.DecodeString(resp.ServerWrapPub)
	if err != nil {
		return nil, err
	}
	serverPub, err := ecdh.X25519().NewPublicKey(serverPubBytes)
	if err != nil {
		return nil, err
	}
	shared, err := appPriv.ECDH(serverPub)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(resp.WrapNonce)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(resp.WrappedKey)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(deriveWrapKey(shared, grantAAD(&resp)))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ciphertext, grantAAD(&resp))
}
