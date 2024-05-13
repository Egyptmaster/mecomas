package services

import (
	"cid.com/content-service/api/contract"
	"cid.com/content-service/backend"
)

type MediaContentQueryService interface {
	contract.QueryMediaContentServer
}

type mediaContentQueryService struct {
	contract.QueryMediaContentServer
	storage backend.ContentStorage
}

func NewMediaContentQueryService(s backend.ContentStorage) (MediaContentQueryService, error) {
	return mediaContentQueryService{storage: s}, nil
}
