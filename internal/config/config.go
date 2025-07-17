package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlDB struct {
		DSN string
	}

	Sequence struct {
		DSN string
	}

	CacheRedis cache.CacheConf

	BaseString string // bas62指定基础字符串

	ShortUrlBlackList []string
	ShortDoamin       string
}
