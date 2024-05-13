package business

import (
	"cid.com/content-service/api/contract"
	"cid.com/content-service/api/services"
	"cid.com/content-service/backend"
	"cid.com/content-service/common/concurrency"
	"cid.com/content-service/common/secrets"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"net/http"
)

type Configuration struct {
	Listen struct {
		Grpc string
		Http string
	}
	FileStorage    backend.FileStorageConfiguration `yaml:"fileStorage"`
	ContentStorage backend.ContentStorageConfiguration
}

type Application interface {
	Run(ctx context.Context) error
}

type application struct {
	dependencies struct {
		storage       backend.ContentStorage
		fileStorage   backend.FileStorage
		queryService  services.MediaContentQueryService
		updateService services.MediaConteUpdateService
		httpService   services.HttpService
	}
	configurations struct {
		grpcAddr string
		httpAddr string
	}
}

func New(ctx context.Context, cfg Configuration, vault secrets.Secrets) (Application, error) {
	storage, err := backend.NewContentStorage(cfg.ContentStorage)
	if err != nil {
		slog.ErrorContext(ctx, "establishing connection to content storage failed", slog.String("err", err.Error()))
		return nil, err
	}
	fileStorage, err := backend.NewFileStorage(cfg.FileStorage, vault)
	if err != nil {
		slog.ErrorContext(ctx, "establishing connection to file storage failed", slog.String("err", err.Error()))
		return nil, err
	}
	queryService, err := services.NewMediaContentQueryService(storage)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create media query service", slog.String("err", err.Error()))
		return nil, err
	}
	updateService, err := services.NewMediaConteUpdateService(storage, fileStorage)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create media update service", slog.String("err", err.Error()))
		return nil, err
	}
	httpService, err := services.NewHttpService()
	if err != nil {
		slog.ErrorContext(ctx, "failed to create HTTP service", slog.String("err", err.Error()))
		return nil, err
	}
	return &application{
		dependencies: struct {
			storage       backend.ContentStorage
			fileStorage   backend.FileStorage
			queryService  services.MediaContentQueryService
			updateService services.MediaConteUpdateService
			httpService   services.HttpService
		}{
			storage:       storage,
			fileStorage:   fileStorage,
			queryService:  queryService,
			updateService: updateService,
			httpService:   httpService,
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
		slog.ErrorContext(ctx, "creating TCP listener", slog.String("err", err.Error()))
		return err
	}

	// gRPC service
	grpcServer := grpc.NewServer()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		contract.RegisterQueryMediaContentServer(grpcServer, a.dependencies.queryService)
		contract.RegisterModifyMediaContentServer(grpcServer, a.dependencies.updateService)

		slog.InfoContext(ctx, fmt.Sprintf("start to listen for grpc requests on %v", listener.Addr()))
		return concurrency.WhenDoneOrError(ctx, func(_ context.Context) error { return grpcServer.Serve(listener) })
	})

	// HTTP service
	httpServer := http.Server{
		Addr:    a.configurations.httpAddr,
		Handler: a.dependencies.httpService,
	}
	g.Go(func() error {
		slog.InfoContext(ctx, fmt.Sprintf("start to listen for HTTP requests on %v", httpServer.Addr))
		return concurrency.WhenDoneOrError(ctx, func(_ context.Context) error { return httpServer.ListenAndServe() })
	})

	// wait until error or shutdown signal
	if err = g.Wait(); err != nil {
		return err
	}

	// shutdown services and listeners
	slog.InfoContext(ctx, "stop listening for grpc requests")
	grpcServer.GracefulStop()
	slog.InfoContext(ctx, "stop listening for HTTP requests")
	return httpServer.Close()
}
