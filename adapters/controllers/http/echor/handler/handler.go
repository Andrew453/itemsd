package handler

import "prjs/itemsd/model/usecase"

// Handler Структура - обработчик запросов.
// Имеет функционал первичной обработки http запросов и вызова реализации бизнеслогики
type Handler struct {
	uc usecase.Usecase
}

func NewHandler(uc usecase.Usecase) *Handler {
	return &Handler{uc: uc}
}
