package config

import (
	"github.com/hashicorp/hcl"
	"github.com/pkg/errors"
	"os"
)

// Config Общая конфигурация приложения
type Config struct {
	Address    string `hcl:"address"`               // адрес сервера
	TLSEnabled bool   `hcl:"tls_enabled"`           // признак использования TLS
	LogFile    string `hcl:"log_file"`              // путь к файлу лога
	LogLevel   int    `hcl:"log_level"`             // уровень логирования
	DBPath     string `hcl:"db_path"`               // путь к csv-файлу базы данных
	CertFile   string `hcl:"cert_file"`             // путь к файлу сертификата
	KeyFile    string `hcl:"key_file"`              // путь к файлу ключа
	MaxItems   int    `hcl:"max_items_per_request"` // максимальное количество элементов в запросе
}

// Validate Проверка конфигурации приложения
func (c *Config) Validate() error {
	if c.Address == "" {
		return errors.New("empty address")
	}
	if c.DBPath == "" {
		return errors.New("empty db_path")
	}
	if c.MaxItems < 1 {
		return errors.New("empty max_items_per_request")
	}
	return nil
}

// LoadConfigHCL Загрузка конфигурации из HCL файла
func LoadConfigHCL(cfgFrom []byte, cfgTo interface{}) (err error) {
	defer func() { err = errors.Wrap(err, "pkg LoadConfigHCL()") }()

	err = hcl.Unmarshal(cfgFrom, cfgTo)
	if err != nil {
		err = errors.Wrap(err, "hcl.Unmarshal(cfgFrom, cfgTo)")
		return err
	}
	return nil
}

// LoadConfigHCLFromFile Загрузка конфигурации из HCL файла
func LoadConfigHCLFromFile(cfgPath string, cfgTo interface{}) (err error) {
	defer func() { err = errors.Wrap(err, "pkg LoadConfigHCLFromFile()") }()

	cfgFrom, err := os.ReadFile(cfgPath)
	if err != nil {
		err = errors.Wrap(err, "os.ReadFile(cfgPath)")
		return err
	}

	err = LoadConfigHCL(cfgFrom, cfgTo)
	if err != nil {
		return err
	}
	return nil
}
