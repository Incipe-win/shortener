package logic

import (
	"context"
	"database/sql"
	"errors"
	"shortener/pkg/base62"

	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Convert 转链：输一个长链接 --> 转为短链接
func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1. 校验数据（使用validator）
	// 1.1 数据不能空
	// 1.2 输入的长链接必须是一能请求通的网址
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("invalid long URL")
	}
	// 1.3 判断之前是否已经转链过（数据库中是否已经存在该长链接）
	// 1.3.1 给长链接生成 md5
	md5Hash := md5.Sum([]byte(req.LongUrl))
	// 1.3.2 拿 md5 去数据库中查是否存在
	u, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Hash, Valid: true})
	if err == nil {
		return &types.ConvertResponse{
			ShortUrl: l.svcCtx.Config.ShortDoamin + "/" + u.Surl.String,
		}, nil
	}
	if err != sqlx.ErrNotFound {
		logx.Errorw("failed to find long URL by md5", logx.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}
	// 1.4 输入不能是一个短链接（避免循环转链）
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool.GetBasePath failed", logx.LogField{Key: "lurl", Value: req.LongUrl}, logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, errors.New("url already is a short link")
		}
		logx.Errorw("failed to find short URL by base path", logx.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}
	var shortUrl string
	for {
		// 2. 取号
		// 每来一个转链请求，使用 replace into 语句往 sequence 表中插入一条数据，并取出主键 id 作为号码
		res, err := l.svcCtx.SequenceModel.ReplaceInto(l.ctx, &model.Sequence{Stub: "a"})
		if err != nil {
			logx.Errorw("failed to replace into sequence", logx.LogField{Key: "error", Value: err.Error()})
			return nil, err
		}
		seq, err := res.LastInsertId()
		if err != nil {
			logx.Errorw("failed to get last insert id from sequence", logx.LogField{Key: "error", Value: err.Error()})
			return nil, err
		}
		// 3. 号码转短链
		// 3.1 安全性，打乱 basestring 顺序
		shortUrl = base62.Int2String(uint64(seq))
		// 3.2 短域名黑名单，避免某些特殊词比如 api、fuck等
		if _, ok := l.svcCtx.ShortUrlBlackList[shortUrl]; !ok {
			break
		}
	}
	// 4. 存储长链接短链接映射关系
	logx.Debugf("short URL generated: %s", shortUrl)
	if _, err := l.svcCtx.ShortUrlModel.Insert(l.ctx, &model.ShortUrlMap{
		Lurl: sql.NullString{String: req.LongUrl, Valid: true},
		Md5:  sql.NullString{String: md5Hash, Valid: true},
		Surl: sql.NullString{String: shortUrl, Valid: true},
	}); err != nil {
		logx.Errorw("failed to insert short URL map", logx.LogField{Key: "error", Value: err.Error()})
	}
	// 5. 返回响应
	// 5.1 返回的是 短域名+短链接  q1mi.cn/1En
	shortUrl = l.svcCtx.Config.ShortDoamin + "/" + shortUrl
	return &types.ConvertResponse{
		ShortUrl: shortUrl,
	}, nil
}
