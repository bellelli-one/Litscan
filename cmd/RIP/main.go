package main

import (
	"RIP/internal/app/config"
	"RIP/internal/app/dsn"
	"RIP/internal/app/handler"
	"RIP/internal/app/redis"
	"RIP/internal/app/repository"
	"RIP/internal/pkg"
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title           API для системы анализа книг
// @version         1.0
// @description     API-сервер для управления книгами, заявками и анализом текстов
// @contact.name    API Support
// @contact.email   support@example.com
// @host            localhost:8090
// @BasePath        /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// main godoc
// @Summary Запуск приложения
// @Description Основная функция инициализации и запуска API сервера
// @Tags system
func main() {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		// Вместо AllowOrigins используем функцию для проверки
		AllowOriginFunc: func(origin string) bool {
			// Список разрешенных адресов
			return origin == "http://localhost:3000" ||
				origin == "tauri://localhost" || // Для macOS/Linux
				origin == "https://tauri.localhost" || // Для Windows
				origin == "http://192.168.1.151:3000" // Для тестов с телефона/сети
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Загрузка конфигурации
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	// Подключение к PostgreSQL
	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	// Инициализация репозитория
	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	// Инициализация Redis
	redisClient, errRedis := redis.New(context.Background(), conf.Redis)
	if errRedis != nil {
		logrus.Fatalf("error initializing redis: %v", errRedis)
	}

	// Инициализация хендлеров
	hand := handler.NewHandler(rep, redisClient, &conf.JWT)

	// Запуск приложения
	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
