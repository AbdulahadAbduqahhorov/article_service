package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/AbdulahadAbduqahhorov/gin/Article/api/docs"
	"github.com/AbdulahadAbduqahhorov/gin/Article/api/handlers"
	"github.com/AbdulahadAbduqahhorov/gin/Article/config"
)

// @contact.url   http://example.com
// @contact.email example@swagger.io
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func SetUpApi(r *gin.Engine, h handlers.Handler, cfg config.Config) {
	docs.SwaggerInfo.Title = cfg.App
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = cfg.Version

	v1 := r.Group("v1")
	{
		v1.GET("/article", h.GetArticle)
		v1.POST("/article", h.CreateArticle)
		v1.PUT("/article", h.UpdateArticle)
		v1.DELETE("/article/:id", h.DeleteArticle)
		v1.GET("/article/:id", h.GetArticleById)

		v1.GET("/author", h.GetAuthor)
		v1.POST("/author", h.CreateAuthor)
		v1.PUT("/author", h.UpdateAuthor)
		v1.DELETE("/author/:id", h.DeleteAuthor)
		v1.GET("/author/:id", h.GetAuthorById)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
