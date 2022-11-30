package service

import (
	"context"

	"github.com/AbdulahadAbduqahhorov/gin/Article/genproto/author_service"
	"github.com/AbdulahadAbduqahhorov/gin/Article/storage"
)

type AuthorService struct {
	author_service.UnimplementedAuthorServiceServer
	Stg storage.StorageI
}

func NewAuthorService(stg storage.StorageI) *AuthorService {
	return &AuthorService{
		Stg: stg,
	}
}

func (s *AuthorService) Create(ctx context.Context, req *author_service.CreateAuthorRequest) (*author_service.CreateAuthorResponse, error) {
	res, err := s.Stg.Author().CreateAuthor(*req)
	if err != nil {
		return &author_service.CreateAuthorResponse{}, err
	}
	return &res, nil
}

func (s *AuthorService) GetAuthor(ctx context.Context, req *author_service.GetAuthorRequest) (*author_service.GetAuthorResponse, error) {
	res, err := s.Stg.Author().GetAuthor(*req)
	if err != nil {
		return &author_service.GetAuthorResponse{}, err
	}
	return &res, nil
}

func (s *AuthorService) GetByIdAuthor(ctx context.Context, req *author_service.GetAuthorByIdResponse) (*author_service.Author, error) {
	res, err := s.Stg.Author().GetAuthorById(*req)
	if err != nil {
		return &author_service.Author{}, err
	}
	return &res, nil
}
