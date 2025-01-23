package devices

import (
	"context"
	"testing"
)

func TestGenMsgToken(t *testing.T) {
	type args struct {
		ctx    context.Context
		nodeID int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{
			ctx:    context.Background(),
			nodeID: 0x12feee53,
		}},
		{args: args{
			ctx:    context.Background(),
			nodeID: 0x12feee53,
		}},
		{args: args{
			ctx:    context.Background(),
			nodeID: 0x12feee53,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenMsgToken(tt.args.ctx, tt.args.nodeID); got != tt.want {
				t.Logf("GenMsgToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
