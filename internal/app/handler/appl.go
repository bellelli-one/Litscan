package handler

import (
	"RIP/internal/app/ds"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const hardcodedUserID = 2

func (h *Handler) AddBookToAppl(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	appl, err := h.Repository.GetDraftAppl(hardcodedUserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newAppl := ds.Application{
			CreatorID: hardcodedUserID,
			Status:    ds.StatusDraft,
		}
		if createErr := h.Repository.CreateAppl(&newAppl); createErr != nil {
			h.errorHandler(c, http.StatusInternalServerError, createErr)
			return
		}
		appl = &newAppl
	} else if err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	if err = h.Repository.AddBookToAppl(appl.ID, uint(bookID)); err != nil {
	}

	c.Redirect(http.StatusFound, "/litscan")
}

func (h *Handler) GetAppl(c *gin.Context) {
	applID, err := strconv.Atoi(c.Param("appl_id"))
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	appl, err := h.Repository.GetApplWithBooks(uint(applID))
	if err != nil {
		h.errorHandler(c, http.StatusNotFound, err)
		return
	}

	if len(appl.BooksLink) == 0 {
		h.errorHandler(c, http.StatusForbidden, errors.New("cannot access an empty appl page, add books first"))
		return
	}

	c.HTML(http.StatusOK, "order.html", appl)
}

func (h *Handler) DeleteAppl(c *gin.Context) {
	applID, _ := strconv.Atoi(c.Param("appl_id"))

	if err := h.Repository.LogicallyDeleteAppl(uint(applID)); err != nil {
		h.errorHandler(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/litscan")
}
