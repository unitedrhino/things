package info

import (
	"context"
	"testing"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicegroup"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"google.golang.org/grpc"
)

type stubDeviceGroup struct {
	devicegroup.DeviceGroup
	readReq  *dm.GroupInfoReadReq
	indexReq *dm.GroupInfoIndexReq
}

func (s *stubDeviceGroup) GroupInfoRead(ctx context.Context, in *dm.GroupInfoReadReq, opts ...grpc.CallOption) (*dm.GroupInfo, error) {
	s.readReq = in
	return &dm.GroupInfo{}, nil
}

func (s *stubDeviceGroup) GroupInfoIndex(ctx context.Context, in *dm.GroupInfoIndexReq, opts ...grpc.CallOption) (*dm.GroupInfoIndexResp, error) {
	s.indexReq = in
	return &dm.GroupInfoIndexResp{}, nil
}

func TestReadLogic_PropagatesWithDevices(t *testing.T) {
	stub := &stubDeviceGroup{}
	logic := NewReadLogic(context.Background(), &svc.ServiceContext{SvrClient: svc.SvrClient{DeviceG: stub}})

	_, err := logic.Read(&types.GroupInfoReadReq{
		ID:           1,
		Purpose:      "default",
		WithChildren: true,
		WithDevices:  true,
	})
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if stub.readReq == nil {
		t.Fatal("expected GroupInfoRead to be called")
	}
	if !stub.readReq.WithDevices {
		t.Fatal("expected WithDevices to be propagated to rpc request")
	}
}

func TestIndexLogic_PropagatesWithDevices(t *testing.T) {
	stub := &stubDeviceGroup{}
	logic := NewIndexLogic(context.Background(), &svc.ServiceContext{SvrClient: svc.SvrClient{DeviceG: stub}})

	_, err := logic.Index(&types.GroupInfoIndexReq{
		AreaID:      9,
		ParentID:    10,
		Name:        "group-name",
		Purpose:     "default",
		WithDevices: true,
	})
	if err != nil {
		t.Fatalf("Index() error = %v", err)
	}
	if stub.indexReq == nil {
		t.Fatal("expected GroupInfoIndex to be called")
	}
	if !stub.indexReq.WithDevices {
		t.Fatal("expected WithDevices to be propagated to rpc request")
	}
}
