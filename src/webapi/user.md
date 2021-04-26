
### 1. N/A

1. 路由定义

- Url: /user/captcha
- Method: POST
- Request: `GetCaptchaReq`
- Response: `GetCaptchaResp`

2. 请求定义


```golang
type GetCaptchaReq struct {
	Data string `json:"data,optional"` //短信验证时填写手机号
	Type string `json:"type,options=sms|image"` //验证方式:短信验证,图片验证码
	Use string `json:"use,options=login|register"` //用途
}
```


3. 返回定义


```golang
type GetCaptchaResp struct {
	CodeID string `json:"codeID"` //验证码编号
	Url string `json:"url,optional"` //图片验证码url
	Expire int64 `json:"expire"` //过期时间
}
```
  


### 2. N/A

1. 路由定义

- Url: /user/login
- Method: POST
- Request: `LoginReq`
- Response: `LoginResp`

2. 请求定义


```golang
type LoginReq struct {
	UserID string `json:"userID,optional"` //登录账号(支持用户名,手机号登录) 账号密码登录时需要填写
	PwdType int32 `json:"pwdtype,optional"` //账号密码登录时需要填写.0,无密码 1，明文 2，md5加密
	Password string `json:"password,optional"` //密码，建议md5转换 密码登录时需要填写
	LoginType string `json:"loginType,options=sms|img|wxopen|wxin|wxminip"` //验证类型 sms 短信验证码 img 图形验证码加账号密码登录 wxopen 微信开放平台登录 wxin 微信内登录 wxminip 微信小程序
	Code string `json:"code,optional"` //验证码    微信登录填code
	CodeID string `json:"codeID,optional"` //验证码编号 微信登录填state
}
```


3. 返回定义


```golang
type LoginResp struct {
	Info UserInfo `json:"info"` //用户信息
	Token JwtToken `json:"token"` //用户token
}

type UserInfo struct {
	Uid int64 `json:"uid,string"` // 用户id
	UserName string `json:"userName,optional,omitempty"` //用户名(唯一)
	NickName string `json:"nickName,optional,omitempty"` // 用户的昵称
	InviterUid int64 `json:"inviterUid,string,optional,omitempty"` // 邀请人用户id
	InviterId string `json:"inviterId,optional,omitempty"` // 邀请码
	Sex int64 `json:"sex,optional,omitempty"` // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City string `json:"city,optional,omitempty"` // 用户所在城市
	Country string `json:"country,optional,omitempty"` // 用户所在国家
	Province string `json:"province,optional,omitempty"` // 用户所在省份
	Language string `json:"language,optional,omitempty"` // 用户的语言，简体中文为zh_CN
	HeadImgUrl string `json:"headImgUrl,optional,omitempty"` // 用户头像
	CreateTime int64 `json:"createTime,string,optional,omitempty"`
}

type JwtToken struct {
	AccessToken string `json:"accessToken,omitempty"`
	AccessExpire int64 `json:"accessExpire,omitempty"`
	RefreshAfter int64 `json:"refreshAfter,omitempty"`
}
```
  


### 3. N/A

1. 路由定义

- Url: /user/registerCore
- Method: POST
- Request: `RegisterCoreReq`
- Response: `RegisterCoreResp`

2. 请求定义


```golang
type RegisterCoreReq struct {
	ReqType string `json:"reqType,options=phone|wxopen|wxin|wxminip"` //注册方式:	phone手机号注册 wxopen 微信开放平台登录 wxin 微信内登录 wxminip 微信小程序
	Note string `json:"note,optional"` //手机号注册时填写手机号
	Code string `json:"code"` //验证码    微信登录填code
	CodeID string `json:"codeID,optional"` //验证码编号 微信登录填state
}
```


3. 返回定义


```golang
type RegisterCoreResp struct {
	Uid int64 `json:"uid,string"` //用户id
JwtToken
}

type JwtToken struct {
	AccessToken string `json:"accessToken,omitempty"`
	AccessExpire int64 `json:"accessExpire,omitempty"`
	RefreshAfter int64 `json:"refreshAfter,omitempty"`
}
```
  


### 4. N/A

1. 路由定义

- Url: /user/register2
- Method: POST
- Request: `Register2Req`
- Response: `-`

2. 请求定义


```golang
type Register2Req struct {
	Token string `json:"token"` //注册第一步的token
	Password string `json:"password,optional"` //明文密码
UserInfo
}

type UserInfo struct {
	Uid int64 `json:"uid,string"` // 用户id
	UserName string `json:"userName,optional,omitempty"` //用户名(唯一)
	NickName string `json:"nickName,optional,omitempty"` // 用户的昵称
	InviterUid int64 `json:"inviterUid,string,optional,omitempty"` // 邀请人用户id
	InviterId string `json:"inviterId,optional,omitempty"` // 邀请码
	Sex int64 `json:"sex,optional,omitempty"` // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City string `json:"city,optional,omitempty"` // 用户所在城市
	Country string `json:"country,optional,omitempty"` // 用户所在国家
	Province string `json:"province,optional,omitempty"` // 用户所在省份
	Language string `json:"language,optional,omitempty"` // 用户的语言，简体中文为zh_CN
	HeadImgUrl string `json:"headImgUrl,optional,omitempty"` // 用户头像
	CreateTime int64 `json:"createTime,string,optional,omitempty"`
}
```


3. 返回定义
  


### 5. N/A

1. 路由定义

- Url: /user/info
- Method: GET
- Request: `-`
- Response: `UserInfo`

2. 请求定义


3. 返回定义


```golang
type UserInfo struct {
	Uid int64 `json:"uid,string"` // 用户id
	UserName string `json:"userName,optional,omitempty"` //用户名(唯一)
	NickName string `json:"nickName,optional,omitempty"` // 用户的昵称
	InviterUid int64 `json:"inviterUid,string,optional,omitempty"` // 邀请人用户id
	InviterId string `json:"inviterId,optional,omitempty"` // 邀请码
	Sex int64 `json:"sex,optional,omitempty"` // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City string `json:"city,optional,omitempty"` // 用户所在城市
	Country string `json:"country,optional,omitempty"` // 用户所在国家
	Province string `json:"province,optional,omitempty"` // 用户所在省份
	Language string `json:"language,optional,omitempty"` // 用户的语言，简体中文为zh_CN
	HeadImgUrl string `json:"headImgUrl,optional,omitempty"` // 用户头像
	CreateTime int64 `json:"createTime,string,optional,omitempty"`
}
```
  


### 6. N/A

1. 路由定义

- Url: /user/modifyUserInfo
- Method: POST
- Request: `ModifyUserInfoReq`
- Response: `-`

2. 请求定义


```golang
type ModifyUserInfoReq struct {
	Info map[string]string `json:"info"` //修改参数key value数组
}
```


3. 返回定义
  

