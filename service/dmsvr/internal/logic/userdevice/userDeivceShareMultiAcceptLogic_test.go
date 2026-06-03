package userdevicelogic

import "testing"

func TestShouldConsumeShareTokenAfterAccept(t *testing.T) {
	tests := []struct {
		name  string
		useBy string
		want  bool
	}{
		{name: "wechat single device token is one-time", useBy: "wechat_single_device", want: true},
		{name: "family token remains reusable", useBy: "family", want: false},
		{name: "empty useBy remains reusable", useBy: "", want: false},
		{name: "unknown useBy remains reusable", useBy: "batch_qr", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldConsumeShareTokenAfterAccept(tt.useBy)
			if got != tt.want {
				t.Fatalf("shouldConsumeShareTokenAfterAccept(%q) = %v, want %v", tt.useBy, got, tt.want)
			}
		})
	}
}
