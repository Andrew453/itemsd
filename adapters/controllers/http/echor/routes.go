package echor

import (
	"github.com/labstack/echo/v4"
	"prjs/itemsd/adapters/controllers/http/echor/handler"
)

// setRoutes Функция выставления обработки запросов по заданным путям
func setRoutes(e *echo.Echo, h *handler.Handler) {
	// e.POST...
	// e.DELETE...
	e.GET("/get-items", h.GetItems)
}
