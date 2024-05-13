package business

import (
	"cid.com/content-service/api/services"
	"cid.com/content-service/backend"
	"cid.com/content-service/common/concurrency"
	"cid.com/content-service/common/slog"
	"context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type Configuration struct {
	Listen struct {
		Grpc string
		Http string
	}
	ContentStorage backend.ContentStorageConfiguration
}

type Application interface {
	Run(ctx context.Context) error
}

type application struct {
	dependencies struct {
		storage      backend.ContentStorage
		queryService services.MediaContentQueryService
		httpService  services.HttpService
	}
	configurations struct {
		grpcAddr string
		httpAddr string
	}
}

func New(ctx context.Context, cfg Configuration) (Application, error) {
	storage, err := backend.NewContentStorage(cfg.ContentStorage)
	if err != nil {
		return nil, slog.ErrorContext(ctx, "establishing connection to content storage failed", err)
	}
	queryService, err := services.NewMediaContentQueryService(cfg.ContentStorage)
	if err != nil {
		return nil, slog.ErrorContext(ctx, "failed to create media query service", err)
	}
	httpService, err := services.NewHttpService()
	if err != nil {
		return nil, slog.ErrorContext(ctx, "failed to create HTTP service", err)
	}
	return &application{
		dependencies: struct {
			storage      backend.ContentStorage
			queryService services.MediaContentQueryService
			httpService  services.HttpService
		}{
			storage:      storage,
			queryService: queryService,
			httpService:  httpService,
		},
		configurations: struct {
			grpcAddr string
			httpAddr string
		}{
			grpcAddr: cfg.Listen.Grpc,
			httpAddr: cfg.Listen.Http,
		},
	}, nil
}

func (a *application) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", a.configurations.grpcAddr)
	if err != nil {
		return slog.ErrorContext(ctx, "creating TCP listener", err)
	}

	// gRPC service
	grpcServer := grpc.NewServer()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		slog.InfoContextf(ctx, "start to listen for grpc requests on %v", listener.Addr())
		return concurrency.WhenDoneOrError(ctx, func(_ context.Context) error { return grpcServer.Serve(listener) })
	})

	// HTTP service
	httpServer := http.Server{
		Addr:    a.configurations.httpAddr,
		Handler: a.dependencies.httpService,
	}
	g.Go(func() error {
		slog.InfoContextf(ctx, "start to listen for HTTP requests on %v", httpServer.Addr)
		return concurrency.WhenDoneOrError(ctx, func(_ context.Context) error { return httpServer.ListenAndServe() })
	})

	// wait until error or shutdown signal
	if err = g.Wait(); err != nil {
		return err
	}

	// shutdown services and listeners
	slog.InfoContextf(ctx, "stop listening for grpc requests")
	grpcServer.GracefulStop()
	slog.InfoContextf(ctx, "stop listening for HTTP requests")
	return httpServer.Close()
}
