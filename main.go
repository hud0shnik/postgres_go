package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
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

func main() {

	initConfig()

	names := make([]Name, 25897)
	data := make([]byte, 6124971)

	file, err := os.Open("names.json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &names)
	if err != nil {
		log.Fatal(err)
	}

	config := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.username"),
		viper.GetString("db.dbname"),
		viper.GetString("db.password"),
		viper.GetString("db.sslmode"))

	db, err := sqlx.Open("postgres", config)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(names); i++ {
		_, err = tx.Exec("INSERT INTO names (id, name, meaning, gender, origin, peoplescount, whenpeoplescount) values ($1, $2, $3, $4, $5, $6, $7)",
			names[i].Id, names[i].Name, names[i].Meaning, names[i].Gender, names[i].Origin, names[i].PeoplesCount, names[i].WhenPeoplesCount)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	tx.Commit()
	fmt.Println("dONE")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
