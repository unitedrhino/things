package config

import "github.com/tal-tech/go-zero/rest"

type Config struct {
	rest.RestConf
	Rej struct {
		AccessSecret string
		AccessExpire int64
	}

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
