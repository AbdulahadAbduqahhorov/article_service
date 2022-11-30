package handlers

import (
	"strconv"

	"github.com/AbdulahadAbduqahhorov/gin/Article/config"
	"github.com/AbdulahadAbduqahhorov/gin/Article/storage"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Stg storage.StorageI
	Cfg config.Config
}

func NewHandler(strg storage.StorageI, cfg config.Config) Handler {
	return Handler{
		Stg: strg,
		Cfg: cfg,
	}
}

func (h *Handler) getOffsetParam(c *gin.Context) (offset int, err error) {
	offsetStr := c.DefaultQuery("offset", h.Cfg.DefaultOffset)
	return strconv.Atoi(offsetStr)
}

func (h *Handler) getLimitParam(c *gin.Context) (offset int, err error) {
	offsetStr := c.DefaultQuery("limit", h.Cfg.DefaultLimit)
	return strconv.Atoi(offsetStr)
}
