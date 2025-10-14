package handler

import (
	"RIP/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GET /api/analysebooks/cart - иконка корзины
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
func (h *Handler) ListAnalyseBooks(c *gin.Context) {
	status := c.Query("status")
	from := c.Query("from")
	to := c.Query("to")

	appList, err := h.Repository.AnalyseBooksListFiltered(status, from, to)
	if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, appList)
}

// GET /api/analysebooks/:id - одна заявка с книгами
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

// const hardcodedUserID = 2

// func (h *Handler) AddBookToAppl(c *gin.Context) {
// 	bookID, err := strconv.Atoi(c.Param("book_id"))
// 	if err != nil {
// 		h.errorHandler(c, http.StatusBadRequest, err)
// 		return
// 	}

// 	appl, err := h.Repository.GetDraftAppl(hardcodedUserID)
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		newAppl := ds.Application{
// 			CreatorID: hardcodedUserID,
// 			Status:    ds.StatusDraft,
// 		}
// 		if createErr := h.Repository.CreateAppl(&newAppl); createErr != nil {
// 			h.errorHandler(c, http.StatusInternalServerError, createErr)
// 			return
// 		}
// 		appl = &newAppl
// 	} else if err != nil {
// 		h.errorHandler(c, http.StatusInternalServerError, err)
// 		return
// 	}

// 	if err = h.Repository.AddBookToAppl(appl.ID, uint(bookID)); err != nil {
// 	}

// 	c.Redirect(http.StatusFound, "/litscan")
// }

// func (h *Handler) GetAppl(c *gin.Context) {
// 	applID, err := strconv.Atoi(c.Param("appl_id"))
// 	if err != nil {
// 		h.errorHandler(c, http.StatusBadRequest, err)
// 		return
// 	}

// 	appl, err := h.Repository.GetApplWithBooks(uint(applID))
// 	if err != nil {
// 		h.errorHandler(c, http.StatusNotFound, err)
// 		return
// 	}

// 	if len(appl.BooksLink) == 0 {
// 		h.errorHandler(c, http.StatusForbidden, errors.New("cannot access an empty appl page, add books first"))
// 		return
// 	}

// 	c.HTML(http.StatusOK, "order.html", appl)
// }

// func (h *Handler) DeleteAppl(c *gin.Context) {
// 	applID, _ := strconv.Atoi(c.Param("appl_id"))

// 	if err := h.Repository.LogicallyDeleteAppl(uint(applID)); err != nil {
// 		h.errorHandler(c, http.StatusInternalServerError, err)
// 		return
// 	}

// 	c.Redirect(http.StatusFound, "/litscan")
// }
