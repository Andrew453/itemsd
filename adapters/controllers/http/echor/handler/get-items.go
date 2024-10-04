package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// GetItems Обработчик запросов на получения информация об объектах item
// /get-items
func (h *Handler) GetItems(c echo.Context) (err error) {
	itemsIDstr := c.Request().Header.Get("Items-IDs")
	itemsIDs, err := validateItemsID(itemsIDstr)

	resp := struct {
		Result string `json:"result,omitempty"`
		Data   any    `json:"data,omitempty"`
	}{}

	if err != nil {
		slog.Error(err.Error())
		resp.Result = ErrorResult
		resp.Data = err.Error()
		return c.JSON(http.StatusNotFound, resp)
	}

	result, err := h.uc.GetItems(itemsIDs)
	if err != nil {
		slog.Error(err.Error())
		resp.Result = ErrorResult
		resp.Data = err.Error()
		return c.JSON(http.StatusNotFound, resp)
	}
	if string(result) == "[]" {
		resp.Result = ErrorResult
		resp.Data = "Items not found"
		return c.JSON(http.StatusNotFound, resp)

	}
	resp.Result = OkResult
	resp.Data = result
	return c.JSON(http.StatusOK, resp)
}

// validateItemsID проверка значение заголовка Items-IDs на валидность и парсинг значений itemsIDs
func validateItemsID(itemsIDstr string) (itemsIDs []int, err error) {
	defer func() {
		err = errors.Wrap(err, "validateItemsID()")
	}()
	if itemsIDstr == "" {
		return nil, errors.New("itemsIDs is empty")
	}
	ids := strings.Split(itemsIDstr, ",")
	if len(ids) == 0 {
		return nil, errors.New("itemsIDs is empty")
	}
	for _, str := range ids {
		id, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		itemsIDs = append(itemsIDs, id)
	}
	if len(itemsIDs) == 0 {
		return nil, errors.New("itemsIDs is invalid")
	}
	return itemsIDs, nil
}
