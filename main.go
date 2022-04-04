package main

import (
	"fmt"
	"os"
	"postgres-go/postgres"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	log.SetFormatter(new(log.JSONFormatter))

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	fmt.Println("Do you need me to initialize data? (y/n) ")
	ans, host := "", ""

	if fmt.Scanln(&ans); ans == "y" {

		fmt.Println("Do you need me to change host adress? (y/n) ")

		if fmt.Scanln(&ans); ans == "y" {
			fmt.Println("Please input custom host adress:")
			fmt.Scanln(&host)
		} else {
			host = viper.GetString("db.host")
		}

		db, err := sqlx.Open("postgres",
			fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
				host,
				viper.GetString("db.port"),
				viper.GetString("db.username"),
				viper.GetString("db.dbname"),
				os.Getenv("DB_PASSWORD"),
				viper.GetString("db.sslmode")))

		if err != nil {
			log.Fatalf("error opening DB: %s", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("failed to ping DB: %s", err)
		}

		if err := postgres.InitData(db); err != nil {
			log.Fatalf("failed to init data: %s", err)

		}
	}

	fmt.Println("dONE")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
