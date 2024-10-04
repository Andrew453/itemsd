// Package handler Пакет первичной обработки http запросов
package handler

import "prjs/itemsd/model/usecase"

// Handler Структура - обработчик запросов.
// Имеет функционал первичной обработки http запросов и вызова реализации бизнес-логики
type Handler struct {
	uc usecase.Usecase // интерфейс реализации бизнес-логики
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}
