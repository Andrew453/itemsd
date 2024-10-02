package csvdb

import (
	"encoding/csv"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// CSVDB Функционал работы с csv файлом
type CSVDB struct {
	filePath string // путь к файлу хранения данных
	headers  map[int]string
}

func NewCSVDB(config Config) (c *CSVDB, err error) {
	defer func() {
		err = errors.Wrap(err, "csvdb (c *CSVDB) NewCSVDB()")
	}()
	c = &CSVDB{}
	c.headers = make(map[int]string)
	c.filePath, err = filepath.Abs(config.FilePath)
	if err != nil {
		err = errors.Wrap(err, "filepath.Abs(confPath)")
		return
	}

	file, err := os.Open(c.filePath)
	if err != nil {
		err = errors.Wrap(err, "os.Open(c.filePath)")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	hs, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for i, h := range hs {
		c.headers[i] = h
	}
	return c, nil
}

// GetItems Функция получения данных из csv файла по списку id itemов
func (c *CSVDB) GetItems(ids []int) (items json.RawMessage, err error) {
	defer func() {
		err = errors.Wrap(err, "csvdb (c *CSVDB) GetItems()")
	}()
	file, err := os.Open(c.filePath)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open(c.filePath)")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var res = make([]interface{}, 0, len(ids))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		id, err := strconv.Atoi(record[1])
		if err != nil {
			continue
		}
		if CheckID(id, ids) {
			j, err := RecordToJSON(c.headers, record)
			if err != nil {
				return nil, err
			}
			res = append(res, j)
		}
	}
	items, err = json.Marshal(res)
	if err != nil {
		err = errors.Wrap(err, "json.Marshal(res)")
		return nil, err
	}

	return items, nil
}

// CheckID Функция проверки id на существование в списке ids
func CheckID(id int, ids []int) (exists bool) {
	for _, i := range ids {
		if id == i {
			return true
		}
	}
	return false
}

// RecordToJSON Функция преобразования записи из csv файла в json
func RecordToJSON(headers map[int]string, values []string) (j json.RawMessage, err error) {
	m := make(map[string]interface{})

	for key, value := range values {
		m[headers[key]] = value
	}
	j, err = json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return j, nil
}
