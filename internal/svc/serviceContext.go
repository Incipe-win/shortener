package svc

import (
	"shortener/internal/config"
	"shortener/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	ShortUrlModel     model.ShortUrlMapModel
	SequenceModel     model.SequenceModel
	ShortUrlBlackList map[string]struct{}
}

func NewServiceContext(c config.Config) *ServiceContext {
	shortUrlModelConn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	sequenceModelConn := sqlx.NewMysql(c.Sequence.DSN)

	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}

	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     model.NewShortUrlMapModel(shortUrlModelConn, c.CacheRedis),
		SequenceModel:     model.NewSequenceModel(sequenceModelConn, c.CacheRedis),
		ShortUrlBlackList: m,
	}
}
