package logic

import (
	"context"
	"database/sql"
	"errors"

	"shortener/internal/svc"
	"shortener/internal/types"
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
	_, err = l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Hash, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {
			return nil, errors.New("long URL already exists")
		}
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
	// 2. 取号
	// 3. 号码转短链
	// 4. 存储长链接短链接映射关系
	// 5. 返回响应
	return
}
