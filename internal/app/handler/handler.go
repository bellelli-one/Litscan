package handler

import (
	"RIP/internal/app/config"
	"RIP/internal/app/redis"
	"RIP/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const hardcodedUserID = 2

type Handler struct {
	Repository *repository.Repository
	Redis      *redis.Client
	JWTConfig  *config.JWTConfig
}

func NewHandler(r *repository.Repository, redis *redis.Client, jwtConfig *config.JWTConfig) *Handler {
	return &Handler{
		Repository: r,
		Redis:      redis,
		JWTConfig:  jwtConfig,
	}
}

func (h *Handler) RegisterAPI(r *gin.RouterGroup) {

	// Доступны всем
	r.POST("/users", h.Register)
	r.POST("/auth/login", h.Login)
	r.GET("/books", h.GetBooks)
	r.GET("/books/:id", h.GetBook)

	// Эндпоинты, доступные только авторизованным пользователям
	auth := r.Group("/")
	auth.Use(h.AuthMiddleware)
	{
		// Пользователи
		auth.POST("/auth/logout", h.Logout)
		auth.GET("/users/:id", h.GetUserData)
		auth.PUT("/users/:id", h.UpdateUserData)

		// Заявки AnalyseBooks
		auth.POST("/analyse-books/draft/books/:book_id", h.AddBookToDraft)
		auth.GET("/analyse-books/cart", h.GetCartBadge)
		auth.GET("/analyse-books", h.ListAnalyseBooks)
		auth.GET("/analyse-books/:id", h.GetAnalyseBooks)
		auth.PUT("/analyse-books/:id", h.UpdateAnalyseBooks)
		auth.PUT("/analyse-books/:id/form", h.FormAnalyseBooks)
		auth.DELETE("/analyse-books/:id", h.DeleteAnalyseBooks)
		auth.DELETE("/analyse-books/:id/books/:book_id", h.RemoveBookFromAnalyseBooks)
		auth.PUT("/analyse-books/:id/books/:book_id", h.UpdateBookToApplication)
	}

	// Эндпоинты, доступные только модераторам
	moderator := r.Group("/")
	moderator.Use(h.AuthMiddleware, h.ModeratorMiddleware)
	{
		// Управление книгами
		moderator.POST("/books", h.CreateBook)
		moderator.PUT("/books/:id", h.UpdateBook)
		moderator.DELETE("/books/:id", h.DeleteBook)
		moderator.POST("/books/:id/image", h.UploadBookImage)

		// Управление заявками AnalyseBooks
		moderator.PUT("/analyse-books/:id/resolve", h.ResolveAnalyseBooks)
	}
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
