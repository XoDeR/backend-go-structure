package router

import "nexus/internal/adapter/http/v1/handler"

func InitHealthModule() *HealthRouter {
	handler := handler.NewHealthHandler()
	return NewHealthRouter(handler)
}
