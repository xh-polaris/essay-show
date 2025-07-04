// Code generated by hertz generator.

package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/xh-polaris/essay-show/biz/adaptor"
	"github.com/xh-polaris/essay-show/biz/infrastructure/util/log"
	"github.com/xh-polaris/essay-show/provider"
	"github.com/xh-polaris/gopkg/hertz/middleware"
	logx "github.com/xh-polaris/gopkg/util/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

func Init() {
	provider.Init()
	hlog.SetLogger(logx.NewHlogLogger())
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(b3.New(), propagation.Baggage{}, propagation.TraceContext{}))
	http.DefaultTransport = otelhttp.NewTransport(http.DefaultTransport)
}

func main() {
	Init()
	c := provider.Get().Config

	tracer, cfg := tracing.NewServerTracer()
	h := server.New(
		server.WithHostPorts(c.ListenOn),
		server.WithTracer(prometheus.NewServerTracer(":9091", "/server/metrics")),
		tracer,
	)
	h.Use(tracing.ServerMiddleware(cfg), middleware.EnvironmentMiddleware, recovery.Recovery(), func(ctx context.Context, c *app.RequestContext) {
		ctx = adaptor.InjectContext(ctx, c)
		c.Next(ctx)
	})

	register(h)
	log.Info("server start")
	h.Spin()
}
