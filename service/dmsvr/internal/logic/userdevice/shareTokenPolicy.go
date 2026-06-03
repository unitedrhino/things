package userdevicelogic

const deviceShareUseByWechatSingleDevice = "wechat_single_device"

func shouldConsumeShareTokenAfterAccept(useBy string) bool {
	return useBy == deviceShareUseByWechatSingleDevice
}
