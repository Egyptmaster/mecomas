package services

import (
	"cid.com/content-service/api/contract"
	"cid.com/content-service/backend"
	"context"
)

type MediaContentQueryService interface {
	contract.QueryMediaContentServer
}

type mediaContentQueryService struct {
	contract.QueryMediaContentServer
	storage backend.ContentStorage
}

func NewMediaContentQueryService(s backend.ContentStorage) (MediaContentQueryService, error) {
	return &mediaContentQueryService{storage: s}, nil
}

// MediaItems lists all media items ordered by their upload time which matching the given filter criteria
func (m *mediaContentQueryService) MediaItems(*contract.ListMediaRequest, contract.QueryMediaContent_MediaItemsServer) error {
	return nil
}

// MediaItem returns a concrete media item for the given id
func (m *mediaContentQueryService) MediaItem(context.Context, *contract.Id) (*contract.MediaItemDetails, error) {
	return nil, nil
}

// Like returns the detailed information of a concrete like
func (m *mediaContentQueryService) Like(context.Context, *contract.Id) (*contract.LikeDetails, error) {
	return nil, nil
}

// Likes returns a detailed list of all likes ordered by time given for the specific media item
func (m *mediaContentQueryService) Likes(*contract.ListLikesRequest, contract.QueryMediaContent_LikesServer) error {
	return nil
}

// Comment returns the detailed list of a concrete comment
func (m *mediaContentQueryService) Comment(context.Context, *contract.Id) (*contract.CommentDetails, error) {
	return nil, nil
}

// Comments returns a detailed list of all comments ordered by time given for the specific media item
func (m *mediaContentQueryService) Comments(*contract.ListCommentsRequest, contract.QueryMediaContent_CommentsServer) error {
	return nil
}
