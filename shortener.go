package main

import (
	"flag"
	"fmt"
	"net/http"

	"shortener/internal/config"
	"shortener/internal/handler"
	"shortener/internal/svc"
	"shortener/pkg/base62"
	myotel "shortener/pkg/otel"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/shortener-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	// base62模块的初始化
	base62.MustInit(c.BaseString)

	// 初始化 OpenTelemetry（如果配置了）
	if c.Otel.Endpoint != "" {
		shutdown, err := myotel.InitTracer(c.Otel.Name, c.Otel.Endpoint, c.Otel.Sampler)
		if err != nil {
			logx.Errorf("failed to init OpenTelemetry tracer: %v", err)
		} else {
			defer shutdown()
			logx.Info("OpenTelemetry tracer initialized")
		}
	}

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册 Prometheus metrics 端点
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/metrics",
		Handler: promhttp.Handler().ServeHTTP,
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
