package pubJob

import (
	"context"
	"github.com/i-Things/things/shared/conf"
)

type (
	PubJob struct {
		natsJs *natsJsClient
	}
)

func NewPubJob(c conf.EventConf) (*PubJob, error) {
	natsJs, err := newNatsJsClient(c.Nats)
	if err != nil {
		return nil, err
	}
	pj := PubJob{natsJs: natsJs}
	return &pj, nil
}
func (p *PubJob) Publish(ctx context.Context, pubType string, topic string, payload []byte) error {
	if pubType == conf.EventModeNatsJs {
		return p.natsJs.Publish(ctx, topic, payload)
	}
	return nil
}
