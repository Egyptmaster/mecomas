package services

import (
	"cid.com/content-service/api/contract"
	"cid.com/content-service/backend"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type MediaConteUpdateService interface {
	contract.ModifyMediaContentServer
}

type mediaConteUpdateService struct {
	contract.ModifyMediaContentServer
	storage     backend.ContentStorage
	fileStorage backend.FileStorage
}

func NewMediaConteUpdateService(s backend.ContentStorage, fs backend.FileStorage) (MediaConteUpdateService, error) {
	return &mediaConteUpdateService{storage: s, fileStorage: fs}, nil
}

// GeneratePreSignedUrl allows to request a temporary direct upload link for a media content
func (m *mediaConteUpdateService) GeneratePreSignedUrl(ctx context.Context, id *contract.Id) (*contract.PreSignedUrl, error) {
	slog.Debug("incoming request for GeneratePreSignedUrl")
	if err := id.Validate(); err != nil {

		slog.Error("bad request", slog.String("err", err.Error()))
		return nil, err
	}

	// the media item must exist before the url will be requested
	mediaItem, err := m.storage.GetMediaItem(ctx, id.Id)
	if err != nil {
		slog.Error("failed to generate pre-signed url", slog.String("err", err.Error()))
		return nil, status.Error(codes.NotFound, err.Error())
	}
	// generated the pre-signed url
	url, err := m.fileStorage.GeneratePreSignedUrl(mediaItem.Key())
	if err != nil {
		slog.Error("failed to generate pre-signed url", slog.String("err", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &contract.PreSignedUrl{Url: url}, nil
}

// AddOrUpdateMediaItem allows to either add a new item or modify an existing item
func (m *mediaConteUpdateService) AddOrUpdateMediaItem(context.Context, *contract.WriteMediaItem) (*contract.Id, error) {
	return nil, nil
}

// AddLike allows to either add a like or a dislike to a concrete media item
func (m *mediaConteUpdateService) AddLike(context.Context, *contract.WriteLike) (*contract.Id, error) {
	return nil, nil
}

// AddComment allows to add a comment to a concrete media item
func (m *mediaConteUpdateService) AddComment(context.Context, *contract.WriteComment) (*contract.Id, error) {
	return nil, nil
}

// DeleteMediaItem allows to remove a media item
func (m *mediaConteUpdateService) DeleteMediaItem(context.Context, *contract.Id) (*contract.Id, error) {
	return nil, nil
}

// DeleteLike allows to remove an given like
func (m *mediaConteUpdateService) DeleteLike(context.Context, *contract.Id) (*contract.Id, error) {
	return nil, nil
}

// DeleteComment allows to remove a comment
func (m *mediaConteUpdateService) DeleteComment(context.Context, *contract.Id) (*contract.Id, error) {
	return nil, nil
}
