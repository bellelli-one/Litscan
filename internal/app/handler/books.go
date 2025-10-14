package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/sirupsen/logrus"
)

// GET /api/books - список книг с фильтрацией
func (h *Handler) GetBooks(c *gin.Context) {
	title := c.Query("title")

	books, total, err := h.Repository.BooksList(title)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	var bookDTOs []ds.BookDTO
	for _, b := range books {
		bookDTOs = append(bookDTOs, ds.BookDTO{
			ID:               b.ID,
			Title:            b.Title,
			Text:             b.Text,
			Image:            b.Image,
			AvgWordLen:       b.AvgWordLen,
			LexicalDiversity: b.LexicalDiversity,
			ConjunctionFreq:  b.ConjunctionFreq,
			AvgSentenceLen:   b.AvgSentenceLen,
			Status:           b.Status,
		})
	}

	c.JSON(http.StatusOK, ds.PaginatedResponse{
		Items: bookDTOs,
		Total: total,
	})
}

// GET /api/books/:id - одна книга
func (h *Handler) GetBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	book, err := h.Repository.GetBookByID(id)
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	bookDTO := ds.BookDTO{
		ID:               book.ID,
		Title:            book.Title,
		Text:             book.Text,
		Image:            book.Image,
		AvgWordLen:       book.AvgWordLen,
		LexicalDiversity: book.LexicalDiversity,
		ConjunctionFreq:  book.ConjunctionFreq,
		AvgSentenceLen:   book.AvgSentenceLen,
		Status:           book.Status,
	}

	c.JSON(http.StatusOK, bookDTO)
}

// POST /api/books - создание книги
func (h *Handler) CreateBook(c *gin.Context) {
	var req ds.BookCreateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	statusValue := false

	book := ds.Books{
		Title:            req.Title,
		Text:             req.Text,
		Image:            req.Image,
		AvgWordLen:       req.AvgWordLen,
		LexicalDiversity: req.LexicalDiversity,
		ConjunctionFreq:  req.ConjunctionFreq,
		AvgSentenceLen:   req.AvgSentenceLen,
		Status:           &statusValue,
	}

	if err := h.Repository.CreateBook(&book); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	bookDTO := ds.BookDTO{
		ID:               book.ID,
		Title:            book.Title,
		Text:             book.Text,
		Image:            book.Image,
		AvgWordLen:       book.AvgWordLen,
		LexicalDiversity: book.LexicalDiversity,
		ConjunctionFreq:  book.ConjunctionFreq,
		AvgSentenceLen:   book.AvgSentenceLen,
		Status:           book.Status,
	}

	c.JSON(http.StatusCreated, bookDTO)
}

// PUT /api/books/:id - обновление книги
func (h *Handler) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	var req ds.BookUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	book, err := h.Repository.UpdateBook(uint(id), req)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	bookDTO := ds.BookDTO{
		ID:               book.ID,
		Title:            book.Title,
		Text:             book.Text,
		Image:            book.Image,
		AvgWordLen:       book.AvgWordLen,
		LexicalDiversity: book.LexicalDiversity,
		ConjunctionFreq:  book.ConjunctionFreq,
		AvgSentenceLen:   book.AvgSentenceLen,
		Status:           book.Status,
	}

	c.JSON(http.StatusOK, bookDTO)
}

// DELETE /api/books/:id - удаление книги
func (h *Handler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.DeleteBook(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Книга удалена",
	})
}

// POST /api/analyse-books/draft/books/:book_id - добавление книги в черновик
func (h *Handler) AddBookToDraft(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.AddBookToDraft(hardcodedUserID, uint(bookID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Черновик создан. Книга добавлена в черновик.",
	})
}

// POST /api/books/:id/image - загрузка изображения книги
func (h *Handler) UploadBookImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	imageURL, err := h.Repository.UploadBookImage(uint(id), file)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"image": imageURL})
}
