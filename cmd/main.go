// main - пакет с точкой входа приложения
package main

import (
	"os"
	"pora/internal/infrastructure/database"
	"pora/internal/infrastructure/logger"
	"pora/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// main запускает сервер для обработки Api запросов.
// В main подключаются middleware, логгер, БД, репозитории, сервисы, хэндлеры
func main() {
	// Создание движка Gin
	r := gin.Default()

	// Подключение Middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.ErrorHandler())

	// Подключение логгера
	logger.InitLogger()

	// Чтение .env файла
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal("Ошибка при загрузке .env файла: ", err)
	}

	// Подключение к БД
	db, err := database.NewPostgresConnection()
	if err != nil {
		logger.Log.Fatal("Ошибка подключения к БД: ", err)
	}
	logger.Log.Info("Успешно осуществлено подключение к БД")

	// Запуск миграций БД
	err = database.RunMigrations(db)
	if err != nil {
		logger.Log.Fatal("Ошибка применения миграций: ", err)
	}
	logger.Log.Info("Миграции успешно применены")

	// Запуск сервера
	port := os.Getenv("PORT")
	logger.Log.Info("Сервер запущен на порту :" + port)

	if err := r.Run(":" + port); err != nil {
		logger.Log.Error(err)
	}
}
