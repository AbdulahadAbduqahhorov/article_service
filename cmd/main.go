package main

import (
	"fmt"
	"net"

	"github.com/AbdulahadAbduqahhorov/gin/Article/config"
	"github.com/AbdulahadAbduqahhorov/gin/Article/genproto/author_service"
	"github.com/AbdulahadAbduqahhorov/gin/Article/service"
	"github.com/AbdulahadAbduqahhorov/gin/Article/storage"
	"github.com/AbdulahadAbduqahhorov/gin/Article/storage/postgres"
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
	s := service.NewAuthorService(stg)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Error("error while listening: %v", err)
		return
	}

	service := grpc.NewServer()
	author_service.RegisterAuthorServiceServer(service,s)
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
