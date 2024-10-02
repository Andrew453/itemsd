package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) GetItems(c echo.Context) (err error) {
	itemsIDstr := c.Request().Header.Get("Items-IDs")
	itemsIDs, err := validateItemsID(itemsIDstr)
	if err != nil {
		slog.Error(err.Error())
		return c.NoContent(http.StatusNotFound)
	}

	result, err := h.uc.GetItems(itemsIDs)
	if err != nil {
		slog.Error(err.Error())
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, result)
}

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
