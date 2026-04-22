package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"shortener/internal/config"
	"shortener/internal/consumer"
	"shortener/internal/handler"
	"shortener/internal/svc"
	"shortener/pkg/base62"
	"shortener/pkg/mq"
	myotel "shortener/pkg/otel"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
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

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册 Prometheus metrics 端点
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/metrics",
		Handler: promhttp.Handler().ServeHTTP,
	})

	// 注册 pprof 性能分析端点（仅开发/压测时启用）
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/",
		Handler: http.DefaultServeMux.ServeHTTP,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/debug/pprof/:profile",
		Handler: http.DefaultServeMux.ServeHTTP,
	})

	// 使用 ServiceGroup 管理 HTTP Server 和 Kafka Consumers 的生命周期
	group := service.NewServiceGroup()
	defer group.Stop()

	// 注册 HTTP Server
	group.Add(server)

	// 注册 Kafka Consumers（仅在 Kafka 启用时）
	if c.Kafka.Enabled {
		brokers := c.Kafka.Brokers

		// AI 分析消费者
		if ctx.LLMClient != nil {
			aiConsumer := mq.NewKafkaConsumer(
				brokers,
				c.Kafka.Topics.AIAnalysis,
				"shortener-ai-analysis",
				consumer.AIAnalysisHandler(ctx.LLMClient, ctx.ShortUrlModel),
			)
			group.Add(aiConsumer)
			logx.Info("Kafka AI analysis consumer registered")
		}

		// 点击事件消费者
		clickConsumer := mq.NewKafkaConsumer(
			brokers,
			c.Kafka.Topics.ClickEvent,
			"shortener-click-events",
			consumer.ClickEventHandler(ctx.ShortUrlModel),
		)
		group.Add(clickConsumer)
		logx.Info("Kafka click event consumer registered")

		// 安全告警消费者
		safetyConsumer := mq.NewKafkaConsumer(
			brokers,
			c.Kafka.Topics.SafetyAlert,
			"shortener-safety-alerts",
			consumer.SafetyAlertHandler(),
		)
		group.Add(safetyConsumer)
		logx.Info("Kafka safety alert consumer registered")
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	group.Start()
}
