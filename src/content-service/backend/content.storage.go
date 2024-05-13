package backend

import "log/slog"

type ContentStorageConfiguration struct {
}

type ContentStorage interface {
}

func NewContentStorage(cfg ContentStorageConfiguration) (ContentStorage, error) {
	slog.Info("establishing connection to content storage")
	return &contentStorage{}, nil
}

type contentStorage struct {
}
