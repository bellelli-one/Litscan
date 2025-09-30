// package handler

// import (
//   "github.com/gin-gonic/gin"
//   "github.com/sirupsen/logrus"
//   "RIP/internal/app/repository"
//   "net/http"
//   "strconv"
//   "fmt"
// )

// type Handler struct {
//   Repository *repository.Repository
// }

// func NewHandler(r *repository.Repository) *Handler {
//   return &Handler{
//     Repository: r,
//   }
// }

// func (h *Handler) GetBooks(ctx *gin.Context) {
// 	var books []repository.Books
// 	var err error
// 	searchQuery := ctx.Query("query")
// 	if searchQuery == "" {         
// 		books, err = h.Repository.GetBooks()
// 		if err != nil {
// 			logrus.Error(err)
// 		}
// 	} else {
// 		books, err = h.Repository.GetBooksByTitle(searchQuery)
// 		if err != nil {
// 			logrus.Error(err)
// 		}
// 	}
// 	items, _ := h.Repository.GetBooksInOrder(1)
// 	count := len(items.Books)
// 	ctx.HTML(http.StatusOK, "index.html", gin.H{
// 		"books": books,
// 		"query":  searchQuery,
// 		"count": count,
// 	})
// }

// func (h *Handler) GetBook(ctx *gin.Context) {
// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	book, err := h.Repository.GetBook(id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	ctx.HTML(http.StatusOK, "book.html", gin.H{
// 		"book": book,
// 	})
// }

// func (h *Handler) GetOrder(ctx *gin.Context) {

// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	order, err := h.Repository.GetBooksInOrder(id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	BooksInArray, err := h.Repository.GetArrayOfBooks(id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	BooksInApplication := order.Books
// 	fmt.Println(BooksInArray)
// 	ctx.HTML(http.StatusOK, "order.html", gin.H{
// 		"Books": BooksInArray,
// 		"BooksInApplication": BooksInApplication,
// 		"AvgWordLen": order.AvgWordLen,
// 		"LexicalDiversity": order.LexicalDiversity,
// 		"ConjunctionFreq": order.ConjunctionFreq,
// 		"AvgSentenceLen": order.AvgSentenceLen,
// 		"Result": order.Result,
// 	})
// }

package handler

import (
	"RIP/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/litscan", h.GetAllBooks)
	router.GET("/book/:id", h.GetBookByID)
	router.GET("/order/:appl_id", h.GetAppl)
	router.POST("/order/add/book/:book_id", h.AddBookToAppl)
	router.POST("/order/:appl_id/delete", h.DeleteAppl)

}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/resources", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}