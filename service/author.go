package service

import (
	"context"

	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/genproto/author_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *AuthorService) CreateAuthor(ctx context.Context, req *author_service.CreateAuthorRequest) (*author_service.CreateAuthorResponse, error) {
	id, err := s.Stg.Author().CreateAuthor(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method CreateAuthor: %v", err)

	}
	return &author_service.CreateAuthorResponse{Id: id}, nil
}

func (s *AuthorService) GetAuthor(ctx context.Context, req *author_service.GetAuthorRequest) (*author_service.GetAuthorResponse, error) {
	res, err := s.Stg.Author().GetAuthor(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method GetAuthor: %v", err)

	}
	return res, nil
}

func (s *AuthorService) GetAuthorById(ctx context.Context, req *author_service.GetAuthorByIdRequest) (*author_service.Author, error) {
	res, err := s.Stg.Author().GetAuthorById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method GetAuthorById: %v", err)

	}
	return res, nil
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, req *author_service.UpdateAuthorRequest) (*author_service.UpdateAuthorResponse, error) {
	err := s.Stg.Author().UpdateAuthor(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method UpdateAuthor: %v", err)

	}
	return &author_service.UpdateAuthorResponse{Status: "Updated"}, nil
}
func (s *AuthorService) DeleteAuthor(ctx context.Context, req *author_service.DeleteAuthorRequest) (*author_service.DeleteAuthorResponse, error) {
	err := s.Stg.Author().DeleteAuthor(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method DeleteAuthor: %v", err)

	}
	return &author_service.DeleteAuthorResponse{Status: "Deleted"}, nil
}
