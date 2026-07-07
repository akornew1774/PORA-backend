// main - пакет с точкой входа приложения
package main

import (
	"os"
	"pora/internal/config"
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

	// Загрузка конфигураций
	cfg := config.Load()

	// Подключение репозиториев
	userRepo := repositories.NewUserRepository(db)
	otpRepo := repositories.NewOtpRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	familyRepo := repositories.NewFamilyRepository(db)
	memberRepo := repositories.NewMemberRepository(db)
	listRepo := repositories.NewListRepository(db)
	itemRepo := repositories.NewItemRepository(db)

	// Подключение сервисов
	tokenService := services.NewTokenService(userRepo, refreshTokenRepo)
	fileService := services.NewFileService(storage)
	authService := services.NewAuthService(userRepo, refreshTokenRepo, tokenService)
	otpService := services.NewOTPService(otpRepo, userRepo, tokenService)
	listService := services.NewListService(userRepo, memberRepo, familyRepo, listRepo, itemRepo)
	userService := services.NewUserService(userRepo, fileService, listService)
	familyService := services.NewFamilyService(userRepo, memberRepo, familyRepo, listService, cfg.DeepLink)
	itemService := services.NewItemService(familyRepo, listRepo, itemRepo)

	// Подключение хэндлеров
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	otpHandler := handlers.NewOtpHandler(otpService)
	userHandler := handlers.NewUserHandler(userService, tokenService)
	familyHandler := handlers.NewFamilyHandler(familyService, tokenService)
	listHandler := handlers.NewListHandler(listService, tokenService)
	itemHandler := handlers.NewItemHandler(itemService, tokenService)
	appLinkHandler := handlers.NewAppLinkHandler(cfg.Android, cfg.DeepLink)

	// Регистрация маршрутов
	api := r.Group("/api")

	authorize := api.Group("/authorize")
	{
		authorize.GET("/check-user", authHandler.CheckUser)
		authorize.GET("/refresh", authHandler.RefreshTokens)
		authorize.POST("/logout", authHandler.Logout)

		authorize.POST("/send-otp", otpHandler.SendOTP)
		authorize.POST("/verify-otp", otpHandler.VerifyOTP)
	}

	user := api.Group("/user")
	{
		user.PATCH("/update", userHandler.UpdateUserInfo)
		user.POST("/save-image", userHandler.SaveImage)
		user.GET("/me", userHandler.GetMyInfo)
	}

	families := api.Group("/families")
	{
		families.GET("", familyHandler.GetFamilies)
		families.GET("/:family_id/lists", familyHandler.GetLists)
		families.GET("/link-code", familyHandler.GetFamilyLink)

		families.POST("/create-family", familyHandler.CreateFamily)
		families.POST("/join/:link_code", familyHandler.JoinFamily)
	}

	lists := api.Group("/lists")
	{
		lists.GET("/:list_id", listHandler.GetListInfo)
		lists.POST("/create-list", listHandler.CreateList)
		lists.POST("/:list_id/items", listHandler.AddItem)
	}

	items := api.Group("/items")
	{
		items.PUT("/:item_id", itemHandler.ChangeItem)
		items.DELETE("/:item_id", itemHandler.DeleteItem)
		items.PATCH("/:item_id/bought", itemHandler.MarkAsBought)
		items.POST("/:item_id/notify", itemHandler.NotifyMembers)
	}

	// Маршруты для реализации внутренних ссылок для приложения
	r.GET("/.well-known/assetlinks.json", appLinkHandler.AssetLinks)

	// Маршруты для переадресации на скачивание приложения
	families.GET("/join/:link_code", appLinkHandler.OpenInvite)

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
