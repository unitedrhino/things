package devicemsglogic

import (
	"testing"

	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

func TestNormalizeSendLogIndexPageDefaultsOrdinaryMissingPageTo20(t *testing.T) {
	page := normalizeSendLogIndexPage(&dm.SendLogIndexReq{})

	if page.GetPage() != 1 || page.GetSize() != 20 {
		t.Fatalf("expected page=1 size=20, got page=%d size=%d", page.GetPage(), page.GetSize())
	}
}

func TestNormalizeSendLogIndexPageDefaultsFirstPropertyControlSendLookupTo1(t *testing.T) {
	page := normalizeSendLogIndexPage(&dm.SendLogIndexReq{
		ProductID:  "product-1",
		DeviceName: "device-1",
		Actions:    []string{"propertyControlSend"},
		ResultCode: 200,
		DataIDs:    []string{"hc_on"},
	})

	if page.GetPage() != 1 || page.GetSize() != 1 {
		t.Fatalf("expected page=1 size=1, got page=%d size=%d", page.GetPage(), page.GetSize())
	}
}

func TestNormalizeSendLogIndexPagePreservesExplicitSize(t *testing.T) {
	tests := []struct {
		name string
		page int64
		size int64
	}{
		{name: "device log page", page: 2, size: 20},
		{name: "caller large page", page: 3, size: 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := normalizeSendLogIndexPage(&dm.SendLogIndexReq{
				ProductID:  "product-1",
				DeviceName: "device-1",
				Actions:    []string{"propertyControlSend"},
				ResultCode: 200,
				DataIDs:    []string{"hc_on"},
				Page: &dm.PageInfo{
					Page: tt.page,
					Size: tt.size,
				},
			})

			if page.GetPage() != tt.page || page.GetSize() != tt.size {
				t.Fatalf("expected page=%d size=%d, got page=%d size=%d", tt.page, tt.size, page.GetPage(), page.GetSize())
			}
		})
	}
}

func TestNormalizeSendLogIndexPageDoesNotTreatOtherQueriesAsFirstLookup(t *testing.T) {
	tests := []struct {
		name string
		req  *dm.SendLogIndexReq
	}{
		{
			name: "other action",
			req: &dm.SendLogIndexReq{
				ProductID:  "product-1",
				DeviceName: "device-1",
				Actions:    []string{"statusLog"},
				ResultCode: 200,
				DataIDs:    []string{"hc_on"},
			},
		},
		{
			name: "missing dataIDs",
			req: &dm.SendLogIndexReq{
				ProductID:  "product-1",
				DeviceName: "device-1",
				Actions:    []string{"propertyControlSend"},
				ResultCode: 200,
			},
		},
		{
			name: "non-success result",
			req: &dm.SendLogIndexReq{
				ProductID:  "product-1",
				DeviceName: "device-1",
				Actions:    []string{"propertyControlSend"},
				ResultCode: 500,
				DataIDs:    []string{"hc_on"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := normalizeSendLogIndexPage(tt.req)

			if page.GetPage() != 1 || page.GetSize() != 20 {
				t.Fatalf("expected page=1 size=20, got page=%d size=%d", page.GetPage(), page.GetSize())
			}
		})
	}
}
