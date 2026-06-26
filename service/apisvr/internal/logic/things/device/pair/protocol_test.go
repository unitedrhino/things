package pair

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestDecodeProductMKUsesFirst16ASCIIBytes(t *testing.T) {
	const productSecret = "Fqx8joXhfN7aLwlWgry2FaykK7g="
	mk, err := DecodeProductMK(productSecret)
	if err != nil {
		t.Fatalf("DecodeProductMK returned error: %v", err)
	}
	if got, want := strings.ToUpper(hex.EncodeToString(mk)), "467178386A6F5868664E37614C776C57"; got != want {
		t.Fatalf("DecodeProductMK key hex = %s, want %s", got, want)
	}
}

func TestDecodeProductMKDoesNotHexDecodeProductSecret(t *testing.T) {
	const productSecret = "00112233445566778899aabbccddeeff"
	mk, err := DecodeProductMK(productSecret)
	if err != nil {
		t.Fatalf("DecodeProductMK returned error: %v", err)
	}
	if got, want := string(mk), "0011223344556677"; got != want {
		t.Fatalf("DecodeProductMK key text = %q, want %q", got, want)
	}

	for _, secret := range []string{
		"",
		"Fqx8joXhfN7aLwl",
	} {
		t.Run(secret, func(t *testing.T) {
			if _, err := DecodeProductMK(secret); err == nil {
				t.Fatalf("DecodeProductMK(%q) succeeded, want error", secret)
			}
		})
	}
}

func TestNormalizeMAC(t *testing.T) {
	mac, macBytes, err := NormalizeMAC("aa:bb-cc dd.ee.ff")
	if err != nil {
		t.Fatalf("NormalizeMAC returned error: %v", err)
	}
	if mac != "AABBCCDDEEFF" {
		t.Fatalf("NormalizeMAC mac = %s, want AABBCCDDEEFF", mac)
	}
	if got := strings.ToUpper(hex.EncodeToString(macBytes)); got != "AABBCCDDEEFF" {
		t.Fatalf("NormalizeMAC bytes = %s", got)
	}

	if _, _, err := NormalizeMAC("AABBCCDDEE"); err == nil {
		t.Fatalf("NormalizeMAC accepted short mac")
	}
	if _, _, err := NormalizeMAC("AABBCCDDEEGG"); err == nil {
		t.Fatalf("NormalizeMAC accepted non-hex mac")
	}
	if _, _, err := NormalizeMAC("AABBCCDDEEFFGG"); err == nil {
		t.Fatalf("NormalizeMAC accepted trailing non-hex bytes")
	}
}

func TestBuildAndVerifyGrant(t *testing.T) {
	mk, err := DecodeProductMK("00112233445566778899aabbccddeeff")
	if err != nil {
		t.Fatal(err)
	}
	signingKey := GrantSigningKey("S01", "00112233445566778899aabbccddeeff")
	grant, err := BuildGrant(GrantInput{
		ProductID:         "S01",
		MAC:               "AABBCCDDEEFF",
		DeviceName:        "AABBCCDDEEFF",
		UserID:            "42",
		ObservedBindEpoch: 7,
		Now:               1700000000,
		Nonce:             mustHex(t, "0102030405060708"),
		MK:                mk,
		SigningKey:        signingKey,
	})
	if err != nil {
		t.Fatalf("BuildGrant returned error: %v", err)
	}
	if grant.Nonce != "0102030405060708" {
		t.Fatalf("Nonce = %s", grant.Nonce)
	}
	if grant.TTLSec != 300 {
		t.Fatalf("TTLSec = %d, want 300", grant.TTLSec)
	}
	if grant.AuthTag == "" || grant.GrantToken == "" {
		t.Fatalf("BuildGrant returned empty auth tag or token: %#v", grant)
	}

	payload, err := VerifyGrant(VerifyGrantInput{
		Token:      grant.GrantToken,
		SigningKey: signingKey,
		Now:        1700000100,
		ProductID:  "S01",
		MAC:        "AABBCCDDEEFF",
		DeviceName: "AABBCCDDEEFF",
		UserID:     "42",
	})
	if err != nil {
		t.Fatalf("VerifyGrant returned error: %v", err)
	}
	if payload.PairKeyHex != "" {
		t.Fatalf("grant token should not carry pair_key_hex")
	}
	if body := decodedGrantTokenBody(t, grant.GrantToken); strings.Contains(string(body), "pair_key_hex") {
		t.Fatalf("grant token payload leaked pair_key_hex: %s", string(body))
	}

	if _, err := VerifyGrant(VerifyGrantInput{
		Token:      grant.GrantToken + "x",
		SigningKey: signingKey,
		Now:        1700000100,
		ProductID:  "S01",
		MAC:        "AABBCCDDEEFF",
		DeviceName: "AABBCCDDEEFF",
		UserID:     "42",
	}); err == nil {
		t.Fatalf("VerifyGrant accepted tampered token")
	}

	if _, err := VerifyGrant(VerifyGrantInput{
		Token:      grant.GrantToken,
		SigningKey: signingKey,
		Now:        1700000301,
		ProductID:  "S01",
		MAC:        "AABBCCDDEEFF",
		DeviceName: "AABBCCDDEEFF",
		UserID:     "42",
	}); err == nil {
		t.Fatalf("VerifyGrant accepted expired token")
	}

	if _, err := VerifyGrant(VerifyGrantInput{
		Token:      grant.GrantToken,
		SigningKey: signingKey,
		Now:        1700000100,
		ProductID:  "S01",
		MAC:        "AABBCCDDEEFF",
		DeviceName: "other",
		UserID:     "42",
	}); err == nil {
		t.Fatalf("VerifyGrant accepted mismatched device name")
	}
}

func TestVerifyPairAck(t *testing.T) {
	mk, err := DecodeProductMK("00112233445566778899aabbccddeeff")
	if err != nil {
		t.Fatal(err)
	}
	signingKey := GrantSigningKey("S01", "00112233445566778899aabbccddeeff")
	grant, err := BuildGrant(GrantInput{
		ProductID:  "S01",
		MAC:        "AABBCCDDEEFF",
		UserID:     "42",
		Now:        1700000000,
		Nonce:      mustHex(t, "0102030405060708"),
		MK:         mk,
		SigningKey: signingKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	payload, err := VerifyGrant(VerifyGrantInput{
		Token:      grant.GrantToken,
		SigningKey: signingKey,
		Now:        1700000100,
		ProductID:  "S01",
		MAC:        "AABBCCDDEEFF",
		DeviceName: "AABBCCDDEEFF",
		UserID:     "42",
	})
	if err != nil {
		t.Fatal(err)
	}
	pairKeyHex := derivePairKeyHexForTest(t, mk, "AABBCCDDEEFF", payload.NonceHex)
	ackHex := BuildTestPairAckHex(t, "AABBCCDDEEFF", 9, pairKeyHex)
	ack, err := VerifyPairAck(ackHex, "AABBCCDDEEFF", pairKeyHex, 7)
	if err != nil {
		t.Fatalf("VerifyPairAck returned error: %v", err)
	}
	if ack.BindEpoch != 9 {
		t.Fatalf("BindEpoch = %d, want 9", ack.BindEpoch)
	}

	if _, err := VerifyPairAck(ackHex[:len(ackHex)-2], "AABBCCDDEEFF", pairKeyHex, 7); err == nil {
		t.Fatalf("VerifyPairAck accepted short payload")
	}
	if _, err := VerifyPairAck("FFFF"+ackHex[4:], "AABBCCDDEEFF", pairKeyHex, 7); err == nil {
		t.Fatalf("VerifyPairAck accepted invalid header")
	}
	if _, err := VerifyPairAck(ackHex, "001122334455", pairKeyHex, 7); err == nil {
		t.Fatalf("VerifyPairAck accepted mismatched mac")
	}
	if _, err := VerifyPairAck(ackHex, "AABBCCDDEEFF", pairKeyHex, 10); err == nil {
		t.Fatalf("VerifyPairAck accepted stale bind epoch")
	}
}

func mustHex(t *testing.T, value string) []byte {
	t.Helper()
	out, err := hex.DecodeString(value)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func BuildTestPairAckHex(t *testing.T, mac string, bindEpoch int64, pairKeyHex string) string {
	t.Helper()
	_, macBytes, err := NormalizeMAC(mac)
	if err != nil {
		t.Fatal(err)
	}
	payload := make([]byte, 16)
	payload[0] = 0xA1
	payload[1] = 0x81
	copy(payload[2:8], macBytes)
	payload[14] = byte(bindEpoch)
	payload[15] = byte(bindEpoch >> 8)
	pairKey, err := hex.DecodeString(pairKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	tag, err := cmac64(pairKey, payload)
	if err != nil {
		t.Fatal(err)
	}
	return strings.ToUpper(hex.EncodeToString(append(payload, tag...)))
}

func decodedGrantTokenBody(t *testing.T, token string) []byte {
	t.Helper()
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		t.Fatalf("grant token parts = %d, want 2", len(parts))
	}
	body, err := b64urlDecode(parts[0])
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func derivePairKeyHexForTest(t *testing.T, mk []byte, mac string, nonceHex string) string {
	t.Helper()
	_, macRaw, err := NormalizeMAC(mac)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := hex.DecodeString(nonceHex)
	if err != nil {
		t.Fatal(err)
	}
	dak, err := cmac128(mk, macRaw)
	if err != nil {
		t.Fatal(err)
	}
	pairKey, err := cmac128(dak, append([]byte("PK"), nonce...))
	if err != nil {
		t.Fatal(err)
	}
	return strings.ToUpper(hex.EncodeToString(pairKey))
}
