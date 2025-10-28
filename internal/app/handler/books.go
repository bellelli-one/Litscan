package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/books - список книг с фильтрацией

// GetBooks godoc
// @Summary      Получить список книг (все)
// @Description  Возвращает постраничный список книг.
// @Tags         books
// @Produce      json
// @Param        title query string false "Фильтр по названию книги"
// @Success      200 {object} ds.PaginatedResponse
// @Failure      500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router       /books [get]
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

// GetBook godoc
// @Summary      Получить одну книгу по ID (все)
// @Description  Возвращает детальную информацию о книге.
// @Tags         books
// @Produce      json
// @Param        id path int true "ID книги"
// @Success      200 {object} ds.BookDTO
// @Failure      404 {object} map[string]string "Книга не найдена"
// @Router       /books/{id} [get]
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

// CreateBook godoc
// @Summary      Создать новую книгу (только модератор)
// @Description  Создает новую запись о книге.
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        bookData body ds.BookCreateRequest true "Данные новой книги"
// @Success      201 {object} ds.BookDTO
// @Failure      400 {object} map[string]string "Ошибка валидации"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен (не модератор)"
// @Router       /books [post]
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

// UpdateBook godoc
// @Summary      Обновить книгу (только модератор)
// @Description  Обновляет информацию о существующей книге.
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "ID книги"
// @Param        updateData body ds.BookUpdateRequest true "Данные для обновления"
// @Success      200 {object} ds.BookDTO
// @Failure      400 {object} map[string]string "Ошибка валидации"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен"
// @Router       /books/{id} [put]
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

// DeleteBook godoc
// @Summary      Удалить книгу (только модератор)
// @Description  Удаляет книгу из системы.
// @Tags         books
// @Security     ApiKeyAuth
// @Param        id path int true "ID книги для удаления"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен"
// @Router       /books/{id} [delete]
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

// AddBookToDraft godoc
// @Summary      Добавить книгу в черновик заявки (все)
// @Description  Находит или создает черновик заявки для текущего пользователя и добавляет в него книгу.
// @Tags         books
// @Security     ApiKeyAuth
// @Param        book_id path int true "ID книги для добавления"
// @Success      201 {object} map[string]string "Сообщение об успехе"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router       /analyse-books/draft/books/{book_id} [post]
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

// UploadBookImage godoc
// @Summary      Загрузить изображение для книги (только модератор)
// @Description  Загружает и привязывает изображение к книге.
// @Tags         books
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "ID книги"
// @Param        file formData file true "Файл изображения"
// @Success      200 {object} map[string]string "URL загруженного изображения"
// @Failure      400 {object} map[string]string "Файл не предоставлен"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен"
// @Router       /books/{id}/image [post]
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
