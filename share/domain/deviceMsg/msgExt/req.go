package msgExt

import (
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
)

type (
	Req struct {
		deviceMsg.CommonMsg
	}

	RegisterReq struct {
		deviceMsg.CommonMsg
		Payload Payload `json:"payload"`
	}

	Payload struct {
		Nonce     int64  `json:"nonce"`     //随机数
		Timestamp int64  `json:"timestamp"` //秒级时间戳
		Signature string `json:"signature"` //签名信息
	}

	RespData struct {
		Len int64 `json:"len"` //payload加密前信息的长度
		/*
		  加密过程将原始 JSON 格式的 payload 转为字符串后进行 AES 加密，再进行 base64 加密。AES 加密算法为 CBC 模式，密钥长度128，取 productSecret 前16位，偏移量为长度16的字符“0”。
		  原始 payload 内容说明：
		  key                value               描述
		  encryptionType     1              加密类型，1表示证书认证，2表示签名认证。
		  psk                1239           设备密钥，当产品认证类型为签名认证时有此参数
		  clientCert         -              设备证书文件字符串格式，当产品认证类型为证书认证时有此参数。
		  clientKey          -              设备私钥文件字符串格式，当产品认证类型为证书认证时有此参数。
		*/
		Payload string `json:"payload"`
	}
)
