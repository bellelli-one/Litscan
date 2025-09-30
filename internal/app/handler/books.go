package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAllBooks(ctx *gin.Context) {
	var books []ds.Books
	var err error

	searchingBooks := ctx.Query("searchingBooks")
	if searchingBooks == "" {
		books, err = h.Repository.GetAllBooks()
	} else {
		books, err = h.Repository.SearchBooksByName(searchingBooks)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	draftAppl, err := h.Repository.GetDraftAppl(2)
	var applID uint = 0
	var applCount int = 0

	if err == nil && draftAppl != nil {
		fullAppl, err := h.Repository.GetApplWithBooks(draftAppl.ID)
		if err == nil {
			applID = fullAppl.ID
			applCount = len(fullAppl.BooksLink)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"books":       books,
		"booksSearch": searchingBooks,
		"applID":      applID,
		"applCount":   applCount,
	})
}

func (h *Handler) GetBookByID(ctx *gin.Context) {
	strId := ctx.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	book, err := h.Repository.GetBookByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "book.html", book)
}
