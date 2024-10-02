package echor

import "github.com/pkg/errors"

// Config Конфигурация HTTP сервера
type Config struct {
	TLSEnable bool   // Включение TLS
	Address   string // HTTP адрес работы сервера
}

// Validate Проверка введенных данных конфигурации на валидность
func (c *Config) Validate() (err error) {
	defer func() {
		err = errors.Wrap(err, "echor (c *Config) Validate()")
	}()
	if c.Address == "" {
		return errors.New("address is empty: c.Address == /'/'")
	}
	return nil
}
