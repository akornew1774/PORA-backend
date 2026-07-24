// main - пакет с точкой входа приложения
package main

import (
	"context"
	"os"
	"os/signal"
	"pora/internal/config"
	"pora/internal/handlers"
	"pora/internal/infrastructure/database"
	"pora/internal/infrastructure/firebase"
	"pora/internal/infrastructure/logger"
	"pora/internal/infrastructure/storage"
	"pora/internal/middleware"
	"pora/internal/repositories"
	"pora/internal/scheduler"
	"pora/internal/services"
	"pora/internal/websocket"
	"syscall"
	"time"

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

	// Создание контекста для единовременной остановки бесконечных процессов
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Остановка приложения после завершения контекста
	go func() {
		<-ctx.Done()
		stop()
	}()

	// Загрузка конфигураций
	cfg := config.Load()

	// Подключение хранилища
	storage := storage.NewLocalStorage(os.Getenv("BASE_PATH"))

	// Создание клиента для отправки Push-уведомлений через Firebase
	firebaseClient, err := firebase.NewClient(ctx, cfg.FireBase)
	if err != nil {
		logger.Log.Fatal("Ошибка при создании клиента для Firebase: ", err)
	}

	// Создание менеджера соединений для Websocket
	hub := websocket.NewHub()

	// Подключение репозиториев
	userRepo := repositories.NewUserRepository(db)
	otpRepo := repositories.NewOtpRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)
	familyRepo := repositories.NewFamilyRepository(db)
	memberRepo := repositories.NewMemberRepository(db)
	listRepo := repositories.NewListRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)
	briefItemRepo := repositories.NewBriefItemRepository(db)

	// Подключение сервисов
	tokenService := services.NewTokenService(userRepo, refreshTokenRepo)
	fileService := services.NewFileService(storage)
	pushService := services.NewPushService(deviceRepo, firebaseClient, cfg.Push)
	authService := services.NewAuthService(userRepo, refreshTokenRepo, tokenService)
	otpService := services.NewOTPService(otpRepo, userRepo, deviceRepo, tokenService)
	listService := services.NewListService(userRepo, memberRepo, familyRepo, listRepo, itemRepo, hub, cfg.Item)
	userService := services.NewUserService(userRepo, deviceRepo, fileService, listService)
	familyService := services.NewFamilyService(userRepo, memberRepo, familyRepo, listService, cfg.DeepLink)
	itemService := services.NewItemService(familyRepo, listRepo, itemRepo, listService, pushService, cfg.Item)
	statisticsService := services.NewStatisticsService(userRepo, briefItemRepo)

	// Подключение хэндлеров
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	otpHandler := handlers.NewOtpHandler(otpService)
	userHandler := handlers.NewUserHandler(userService, tokenService)
	familyHandler := handlers.NewFamilyHandler(familyService, tokenService)
	listHandler := handlers.NewListHandler(listService, tokenService)
	itemHandler := handlers.NewItemHandler(itemService, tokenService)
	statisticsHandler := handlers.NewStatisticsHandler(statisticsService, tokenService)
	appLinkHandler := handlers.NewAppLinkHandler(cfg.Android, cfg.DeepLink)
	wsHandler := handlers.NewWSHandler(hub, tokenService)

	// Создание объектов Scheduler для выполнений задач по расписанию
	reminderScheduler := scheduler.NewScheduler(time.Hour, itemService.ProcessReminders)
	checkedScheduler := scheduler.NewScheduler(time.Hour, itemService.ProcessCheckedItems)

	// Запуск методов созданных Scheduler
	go reminderScheduler.Run(ctx)
	go checkedScheduler.Run(ctx)

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
		user.PUT("/device", userHandler.UpdateDevice)
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
		families.DELETE("/:family_id", familyHandler.DeleteFamily)
	}

	lists := api.Group("/lists")
	{
		lists.GET("/:list_id", listHandler.GetListInfo)
		lists.POST("/create-list", listHandler.CreateList)
		lists.POST("/:list_id/items", listHandler.AddItem)
		lists.DELETE("/:list_id", listHandler.DeleteList)
	}

	items := api.Group("/items")
	{
		items.GET("/:item_id", itemHandler.GetItemInfo)
		items.PUT("/:item_id", itemHandler.ChangeItem)
		items.DELETE("/:item_id", itemHandler.DeleteItem)

		items.PATCH("/:item_id/bought", itemHandler.MarkAsBought)
		items.POST("/:item_id/notify", itemHandler.NotifyMembers)
	}

	statistics := user.Group("/statistics")
	{
		statistics.GET("/products", statisticsHandler.GetUserProducts)
		statistics.GET("/get_brief", statisticsHandler.GetBrief)
		statistics.POST("/brief", statisticsHandler.SaveBrief)
	}

	// Маршрут для создания Websocket-соединения
	api.GET("/websocket", wsHandler.Connect)

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
