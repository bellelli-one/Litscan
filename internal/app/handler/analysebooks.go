package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /api/analysebooks/cart - иконка корзины

// GetCartBadge godoc
// @Summary      Получить информацию для иконки корзины (авторизованный пользователь)
// @Description  Возвращает ID черновика текущего пользователя и количество книг в нем.
// @Tags         analysebooks
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} ds.AnalyseBooksBadgeDTO
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /analysebooks/cart [get]
func (h *Handler) GetCartBadge(c *gin.Context) {
	draft, err := h.Repository.GetDraftAnalyseBooks(hardcodedUserID)
	if err != nil {
		c.JSON(http.StatusOK, ds.AnalyseBooksBadgeDTO{
			ApplicationID: nil,
			Count:         0,
		})
		return
	}

	fullApp, err := h.Repository.GetAnalyseBooksWithBooks(draft.ID)
	if err != nil {
		logrus.Error("Error getting application with books:", err)
		c.JSON(http.StatusOK, ds.AnalyseBooksBadgeDTO{
			ApplicationID: nil,
			Count:         0,
		})
		return
	}

	c.JSON(http.StatusOK, ds.AnalyseBooksBadgeDTO{
		ApplicationID: &fullApp.ID,
		Count:         len(fullApp.BooksLink),
	})
}

// GET /api/analysebooks - список заявок с фильтрацией

// ListAnalyseBooks godoc
// @Summary      Получить список заявок (авторизованный пользователь)
// @Description  Возвращает отфильтрованный список всех сформированных заявок (кроме черновиков и удаленных).
// @Tags         analysebooks
// @Produce      json
// @Security     ApiKeyAuth
// @Param        status query string false "Фильтр по статусу заявки"
// @Param        from query string false "Фильтр по дате 'от' (формат YYYY-MM-DD)"
// @Param        to query string false "Фильтр по дате 'до' (формат YYYY-MM-DD)"
// @Success      200 {array} ds.AnalyseBooksDTO
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /analysebooks [get]
func (h *Handler) ListAnalyseBooks(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.errorHandler(c, http.StatusUnauthorized, err)
		return
	}
	isModerator := isUserModerator(c)

	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")

	appList, err := h.Repository.AnalyseBooksListFiltered(userID, isModerator, status, from, to)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, appList)
}

// GET /api/analysebooks/:id - одна заявка с книгами

// GetAnalyseBooks godoc
// @Summary      Получить одну заявку по ID (авторизованный пользователь)
// @Description  Возвращает полную информацию о заявке, включая привязанные книги.
// @Tags         analysebooks
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Success      200 {object} ds.AnalyseBooksDTO
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      404 {object} map[string]string "Заявка не найдена"
// @Router       /analysebooks/{id} [get]
func (h *Handler) GetAnalyseBooks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	app, err := h.Repository.GetAnalyseBooksWithBooks(uint(id))
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	var books []ds.BookInApplicationDTO
	for _, link := range app.BooksLink {
		books = append(books, ds.BookInApplicationDTO{
			BookID:           link.BookID,
			Title:            link.Book.Title,
			Text:             link.Book.Text,
			Image:            link.Book.Image,
			AvgWordLen:       link.Book.AvgWordLen,
			LexicalDiversity: link.Book.LexicalDiversity,
			ConjunctionFreq:  link.Book.ConjunctionFreq,
			AvgSentenceLen:   link.Book.AvgSentenceLen,
			Description:      link.Description,
		})
	}

	appDTO := ds.AnalyseBooksDTO{
		ID:               app.ID,
		Status:           app.Status,
		CreationDate:     app.CreationDate,
		CreatorID:        app.Creator.ID,
		ModeratorID:      nil,
		FormingDate:      app.FormingDate,
		CompletionDate:   app.ComplitionDate,
		AwgWordLen:       app.AwgWordLen,
		LexicalDiversity: app.LexicalDiversity,
		ConjunctionFreq:  app.ConjunctionFreq,
		AvgSentenceLen:   app.AvgSentenceLen,
		Books:            books,
	}

	if app.ModeratorID != nil {
		appDTO.ModeratorID = &app.Moderator.ID
	}

	c.JSON(http.StatusOK, appDTO)
}

// PUT /api/analysebooks/:id - изменение полей заявки

// UpdateAnalyseBooks godoc
// @Summary      Обновить данные заявки (авторизованный пользователь)
// @Description  Позволяет пользователю обновить поля своей заявки.
// @Tags         analysebooks
// @Accept       json
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Param        updateData body ds.AnalyseBooksUpdateRequest true "Данные для обновления"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /analysebooks/{id} [put]
func (h *Handler) UpdateAnalyseBooks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	var req ds.AnalyseBooksUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.UpdateAnalyseBooksUserFields(uint(id), req); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Данные заявки обновлены",
	})
}

// PUT /api/analysebooks/:id/form - сформировать заявку

// FormAnalyseBooks godoc
// @Summary      Сформировать заявку (авторизованный пользователь)
// @Description  Переводит заявку из статуса "черновик" в "сформирована".
// @Tags         analysebooks
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки (черновика)"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Не все поля заполнены"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /analysebooks/{id}/form [put]
func (h *Handler) FormAnalyseBooks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.FormAnalyseBooks(uint(id), hardcodedUserID); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Заявка сформирована",
	})
}

// PUT /api/analysebooks/:id/resolve - завершить/отклонить заявку

// ResolveAnalyseBooks godoc
// @Summary      Завершить или отклонить заявку (только модератор)
// @Description  Модератор завершает (с расчетом) или отклоняет заявку.
// @Tags         analysebooks
// @Accept       json
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Param        action body ds.AnalyseBooksResolveRequest true "Действие: 'complete' или 'reject'"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Failure      403 {object} map[string]string "Доступ запрещен"
// @Router       /analysebooks/{id}/resolve [put]
func (h *Handler) ResolveAnalyseBooks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	var req ds.AnalyseBooksResolveRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	moderatorID := uint(hardcodedUserID)
	if err := h.Repository.ResolveAnalyseBooks(uint(id), moderatorID, req.Action); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Заявка обработана модератором",
	})
}

// DELETE /api/analysebooks/:id - удаление заявки

// DeleteAnalyseBooks godoc
// @Summary      Удалить заявку (авторизованный пользователь)
// @Description  Логически удаляет заявку, переводя ее в статус "удалена".
// @Tags         analysebooks
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /analysebooks/{id} [delete]
func (h *Handler) DeleteAnalyseBooks(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.LogicallyDeleteAnalyseBooks(uint(id)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Заявка удалена",
	})
}

// DELETE /api/analysebooks/:id/books/:book_id - удаление книги из заявки

// RemoveBookFromAnalyseBooks godoc
// @Summary      Удалить книгу из заявки (авторизованный пользователь)
// @Description  Удаляет связь между заявкой и книгой.
// @Tags         analysebooks-m-m
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Param        book_id path int true "ID книги"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /analysebooks/{id}/books/{book_id} [delete]
func (h *Handler) RemoveBookFromAnalyseBooks(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	if err := h.Repository.RemoveBookFromAnalyseBooks(uint(appID), uint(bookID)); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Книга удалена из заявки",
	})
}

// PUT /api/analysebooks/:id/books/:book_id - изменение м-м связи

// UpdateBookToApplication godoc
// @Summary      Обновить описание книги в заявке (авторизованный пользователь)
// @Description  Изменяет дополнительное описание для конкретной книги в рамках одной заявки.
// @Tags         analysebooks-m-m
// @Accept       json
// @Security     ApiKeyAuth
// @Param        id path int true "ID заявки"
// @Param        book_id path int true "ID книги"
// @Param        updateData body ds.BookToApplicationUpdateRequest true "Новое описание"
// @Success      204 "No Content"
// @Failure      401 {object} map[string]string "Необходима авторизация"
// @Router       /analysebooks/{id}/books/{book_id} [put]
func (h *Handler) UpdateBookToApplication(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	var req ds.BookToApplicationUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	updateData := ds.BookToAppl{
		Description: req.Description,
	}

	if err := h.Repository.UpdateBookToApplication(uint(appID), uint(bookID), updateData); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "Дополнительная информация к книге обновлена",
	})
}
