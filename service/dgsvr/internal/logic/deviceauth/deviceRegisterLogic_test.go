package deviceauthlogic

import (
	"encoding/json"
	"strings"
	"testing"

	"gitee.com/unitedrhino/things/share/devices"
)

func Test_getSignature(t *testing.T) {
	type args struct {
		secret string
		dest   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "legacy base64 encoded hex hmac",
			args: args{
				secret: "HaWo5qNOmisSLb/36oRUfwAY43A=",
				dest:   "deviceName=test&nonce=2428685019&productID=66&timestamp=1756780254",
			},
			want: "ZmQyMDk0YTlmNDE4YWMxN2Y5NDU2ZDgyZjM4MTdiNGJmZmRmMDU1OA==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSignature(tt.args.secret, tt.args.dest); got != tt.want {
				t.Errorf("getSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkSignatureAcceptsLegacyBase64AndHex(t *testing.T) {
	secret := "HaWo5qNOmisSLb/36oRUfwAY43A="
	dest := "deviceName=test&nonce=2428685019&productID=66&timestamp=1756780254"
	legacy := "ZmQyMDk0YTlmNDE4YWMxN2Y5NDU2ZDgyZjM4MTdiNGJmZmRmMDU1OA=="
	hexSign := "fd2094a9f418ac17f9456d82f3817b4bffdf0558"

	if !checkSignature(secret, dest, legacy) {
		t.Fatalf("legacy base64 signature should be accepted")
	}
	if !checkSignature(secret, dest, hexSign) {
		t.Fatalf("plain hex signature should be accepted")
	}
	if !checkSignature(secret, dest, strings.ToUpper(hexSign)) {
		t.Fatalf("plain hex signature should be case insensitive")
	}
	if checkSignature(secret, dest, "bad-signature") {
		t.Fatalf("invalid signature should be rejected")
	}
}

func Test_getPayloadHexReturnsPlainHexSecret(t *testing.T) {
	size, payload, err := getPayload(devices.EncTypeCert, "hex", "mzxe12OY8z/im7S3DNhHsCXdB4o=", "product-secret")
	if err != nil {
		t.Fatalf("getPayload() err=%v", err)
	}
	if size != len(payload) {
		t.Fatalf("size=%d, want len(payload)=%d", size, len(payload))
	}

	var got map[string]any
	if err := json.Unmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("payload should be plain json, got %q: %v", payload, err)
	}
	if got["encryptionType"] != float64(devices.EncTypeKey) {
		t.Fatalf("encryptionType=%v, want %d", got["encryptionType"], devices.EncTypeKey)
	}
	if got["psk"] != "9b3c5ed76398f33fe29bb4b70cd847b025dd078a" {
		t.Fatalf("psk=%v", got["psk"])
	}
	if _, ok := got["pskHex"]; ok {
		t.Fatalf("plain payload should not include pskHex field")
	}
	if _, ok := got["secretFormat"]; ok {
		t.Fatalf("plain payload should not include secretFormat field")
	}
}

func Test_getPayloadKeepsEncryptedLegacyModes(t *testing.T) {
	for _, retEnc := range []string{"", "aes128ecb"} {
		t.Run(retEnc, func(t *testing.T) {
			size, payload, err := getPayload(devices.EncTypeCert, retEnc, "mzxe12OY8z/im7S3DNhHsCXdB4o=", "1234567890123456")
			if err != nil {
				t.Fatalf("getPayload() err=%v", err)
			}
			if size == 0 || payload == "" {
				t.Fatalf("size=%d payload=%q", size, payload)
			}
			if json.Valid([]byte(payload)) {
				t.Fatalf("legacy retEnc=%q should remain encrypted, got plain json %q", retEnc, payload)
			}
		})
	}
}
