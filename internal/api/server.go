package api

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "log"
  "RIP/internal/app/handler"
  "RIP/internal/app/repository"
)

func StartServer() {
  log.Println("Starting server")

  repo, err := repository.NewRepository()
  if err != nil {
    logrus.Error("ошибка инициализации репозитория")
  }

  handler := handler.NewHandler(repo)

  r := gin.Default()
  // добавляем наш html/шаблон
  r.LoadHTMLGlob("templates/*")
  r.Static("/static", "./resources")

  r.GET("/litscan", handler.GetBooks)
  r.GET("/book/:id", handler.GetBook)
  r.GET("/order/:id", handler.GetOrder)

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
  log.Println("Server down")
}