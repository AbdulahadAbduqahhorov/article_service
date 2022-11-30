package handlers

import (
	"errors"
	"net/http"

	"github.com/AbdulahadAbduqahhorov/gRPC/blogpost/article_service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateArticle godoc
// @Summary     Create article
// @Description create article
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       article body     models.CreateArticleModel true "article body"
// @Success     201     {object} models.JSONResult{data=string} "Success"
// @Failure     400     {object} models.JSONErrorResult "Bad request"
// @Failure     500     {object} models.JSONErrorResult "Server error"
// @Router      /v1/article [post]
func (h Handler) CreateArticle(c *gin.Context) {

	var article models.CreateArticleModel

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}
	_, err := h.Stg.Author().GetAuthorById(article.AuthorId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: errors.New("author not found").Error(),
		})
		return
	}
	id := uuid.New().String()

	err = h.Stg.Article().CreateArticle(id, article)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}

	_, err = h.Stg.Article().GetArticleById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.JSONResult{
		Message: "Article has been created",
		Data:    id,
	})
}

// GetArticleList godoc
// @Summary     List Article
// @Description get article
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       limit  query    int    false "10"
// @Param       offset query    int    false "0"
// @Param       search query    string false "string default"
// @Success     200    {object} models.JSONResult{data=[]models.Article} "Success"
// @Failure     400     {object} models.JSONErrorResult "Bad request"
// @Router      /v1/article [get]
func (h Handler) GetArticle(c *gin.Context) {
	search := c.Query("search")

	limit, err := h.getLimitParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}

	offset, err := h.getOffsetParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}
	res, err := h.Stg.Article().GetArticle(limit, offset, search)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.JSONResult{
		Message: "Article List",
		Data:    res,
	})
}

// GetArticleById godoc
// @Summary     Get article by id
// @Description get an article by id
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       id  path     string true "article id"
// @Success     200 {object} models.JSONResult{data=models.GetArticleByIdModel} "Success"
// @Failure     400 {object} models.JSONErrorResult "Bad request"
// @Router      /v1/article/{id} [get]
func (h Handler) GetArticleById(c *gin.Context) {
	id := c.Param("id")

	res, err := h.Stg.Article().GetArticleById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.JSONResult{
		Message: "OK",
		Data:    res,
	})
}

// UpdateArticle godoc
// @Summary     Update article
// @Description update article
// @Tags        articles
// @Accept      json
// @Produce     json
// @Param       article body     models.UpdateArticleModel true "article body"
// @Success     200     {object} models.JSONResult{message=string} "Success"
// @Failure     400     {object} models.JSONErrorResult "Bad request"
// @Failure     500     {object} models.JSONErrorResult "Server error"
// @Router      /v1/article [put]
func (h Handler) UpdateArticle(c *gin.Context) {
	var article models.UpdateArticleModel

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}

	err := h.Stg.Article().UpdateArticle(article)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}

	_, err = h.Stg.Article().GetArticleById(article.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.JSONResult{
		Message: "Article has been  Updated",
	})

}

// DeleteArticle godoc
// @Summary     Delete article
// @Description delete an article
// @Tags        articles
// @Produce     json
// @Param       id  path     string true "article id"
// @Success     200 {object} models.JSONResult{message=string} "Success"
// @Failure     400 {object} models.JSONErrorResult "Bad Request"
// @Router      /v1/article/{id} [delete]
func (h Handler) DeleteArticle(c *gin.Context) {

	id := c.Param("id")
	err := h.Stg.Article().DeleteArticle(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.JSONErrorResult{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.JSONResult{
		Message: "Article has been Deleted",
	})
}
