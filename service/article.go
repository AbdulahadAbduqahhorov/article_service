package service

import (
	"context"

	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/genproto/article_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ArticleService struct {
	article_service.UnimplementedArticleServiceServer
	Stg storage.StorageI
}

func NewArticleService(stg storage.StorageI) *ArticleService {
	return &ArticleService{
		Stg: stg,
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, req *article_service.CreateArticleRequest) (*article_service.CreateArticleResponse, error) {

	_, err := s.Stg.Author().GetAuthorById(req.AuthorId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "method GetAuthorById: %v",err)

	}
	id := uuid.New().String()

	err = s.Stg.Article().CreateArticle(id,req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method CreateArticle: %v",err)
	}

	_, err = s.Stg.Article().GetArticleById(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method GetArticleById: %v",err)

	}

	return &article_service.CreateArticleResponse{
		Id: id,
		},nil
}
func (s *ArticleService) GetArticleList(ctx context.Context,req *article_service.GetArticleListRequest) (*article_service.GetArticleListResponse, error) {

	res, err := s.Stg.Article().GetArticleList(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method GetArticleList: %v",err)

	}
	return res,nil
}
func (s *ArticleService) GetArticleById(ctx context.Context,req *article_service.GetArticleByIdRequest) (*article_service.GetArticleByIdResponse, error) {

	res, err := s.Stg.Article().GetArticleById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method GetArticleById: %v",err)

	}
	return res,nil

}
func (s *ArticleService) UpdateArticle(ctx context.Context,req *article_service.UpdateArticleRequest) (*article_service.UpdateArticleResponse, error) {

	
	err := s.Stg.Article().UpdateArticle(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method UpdateArticle: %v",err)
	}

	_, err = s.Stg.Article().GetArticleById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method GetArticleById: %v",err)
	
	}

	
	return &article_service.UpdateArticleResponse{
		Status: "Updated",
	}, nil
}
func (s *ArticleService) DeleteArticle(ctx context.Context,req *article_service.DeleteArticleRequest) (*article_service.DeleteArticleResponse, error) {

	err := s.Stg.Article().DeleteArticle(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "method DeleteArticle : %v",err)
	
	}

	return &article_service.DeleteArticleResponse{
		Status: "Deleted",
	},nil
	

}
