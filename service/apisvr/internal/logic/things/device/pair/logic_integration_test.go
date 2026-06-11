package pair

import (
	"context"
	"testing"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/productmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"google.golang.org/grpc"
)

func TestGrantConfirmLogicBindsDevice(t *testing.T) {
	const productSecret = "00112233445566778899aabbccddeeff"
	productM := &fakeProductManage{secret: productSecret}
	deviceM := &fakeDeviceManage{}
	svcCtx := &svc.ServiceContext{}
	svcCtx.ProductM = productM
	svcCtx.DeviceM = deviceM
	ctx := ctxs.SetUserCtx(context.Background(), &ctxs.UserCtx{
		UserID:     42,
		TenantCode: "default",
	})

	grantResp, err := NewGrantLogic(ctx, svcCtx).Grant(&types.DevicePairGrantReq{
		ProductID:         "S01",
		Mac:               "aa:bb:cc:dd:ee:ff",
		ObservedBindEpoch: 7,
	})
	if err != nil {
		t.Fatalf("Grant returned error: %v", err)
	}
	if grantResp.DeviceName != "AABBCCDDEEFF" {
		t.Fatalf("Grant deviceName = %s, want normalized MAC", grantResp.DeviceName)
	}

	grantPayload, err := VerifyGrant(VerifyGrantInput{
		Token:      grantResp.GrantToken,
		SigningKey: GrantSigningKey("S01", productSecret),
		ProductID:  "S01",
		MAC:        "AABBCCDDEEFF",
		DeviceName: "AABBCCDDEEFF",
		UserID:     "42",
	})
	if err != nil {
		t.Fatalf("VerifyGrant returned error: %v", err)
	}
	ackHex := BuildTestPairAckHex(t, "AABBCCDDEEFF", 9, grantPayload.PairKeyHex)

	confirmResp, err := NewConfirmLogic(ctx, svcCtx).Confirm(&types.DevicePairConfirmReq{
		ProductID:      "S01",
		Mac:            "AABBCCDDEEFF",
		GrantToken:     grantResp.GrantToken,
		PairAckPayload: ackHex,
		AreaID:         88,
	})
	if err != nil {
		t.Fatalf("Confirm returned error: %v", err)
	}
	if confirmResp.BindEpoch != 9 || confirmResp.BleSecVer != 2 || confirmResp.BlePairKey != grantPayload.PairKeyHex || confirmResp.Message != "bind_confirmed" {
		t.Fatalf("Confirm response unexpected: %#v", confirmResp)
	}
	if productM.readCalls != 2 {
		t.Fatalf("ProductInfoRead calls = %d, want 2", productM.readCalls)
	}
	if deviceM.bindCalls != 1 {
		t.Fatalf("DeviceInfoBind calls = %d, want 1", deviceM.bindCalls)
	}
	if deviceM.bindReq.AreaID != 88 || !deviceM.bindReq.IsIgnoreOffline {
		t.Fatalf("Bind options unexpected: %#v", deviceM.bindReq)
	}
	if deviceM.bindReq.Device.ProductID != "S01" || deviceM.bindReq.Device.DeviceName != "AABBCCDDEEFF" {
		t.Fatalf("Bind device unexpected: %#v", deviceM.bindReq.Device)
	}
}

type fakeProductManage struct {
	productmanage.ProductManage
	secret    string
	readCalls int
}

func (f *fakeProductManage) ProductInfoRead(ctx context.Context, in *dm.ProductInfoReadReq, opts ...grpc.CallOption) (*dm.ProductInfo, error) {
	f.readCalls++
	if in.ProductID != "S01" {
		return nil, ErrInvalidGrantToken
	}
	return &dm.ProductInfo{
		ProductID: "S01",
		Secret:    f.secret,
	}, nil
}

type fakeDeviceManage struct {
	devicemanage.DeviceManage
	bindReq   *dm.DeviceInfoBindReq
	bindCalls int
}

func (f *fakeDeviceManage) DeviceInfoBind(ctx context.Context, in *dm.DeviceInfoBindReq, opts ...grpc.CallOption) (*dm.Empty, error) {
	f.bindCalls++
	f.bindReq = in
	return &dm.Empty{}, nil
}
