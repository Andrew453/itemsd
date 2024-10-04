// Package csvdb Пакет с описанием функционала работы с CSV файлом
package csvdb

import (
	"encoding/csv"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
)

// CSVDB Функционал работы с csv файлом
type CSVDB struct {
	filePath string         // путь к файлу хранения данных
	headers  map[int]string // Имена столбцов данных в csv файле
	numID    int            // Порядковый номер поля ID, по которому будет осуществляться поиск
	file     *os.File
	mu       sync.Mutex
}

func NewCSVDB(config Config) (c *CSVDB, err error) {
	defer func() {
		err = errors.Wrap(err, "csvdb (c *CSVDB) NewCSVDB()")
	}()
	c = &CSVDB{}
	c.headers = make(map[int]string)
	c.filePath, err = filepath.Abs(config.FilePath)
	c.numID = config.NumID
	if err != nil {
		err = errors.Wrap(err, "filepath.Abs(confPath)")
		return
	}

	c.file, err = os.Open(c.filePath)
	if err != nil {
		err = errors.Wrap(err, "os.Open(c.filePath)")
		return
	}
	//defer file.Close()

	reader := csv.NewReader(c.file)

	hs, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for i, h := range hs {
		c.headers[i] = h
	}
	c.file.Seek(0, 0)
	return c, nil
}

// GetItems Функция получения данных из csv файла по списку id itemов
func (c *CSVDB) GetItems(ids []int) (items json.RawMessage, err error) {
	defer func() {
		err = errors.Wrap(err, "csvdb (c *CSVDB) GetItems()")
	}()
	c.mu.Lock()

	reader := csv.NewReader(c.file)

	var res = make([]interface{}, 0, len(ids))
	var exist bool
	for {

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		id, err := strconv.Atoi(record[c.numID])
		if err != nil {
			continue
		}
		ids, exist = CheckID(id, ids)
		if exist {
			j, err := RecordToJSON(c.headers, record)
			if err != nil {
				return nil, err
			}
			res = append(res, j)
		}
		if len(ids) == 0 {
			break
		}
	}
	c.file.Seek(0, 0)
	c.mu.Unlock()
	items, err = json.Marshal(res)
	if err != nil {
		err = errors.Wrap(err, "json.Marshal(res)")
		return nil, err
	}

	return items, nil
}

// GetItemsAsync Функция получения данных из csv файла по списку id itemов
// UNSAFE
func (c *CSVDB) GetItemsAsync(ids []int) (items json.RawMessage, err error) {
	defer func() {
		err = errors.Wrap(err, "csvdb (c *CSVDB) GetItems()")
	}()

	file, err := os.Open(c.filePath)
	if err != nil {
		err = errors.Wrap(err, "os.Open(c.filePath)")
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var res = make([]interface{}, 0, len(ids))
	var exist bool
	for {

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		id, err := strconv.Atoi(record[c.numID])
		if err != nil {
			continue
		}
		ids, exist = CheckID(id, ids)
		if exist {
			j, err := RecordToJSON(c.headers, record)
			if err != nil {
				return nil, err
			}
			res = append(res, j)
		}
		if len(ids) == 0 {
			break
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
// Если элемент найден в массиве, то он удаляется
func CheckID(id int, ids []int) (newIDs []int, exists bool) {
	var delElemNum int
	for j, i := range ids {
		if id == i {
			exists = true
			delElemNum = j
		}
	}
	if exists {
		ids = slices.Delete(ids, delElemNum, delElemNum+1)
	}
	return ids, exists
}

// RecordToJSON Функция преобразования записи из csv файла в json
func RecordToJSON(headers map[int]string, values []string) (j json.RawMessage, err error) {
	m := make(map[string]interface{})

	for key, value := range values {
		if key == 0 {
			continue
		}
		m[headers[key]] = value
	}
	j, err = json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (c *CSVDB) Stop() (err error) {
	defer func() {
		err = errors.Wrap(err, "csvdb (c *CSVDB) Stop()")
	}()
	if c.file != nil {
		err = c.file.Close()
		if err != nil {
			err = errors.Wrap(err, "c.file.Close()")
			return err
		}
	}
	return nil
}
