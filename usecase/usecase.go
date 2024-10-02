package usecase

import (
	"encoding/json"
	"github.com/pkg/errors"
	"prjs/itemsd/adapters/gates/csvdb"
)

type Usecase struct {
	db     *csvdb.CSVDB
	config Config
}

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

func (u *Usecase) GetItems(ids []int) (result json.RawMessage, err error) {
	defer func() {
		err = errors.Wrap(err, "usecase (u *Usecase) GetItems()")
	}()
	if len(ids) > u.config.MaxItemsIDs {
		return nil, errors.New("too many items")
	}

	result, err = u.db.GetItems(ids)
	if err != nil {
		return nil, err
	}
	return result, nil
}
