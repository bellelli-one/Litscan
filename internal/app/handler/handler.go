package handler

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "RIP/internal/app/repository"
  "net/http"
  "strconv"
)

type Handler struct {
  Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
  return &Handler{
    Repository: r,
  }
}

func (h *Handler) GetBooks(ctx *gin.Context) {
	var books []repository.Books
	var err error
	searchQuery := ctx.Query("query")
	// addItem := ctx.Query("addItem")
	// if addItem != ""{
	// 	itemID, _ := strconv.Atoi(addItem)
    //     h.Repository.AddToCart(itemID)  
	// }
	if searchQuery == "" {         
		books, err = h.Repository.GetBooks()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		books, err = h.Repository.GetBooksByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}
	items, _ := h.Repository.GetBooksInOrder(1)
	count := len(items.Books)
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"books": books,
		"query":  searchQuery,
		"count": count,
	})
}

func (h *Handler) GetBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	book, err := h.Repository.GetBook(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "book.html", gin.H{
		"book": book,
	})
}

func (h *Handler) GetOrder(ctx *gin.Context) {
	// if h.Repository.GetCartCount() == 0 {
    //     // Перенаправляем обратно с сообщением
    //     ctx.Redirect(http.StatusFound, "/hello?error=empty_cart")
    //     return
    // }

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetBooksInOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	BooksInApplication := order.Books
	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"Books": BooksInApplication,
		"AvgWordLen": order.AvgWordLen,
		"LexicalDiversity": order.LexicalDiversity,
		"ConjunctionFreq": order.ConjunctionFreq,
		"AvgSentenceLen": order.AvgSentenceLen,
		"Result": order.Result,
	})
}

// func (h *Handler) GetOrders(ctx *gin.Context) {
// 	order, err := h.Repository.GetOrders()
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	ctx.HTML(http.StatusOK, "index.html", gin.H{
// 		"order": order,
// 	})
// }

// func (h *Handler) GetApplicationComponents(ctx *gin.Context) {
// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	total := 0
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	items, _ := h.Repository.GetApplicationComponents(id)
// 	for _, i := range items {
// 		total += i.Price
// 	}

// 	ctx.HTML(http.StatusOK, "cart.html", gin.H{
// 		"items":     items,
// 		"total":     total,
// 		"cartCount": len(items),
// 	})
// }
