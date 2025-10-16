package main

import (
	"flag"
	"os"

	"mshop/pkg/nacosx"
	"mshop/service/order/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	env   string
	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&env, "env", "dev", "config path, eg: -env dev")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		// "ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		// "service.id", id,
		// "service.name", Name,
		// "service.version", Version,
		// "trace.id", tracing.TraceID(),
		// "span.id", tracing.SpanID(),
	)

	sources, err := nacosx.NewNacosConfigSource(
		nacosx.WithNamespace("mshop"),
		nacosx.WithEnv(env),
		nacosx.WithDataIds("order.yaml"),
	)

	if err != nil {
		panic(err)
	}

	c := config.New(
		config.WithSource(sources...),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	defer c.Close()

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Services, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
