// main - пакет с точкой входа приложения
package main

import (
	"os"
	"pora/internal/handlers"
	"pora/internal/infrastructure/database"
	"pora/internal/infrastructure/logger"
	"pora/internal/infrastructure/storage"
	"pora/internal/middleware"
	"pora/internal/repositories"
	"pora/internal/services"

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

	// Подключение хранилища
	storage := storage.NewLocalStorage(os.Getenv("BASE_PATH"))

	// Подключение репозиториев
	userRepo := repositories.NewUserRepository(db)
	otpRepo := repositories.NewOtpRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	memberRepo := repositories.NewMemberRepository(db)
	listRepo := repositories.NewListRepository(db)
	itemRepo := repositories.NewItemRepository(db)

	// Подключение сервисов
	tokenService := services.NewTokenService(userRepo, refreshTokenRepo)
	fileService := services.NewFileService(storage)
	authService := services.NewAuthService(userRepo, refreshTokenRepo, tokenService)
	otpService := services.NewOTPService(otpRepo, userRepo, tokenService)
	listService := services.NewListService(memberRepo, listRepo, itemRepo)
	userService := services.NewUserService(userRepo, fileService, listService)

	// Подключение хэндлеров
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	otpHandler := handlers.NewOtpHandler(otpService)
	userHandler := handlers.NewUserHandler(userService, tokenService)

	// Регистрация маршрутов
	api := r.Group("/api")

	authorize := api.Group("/authorize")
	{
		authorize.GET("/check-user", authHandler.CheckUser)
		authorize.GET("/refresh", authHandler.RefreshTokens)

		authorize.POST("/send-otp", otpHandler.SendOTP)
		authorize.POST("/verify-otp", otpHandler.VerifyOTP)
	}

	user := api.Group("/user")
	{
		user.POST("/update", userHandler.UpdateUserInfo)
		user.POST("/save-image", userHandler.SaveImage)
	}

	// Маршруты для получения изображений
	uploads := r.Group("/uploads")
	{
		uploads.Static("/profiles", "./uploads/profiles")
	}

	// Запуск сервера
	port := os.Getenv("PORT")
	logger.Log.Info("Сервер запущен на порту :" + port)

	if err := r.Run(":" + port); err != nil {
		logger.Log.Error(err)
	}
}
