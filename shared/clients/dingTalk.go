package clients

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/zhaoyunxing92/dingtalk/v2"
)

type DingTalk = dingtalk.DingTalk

func NewDingTalkClient(c *conf.ThirdConf) (*DingTalk, error) {
	if c == nil {
		return nil, nil
	}
	cli, err := dingtalk.NewClient(c.AppKey, c.AppSecret)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	return cli, nil
}

//
//
//
//
//
///**
// * 使用 Token 初始化账号Client
// * @return Client
// * @throws Exception
// */
//func CreateClient() (_result *dingtalkoauth2_1_0.Client, _err error) {
//	config := &openapi.Config{}
//	config.Protocol = tea.String("https")
//	config.RegionId = tea.String("central")
//	_result = &dingtalkoauth2_1_0.Client{}
//	_result, _err = dingtalkoauth2_1_0.NewClient(config)
//	return _result, _err
//}
//
//func _main(args []*string) (_err error) {
//	client, _err := CreateClient()
//	if _err != nil {
//		return _err
//	}
//
//	getAccessTokenRequest := &dingtalkoauth2_1_0.GetAccessTokenRequest{
//		AppKey:    tea.String("dingu3zh1aifodxsgtye"),
//		AppSecret: tea.String("MLPuoBoeMdko2xVFLgkrfLYe_GU_nk0urFx1fPgV_heCgAIh0VXzikxWa43zGpIM"),
//	}
//	tryErr := func() (_e error) {
//		defer func() {
//			if r := tea.Recover(recover()); r != nil {
//				_e = r
//			}
//		}()
//		_, _err = client.GetAccessToken(getAccessTokenRequest)
//		if _err != nil {
//			return _err
//		}
//
//		return nil
//	}()
//
//	if tryErr != nil {
//		var err = &tea.SDKError{}
//		if _t, ok := tryErr.(*tea.SDKError); ok {
//			err = _t
//		} else {
//			err.Message = tea.String(tryErr.Error())
//		}
//		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {
//			// err 中含有 code 和 message 属性，可帮助开发定位问题
//		}
//
//	}
//	return _err
//}
//
//func _main2(args []*string) (_err error) {
//	client, _err := CreateClient()
//	if _err != nil {
//		return _err
//	}
//
//	getUserTokenRequest := &dingtalkoauth2_1_0.GetUserTokenRequest{
//		ClientSecret: tea.String("wgea"),
//		ClientId:     tea.String("erqFDGS"),
//		Code:         tea.String("afae"),
//		GrantType:    tea.String("gwegweg"),
//		RefreshToken: tea.String("asdfe"),
//	}
//	tryErr := func() (_e error) {
//		defer func() {
//			if r := tea.Recover(recover()); r != nil {
//				_e = r
//			}
//		}()
//		_, _err = client.GetUserToken(getUserTokenRequest)
//		if _err != nil {
//			return _err
//		}
//
//		return nil
//	}()
//
//	if tryErr != nil {
//		var err = &tea.SDKError{}
//		if _t, ok := tryErr.(*tea.SDKError); ok {
//			err = _t
//		} else {
//			err.Message = tea.String(tryErr.Error())
//		}
//		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {
//			// err 中含有 code 和 message 属性，可帮助开发定位问题
//		}
//
//	}
//	return _err
//}
//
//func main() {
//	err := _main(tea.StringSlice(os.Args[1:]))
//	if err != nil {
//		panic(err)
//	}
//}
