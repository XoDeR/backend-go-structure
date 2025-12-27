package router

import "github.com/gin-gonic/gin"

type V1Router struct {
	healthRouter *HealthRouter
}

func NewV1Router(
	healthRouter *HealthRouter,
) *V1Router {
	return &V1Router{
		healthRouter: healthRouter,
	}
}

func (r *V1Router) Setup(rg *gin.RouterGroup) {
	r.healthRouter.Setup(rg)
}
