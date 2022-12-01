package main

import (
	"fmt"
	"net"

	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/config"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/genproto/article_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/genproto/author_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/service"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/storage"
	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/storage/postgres"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
)

func main() {
	var stg storage.StorageI
	var err error
	cfg := config.Load()
	stg, err = postgres.NewPostgres(fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase))
	if err != nil {
		panic(err)
	}
	authorSrv := service.NewAuthorService(stg)
	articleSrv:=service.NewArticleService(stg)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Error("error while listening: %v", err)
		return
	}

	service := grpc.NewServer()
	author_service.RegisterAuthorServiceServer(service, authorSrv)
	article_service.RegisterArticleServiceServer(service,articleSrv)
	if err := service.Serve(lis); err != nil {
		log.Error("error while listening: %v", err)
	}

	// h := handlers.NewHandler(stg,cfg)
	// switch cfg.Environment {
	// case "dev":
	// 	gin.SetMode(gin.DebugMode)
	// case "test":
	// 	gin.SetMode(gin.TestMode)
	// default:
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	// router := gin.New()

	// if cfg.Environment != "production" {
	// 	router.Use(gin.Logger(), gin.Recovery())
	// }

	// api.SetUpApi(router, h, cfg)

	// router.Run(cfg.HTTPPort)
}
