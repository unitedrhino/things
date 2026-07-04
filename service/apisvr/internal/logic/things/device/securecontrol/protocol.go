// Package securecontrol 提供 App 安全控制 key-grant 的密钥派生与包装能力。
package securecontrol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	// VersionAppSecControlV1 是 App 安全控制 key-grant 的协议版本。
	VersionAppSecControlV1 = "app-sec-ctrl-v1"
	// AuthAlgAESCMAC64 是低端 MCU 默认支持的 AES-CMAC-64 认证算法。
	AuthAlgAESCMAC64 = "aes-cmac-64"
	// KeyEncodingPlainBase64 表示响应体直接返回 Base64 编码的 App 控制 key。
	KeyEncodingPlainBase64 = "plain-base64-v1"
	// KeyEncodingWrapped 表示响应体返回由 App 一次性公钥包装过的 App 控制 key。
	KeyEncodingWrapped = "wrapped-key-v1"
	// DefaultScopeBits 是第一版 App 手动属性控制与行为调用的默认授权范围。
	DefaultScopeBits = "0x00000003"
	// GrantTTLSeconds 是 App 本地保存 key-grant 响应的建议有效期。
	GrantTTLSeconds = 86400
)

var (
	// ErrInvalidDeviceSecret 表示设备密钥不是合法 Base64 或不足 16 字节。
	ErrInvalidDeviceSecret = errors.New("invalid_device_secret")
	// ErrUnsupportedAuthAlg 表示设备请求了 key-grant 不支持的认证算法。
	ErrUnsupportedAuthAlg = errors.New("unsupported_auth_alg")
	// ErrInvalidKeyVersion 表示 keyVersion 超出 v1 可派生范围。
	ErrInvalidKeyVersion = errors.New("invalid_key_version")
	// ErrInvalidWrapPublicKey 表示 App 的一次性包装公钥格式不正确。
	ErrInvalidWrapPublicKey = errors.New("invalid_app_key_wrap_pub")
	// ErrInvalidWrapNonce 表示包装 nonce 格式不正确。
	ErrInvalidWrapNonce = errors.New("invalid_wrap_nonce")
)

// GrantKeyInput 描述从 deviceSecret 派生 App 控制 key 所需的稳定上下文。
type GrantKeyInput struct {
	ProductID    string
	DeviceName   string
	DeviceSecret string
	KeyVersion   int64
	AuthAlg      string
}

// GrantInput 描述生成 key-grant 响应需要的设备密钥、算法和可选包装参数。
type GrantInput struct {
	ProductID      string
	DeviceName     string
	DeviceSecret   string
	KeyVersion     int64
	AuthAlg        string
	AppKeyWrapPub  string
	ServerWrapPriv *ecdh.PrivateKey
	WrapNonce      []byte
}

// GrantResponse 是 key-grant 逻辑层返回给 handler 的响应模型。
type GrantResponse struct {
	Version       string `json:"version"`
	ProductID     string `json:"productID"`
	DeviceName    string `json:"deviceName"`
	AuthAlg       string `json:"authAlg"`
	KeyVersion    int64  `json:"keyVersion"`
	KeyEncoding   string `json:"keyEncoding"`
	AppControlKey string `json:"appControlKey,omitempty"`
	WrappedKey    string `json:"wrappedKey,omitempty"`
	WrapNonce     string `json:"wrapNonce,omitempty"`
	ServerWrapPub string `json:"serverWrapPub,omitempty"`
	ScopeBits     string `json:"scopeBits"`
	TTLSec        int64  `json:"ttlSec"`
}

// DeriveAppControlKey 使用 deviceSecret 和设备身份上下文派生 K_app_control。
func DeriveAppControlKey(in GrantKeyInput) ([]byte, error) {
	authAlg := normalizeAuthAlg(in.AuthAlg)
	if authAlg != AuthAlgAESCMAC64 {
		return nil, ErrUnsupportedAuthAlg
	}
	keyVersion, err := normalizeKeyVersion(in.KeyVersion)
	if err != nil {
		return nil, err
	}
	secret, err := base64.StdEncoding.DecodeString(in.DeviceSecret)
	if err != nil || len(secret) < aes.BlockSize {
		return nil, ErrInvalidDeviceSecret
	}
	context := appControlContext(in.ProductID, in.DeviceName, keyVersion)
	return cmac128(secret[:aes.BlockSize], context)
}

// BuildKeyGrant 根据设备密钥派生 App 控制 key，并按请求选择明文 Base64 或 wrappedKey 返回。
func BuildKeyGrant(in GrantInput) (*GrantResponse, error) {
	authAlg := normalizeAuthAlg(in.AuthAlg)
	keyVersion, err := normalizeKeyVersion(in.KeyVersion)
	if err != nil {
		return nil, err
	}
	appControlKey, err := DeriveAppControlKey(GrantKeyInput{
		ProductID:    in.ProductID,
		DeviceName:   in.DeviceName,
		DeviceSecret: in.DeviceSecret,
		KeyVersion:   keyVersion,
		AuthAlg:      authAlg,
	})
	if err != nil {
		return nil, err
	}
	resp := &GrantResponse{
		Version:     VersionAppSecControlV1,
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		AuthAlg:     authAlg,
		KeyVersion:  keyVersion,
		ScopeBits:   DefaultScopeBits,
		TTLSec:      GrantTTLSeconds,
		KeyEncoding: KeyEncodingPlainBase64,
	}
	if in.AppKeyWrapPub == "" {
		resp.AppControlKey = base64.StdEncoding.EncodeToString(appControlKey)
		return resp, nil
	}
	if err := wrapAppControlKey(resp, appControlKey, in); err != nil {
		return nil, err
	}
	return resp, nil
}

// normalizeAuthAlg 统一 key-grant 默认算法。
func normalizeAuthAlg(authAlg string) string {
	if authAlg == "" {
		return AuthAlgAESCMAC64
	}
	return authAlg
}

// normalizeKeyVersion 统一 keyVersion 默认值并限制到 uint32 派生空间。
func normalizeKeyVersion(keyVersion int64) (int64, error) {
	if keyVersion == 0 {
		return 1, nil
	}
	if keyVersion < 0 || keyVersion > int64(^uint32(0)) {
		return 0, ErrInvalidKeyVersion
	}
	return keyVersion, nil
}

// appControlContext 生成 App 控制 key 派生的稳定上下文字节。
func appControlContext(productID, deviceName string, keyVersion int64) []byte {
	out := make([]byte, 0, len("YK_APP_SEC_CTRL_V1")+len(productID)+len(deviceName)+7)
	out = append(out, []byte("YK_APP_SEC_CTRL_V1")...)
	out = append(out, 0x00)
	out = append(out, []byte(productID)...)
	out = append(out, 0x00)
	out = append(out, []byte(deviceName)...)
	out = append(out, 0x00)
	var version [4]byte
	binary.LittleEndian.PutUint32(version[:], uint32(keyVersion))
	out = append(out, version[:]...)
	return out
}

// wrapAppControlKey 使用 App 的一次性 X25519 公钥包装 App 控制 key。
func wrapAppControlKey(resp *GrantResponse, appControlKey []byte, in GrantInput) error {
	appPubBytes, err := base64.StdEncoding.DecodeString(in.AppKeyWrapPub)
	if err != nil {
		return ErrInvalidWrapPublicKey
	}
	curve := ecdh.X25519()
	appPub, err := curve.NewPublicKey(appPubBytes)
	if err != nil {
		return ErrInvalidWrapPublicKey
	}
	serverPriv := in.ServerWrapPriv
	if serverPriv == nil {
		serverPriv, err = curve.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
	}
	nonce := in.WrapNonce
	if len(nonce) == 0 {
		nonce = make([]byte, 12)
		if _, err := rand.Read(nonce); err != nil {
			return err
		}
	}
	if len(nonce) != 12 {
		return ErrInvalidWrapNonce
	}
	shared, err := serverPriv.ECDH(appPub)
	if err != nil {
		return ErrInvalidWrapPublicKey
	}
	resp.KeyEncoding = KeyEncodingWrapped
	resp.AppControlKey = ""
	resp.ServerWrapPub = base64.StdEncoding.EncodeToString(serverPriv.PublicKey().Bytes())
	resp.WrapNonce = base64.StdEncoding.EncodeToString(nonce)
	wrapKey := deriveWrapKey(shared, grantAAD(resp))
	block, err := aes.NewCipher(wrapKey)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	resp.WrappedKey = base64.StdEncoding.EncodeToString(gcm.Seal(nil, nonce, appControlKey, grantAAD(resp)))
	return nil
}

// grantAAD 生成 wrappedKey 绑定的响应元数据。
func grantAAD(resp *GrantResponse) []byte {
	return []byte(fmt.Sprintf("%s\n%s\n%s\n%s\n%d\n%s\n%s",
		resp.Version,
		resp.ProductID,
		resp.DeviceName,
		resp.AuthAlg,
		resp.KeyVersion,
		resp.KeyEncoding,
		resp.ScopeBits,
	))
}

// deriveWrapKey 用 HKDF-SHA256 从 X25519 共享密钥派生 AES-GCM 包装 key。
func deriveWrapKey(sharedSecret, info []byte) []byte {
	salt := []byte("YK_APP_SEC_CTRL_WRAP_V1")
	prkMac := hmac.New(sha256.New, salt)
	prkMac.Write(sharedSecret)
	prk := prkMac.Sum(nil)

	expander := hmac.New(sha256.New, prk)
	expander.Write(info)
	expander.Write([]byte{0x01})
	return expander.Sum(nil)
}

// cmac128 计算 AES-CMAC 完整 16 字节认证值。
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

// cmacSubkeys 生成 RFC 4493 定义的 CMAC 子密钥。
func cmacSubkeys(block cipher.Block) ([]byte, []byte) {
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

// leftShiftOneBit 对 16 字节块执行按位左移一位。
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

// xorBlock 对两个 16 字节块执行异或。
func xorBlock(a, b []byte) []byte {
	out := make([]byte, aes.BlockSize)
	for i := 0; i < aes.BlockSize; i++ {
		out[i] = a[i] ^ b[i]
	}
	return out
}
