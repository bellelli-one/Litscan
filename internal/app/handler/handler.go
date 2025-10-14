package handler

import (
	"RIP/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const hardcodedUserID = 2

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterAPI(r *gin.RouterGroup) {
	// Домен книг
	r.GET("/books", h.GetBooks)
	r.GET("/books/:id", h.GetBook)
	r.POST("/books", h.CreateBook)
	r.PUT("/books/:id", h.UpdateBook)
	r.DELETE("/books/:id", h.DeleteBook)
	r.POST("/analyse-books/draft/books/:book_id", h.AddBookToDraft)
	r.POST("/books/:id/image", h.UploadBookImage)

	// Домен заявок (AnalyseBooks)
	r.GET("/analyse-books/cart", h.GetCartBadge)
	r.GET("/analyse-books", h.ListAnalyseBooks)
	r.GET("/analyse-books/:id", h.GetAnalyseBooks)
	r.PUT("/analyse-books/:id", h.UpdateAnalyseBooks)
	r.PUT("/analyse-books/:id/form", h.FormAnalyseBooks)
	r.PUT("/analyse-books/:id/resolve", h.ResolveAnalyseBooks)
	r.DELETE("/analyse-books/:id", h.DeleteAnalyseBooks)

	// Домен м-м (связь заявок и книг)
	r.DELETE("/analyse-books/:id/books/:book_id", h.RemoveBookFromAnalyseBooks)
	r.PUT("/analyse-books/:id/books/:book_id", h.UpdateBookToApplication)

	// Домен пользователь
	r.POST("/users", h.Register)
	r.GET("/users/:id", h.GetUserData)
	r.PUT("/users/:id", h.UpdateUserData)
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
