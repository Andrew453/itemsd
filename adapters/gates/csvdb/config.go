package csvdb

// Config Конфигурация шлюза работы с данными из CSV
type Config struct {
	FilePath string // Путь к файлу CSV
	// Порядковый номер столбца, который хранит идентификаторы itemID.
	// В данном случае поле id. Отсчет начинается с 0
	NumID int
}
