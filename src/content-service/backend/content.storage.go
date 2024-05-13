package backend

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type ContentStorageConfiguration struct {
}

type ContentStorage interface {
	GetMediaItem(context.Context, string) (MediaItem, error)
}

func NewContentStorage(cfg ContentStorageConfiguration) (ContentStorage, error) {
	slog.Info("establishing connection to content storage")
	return &contentStorage{}, nil
}

type contentStorage struct {
}

func (c contentStorage) GetMediaItem(ctx context.Context, mediaId string) (MediaItem, error) {
	//TODO fake
	if mediaId == "unknown" {
		return MediaItem{}, errors.New("unknown media id")
	}

	return MediaItem{
		Id:         mediaId,
		UserId:     "me",
		UploadDate: time.Date(2024, 04, 30, 0, 0, 0, 0, time.UTC),
	}, nil
}
