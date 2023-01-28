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

	// Инициализация конфига
	err := initConfig()
	if err != nil {
		log.Fatalf("error initializing config: %s", err)
	}

	// Инициализация переменных окружения
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	// Проверка на новый адрес БД
	fmt.Println("Do you need me to change host address? (y/n) ")
	ans, host := "", ""
	if fmt.Scanln(&ans); ans == "y" {
		fmt.Println("Input custom host address:")
		fmt.Scanln(&host)
	} else {
		host = viper.GetString("db.host")
	}

	// Подключение к БД
	fmt.Println("Connecting to DB ...")
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

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping DB: %s", err)
	}

	// Проверка на инициализацию данных
	fmt.Println("Successful connection to DB")
	fmt.Println("Do you need me to initialize data? (y/n) ")
	if fmt.Scanln(&ans); ans == "y" {

		// Инициализация данных
		err = postgres.InitData(db)
		if err != nil {
			log.Fatalf("failed to init data: %s", err)
		}

	}

	// Получение имени
	fmt.Println("Please insert name")
	fmt.Scanln(&ans)

	// Вывод данных
	fmt.Println(postgres.GetName(db, ans))

	fmt.Println("Do you need me to add name? (y/n) ")
	if fmt.Scanln(&ans); ans == "y" {
		err = postgres.AddName(db, postgres.Name{
			Id:               8,
			Name:             "Danila",
			Meaning:          "-",
			Gender:           "Male",
			Origin:           "Russian",
			PeoplesCount:     1,
			WhenPeoplesCount: "2023-01-28T28 : 13 : 25.467",
		})
		if err != nil {
			log.Fatal("failed to add data: %s", err)
	fmt.Println("\nDone.")
}

// Функция инициализации конфига
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
