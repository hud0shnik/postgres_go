package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
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

	log.SetFormatter(new(log.JSONFormatter))

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err)
	}

	names := make([]Name, 25897)
	data := make([]byte, 6124971)

	file, err := os.Open("names.json")
	if err != nil {
		log.Fatalf("error opening names.json file: %s", err)
	}

	_, err = file.Read(data)
	if err != nil {
		log.Fatalf("error reading names.json file: %s", err)
	}

	if err := json.Unmarshal(data, &names); err != nil {
		log.Fatalf("error unmarshaling structers from names.json fie: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	config := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.username"),
		viper.GetString("db.dbname"),
		os.Getenv("DB_PASSWORD"),
		viper.GetString("db.sslmode"))

	db, err := sqlx.Open("postgres", config)
	if err != nil {
		log.Fatalf("error opening DB: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping DB: %s", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("transaction error: %s", err)
	}

	for _, name := range names {
		_, err = tx.Exec("INSERT INTO names (id, name, meaning, gender, origin, peoplescount, whenpeoplescount) values ($1, $2, $3, $4, $5, $6, $7)",
			name.Id, name.Name, name.Meaning, name.Gender, name.Origin, name.PeoplesCount, name.WhenPeoplesCount)
		if err != nil {
			tx.Rollback()
			log.Fatalf("error inserting into db: %s", err)
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
