package handler

import (
	"net/http"
	"nexus/internal/infrastructure/config/adapter/http/shared/response"
	"nexus/pkg/logger/version"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{
		"status":  "ok",
		"service": version.ServiceID,
		"time":    time.Now().Unix(),
	})
}
