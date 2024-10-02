package usecase

import "github.com/pkg/errors"

// Config Структура конфигурации бизнес-логики
type Config struct {
	MaxItemsIDs int // Максимальное количество items, которые будут обработаны в одном запросе
}

func (c *Config) Validate() error {
	if c.MaxItemsIDs <= 0 {
		return errors.New("MaxItemsIDs must be > 0")
	}
	return nil
}
