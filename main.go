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

	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			viper.GetString("db.host"),
			viper.GetString("db.port"),
			viper.GetString("db.username"),
			viper.GetString("db.dbname"),
			os.Getenv("DB_PASSWORD"),
			viper.GetString("db.sslmode")))

	if err != nil {
		log.Fatalf("error opening DB: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping DB: %s", err)
	}

	if err := postgres.InitData(db); err != nil {
		log.Fatalf("failed to init data: %s", err)

	}
	fmt.Println("dONE")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
