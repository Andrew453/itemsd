package usecase

import (
	"encoding/json"
	"github.com/pkg/errors"
	"prjs/itemsd/adapters/gates/csvdb"
)

// Usecase Реализация бизнес-логики. Здесь осуществляется вся необходимая работа с данными, шлюзами к данным
type Usecase struct {
	db     *csvdb.CSVDB // Шлюз к данным
	config Config       // конфигурация
}

// NewUsecase Конструктор
func NewUsecase(db *csvdb.CSVDB, config Config) (uc *Usecase, err error) {
	uc = &Usecase{
		db:     db,
		config: config,
	}
	err = uc.config.Validate()
	if err != nil {
		return nil, err
	}
	return uc, nil
}

// GetItems Получение и обработка данных об items по списку id
func (u *Usecase) GetItems(ids []int) (result json.RawMessage, err error) {
	defer func() {
		err = errors.Wrap(err, "usecase (u *Usecase) GetItems()")
	}()
	if len(ids) > u.config.MaxItemsIDs {
		return nil, errors.New("too many items")
	}
	ids = removeDuplicate(ids)
	result, err = u.db.GetItems(ids)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Функция удаления дубликатов из среза данных
func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
