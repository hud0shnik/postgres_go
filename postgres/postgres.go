package postgres

import (
	"encoding/json"
	"os"

	"github.com/jmoiron/sqlx"
)

type Name struct {
	Id               int    `json:"id" db:"id"`
	Name             string `json:"name" binding:"required"`
	Meaning          string `json:"meaning" binding:"required"`
	Gender           string `json:"gender" binding:"required"`
	Origin           string `json:"origin" binding:"required"`
	PeoplesCount     int    `json:"PeoplesCount" binding:"required"`
	WhenPeoplesCount string `json:"WhenPeoplesCount" binding:"required"`
}

func InitData(db *sqlx.DB) error {

	data := make([]byte, 6124971)
	names := make([]Name, 25897)

	file, err := os.Open("names.json")
	if err != nil {
		return err
	}

	if _, err = file.Read(data); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &names); err != nil {
		return err

	}

	tx, err := db.Begin()
	if err != nil {
		return err

	}

	for _, name := range names {
		_, err = tx.Exec("INSERT INTO names (id, name, meaning, gender, origin, peoplescount, whenpeoplescount) values ($1, $2, $3, $4, $5, $6, $7)",
			name.Id, name.Name, name.Meaning, name.Gender, name.Origin, name.PeoplesCount, name.WhenPeoplesCount)
		if err != nil {
			tx.Rollback()
			return err

		}
	}

	tx.Commit()
	return nil
}
