package svc

import (
	"shortener/internal/config"
	"shortener/model"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	ShortUrlModel     model.ShortUrlMapModel
	SequenceModel     model.SequenceModel
	ShortUrlBlackList map[string]struct{}

	// bloom filter
	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	shortUrlModelConn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	sequenceModelConn := sqlx.NewMysql(c.Sequence.DSN)

	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}

	store := redis.New(c.CacheRedis[0].Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})

	shortUrlModel := model.NewShortUrlMapModel(shortUrlModelConn, c.CacheRedis)

	// 初始化 bloom filter
	filter := bloom.New(store, "bloom_filter", 20*(1<<20))
	loadDataToBloomFilter(filter, shortUrlModel)

	return &ServiceContext{
		Config:            c,
		ShortUrlModel:     shortUrlModel,
		SequenceModel:     model.NewSequenceModel(sequenceModelConn, c.CacheRedis),
		ShortUrlBlackList: m,
		Filter:            filter,
	}
}

func loadDataToBloomFilter(filter *bloom.Filter, model model.ShortUrlMapModel) {
	urls, err := model.FindAll()
	if err != nil {
		panic("failed to load data into bloom filter: " + err.Error())
	}
	for _, url := range urls {
		filter.Add([]byte(url))
	}
}
