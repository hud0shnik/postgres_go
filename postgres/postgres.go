package postgres

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

// Структура имени
type Name struct {
	Id               int    `json:"id" db:"id"`
	Name             string `json:"name" binding:"required"`
	Meaning          string `json:"meaning" binding:"required"`
	Gender           string `json:"gender" binding:"required"`
	Origin           string `json:"origin" binding:"required"`
	PeoplesCount     int    `json:"PeoplesCount" binding:"required"`
	WhenPeoplesCount string `json:"WhenPeoplesCount" binding:"required"`
}

// Функция записи данных в бд
func InitData(db *sqlx.DB) error {

	// Слайс символов и слайс имён
	data := make([]byte, 6124971) //если не работает, попробовать 6113867 4429392 6346950
	names := make([]Name, 25897)

	// Открытие файла с именами
	file, err := os.Open("names.json")
	if err != nil {
		return err
	}

	// Запись текста файла в слайс символов
	if _, err = file.Read(data); err != nil {
		return err
	}

	// Перевод слайса символов в слайс имён
	if err := json.Unmarshal(data, &names); err != nil {
		return err

	}

	// Создание транзакции
	tx, err := db.Begin()
	if err != nil {
		return err

	}

	// Проход по всем именам
	for i, name := range names {

		// Вывод в консоль уведомления о добавлении имени в транзакцию
		fmt.Println("Inserting...\t", i, "/", 25897)

		// Добавление запроса к транзакции его проверка
		_, err = tx.Exec("INSERT INTO names (id, name, meaning, gender, origin, peoplescount, whenpeoplescount) values ($1, $2, $3, $4, $5, $6, $7)",
			name.Id, name.Name, name.Meaning, name.Gender, name.Origin, name.PeoplesCount, name.WhenPeoplesCount)
		if err != nil {
			tx.Rollback()
			fmt.Println("Exec error:", err)
			return err

		}

	}

	// Исполнение транзакции
	tx.Commit()

	return nil
}

// Функция получения имени из БД
func GetName(db *sqlx.DB, name string) Name {

	var item Name

	db.Get(&item, "SELECT * FROM names WHERE name = $1", name)

	return item
}

// Функция добавления имени в БД
func AddName(db *sqlx.DB, name Name) error {

	// Добавление имени в БД
	_, err := db.Exec("INSERT INTO names (id, name, meaning, gender, origin, peoplescount, whenpeoplescount) values ($1, $2, $3, $4, $5, $6, $7)",
		name.Id, name.Name, name.Meaning, name.Gender, name.Origin, name.PeoplesCount, name.WhenPeoplesCount)
	if err != nil {
		fmt.Println("Exec error:", err)
		return err

	}

	return nil

}
