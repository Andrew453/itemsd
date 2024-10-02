package usecase

import "encoding/json"

// Usecase Описание функционала бизнес-логики
type Usecase interface {
	GetItems(ids []int) (result json.RawMessage, err error)
}
