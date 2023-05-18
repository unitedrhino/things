package oss

import (
	"context"
	"sync"

	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/oss/common"
)

type Client struct {
	Handle
}

var (
	client  *Client
	newOnce sync.Once
)

func NewOssClient(c conf.OssConf) *Client {
	newOnce.Do(func() {
		ossManager, err := newOssManager(c)
		if err != nil {
			panic(err)
		}
		client = &Client{
			ossManager,
		}
	})

	return client
}

type OpOption func(*common.OptionKv)

func (c *Client) getDefaultOption(ctx context.Context) OpOption {
	return func(option *common.OptionKv) {
		option.SetHttpParams("x-process", "xxxxx")
	}
}
