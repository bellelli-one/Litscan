package repository

import (
  "fmt"
  "strings"
)



type Repository struct {
}

func NewRepository() (*Repository, error) {
  return &Repository{
  }, nil
}


type Books struct {
  ID int 
  Title string
  ImageUrl string
  Description string
  AvgWordLen       float64
  LexicalDiversity float64
  ConjunctionFreq  float64
  AvgSentenceLen   float64
}

type BooksToApplication struct{
	Book Books
	Description string
}

type Application struct { 
	Books []BooksToApplication
	AvgWordLen       float64
	LexicalDiversity float64
	ConjunctionFreq  float64
	AvgSentenceLen   float64
	Result string
}

var books = []Books{ 
    {
      ID:    1,
      Title: "Капитанская дочка",
	  ImageUrl: "http://172.18.0.4:9000/test/kapitanskayadochka.jpeg",
	  Description: `«Капитанская дочка» — это историческая повесть Александра Пушкина
					о чести, долге и любви, разворачивающаяся на фоне пугачёвского
					восстания и представленная как семейные записки молодого
					офицера Петра Гринёва.`,
	  AvgWordLen: 5.15,
	  LexicalDiversity: 0.28,
	  ConjunctionFreq: 0.075,
	  AvgSentenceLen: 9.89,
    },
    {
      ID:    2,
      Title: "Война и мир",
	  ImageUrl: "http://172.18.0.4:9000/test/voinaimir.jpg",
	  Description: `«Война и мир» — это масштабная эпопея Льва Толстого о судьбах
					русского общества на фоне Наполеоновских войн,исследующая вечные
					вопросы истории, свободы воли, любви и смысла человеческого существования.`,
	  AvgWordLen: 5.08,
	  LexicalDiversity: 0.14,
	  ConjunctionFreq: 0.093,
	  AvgSentenceLen: 12.78,
    },
    {
      ID:    3,
      Title: "Грокаем алгоритмы",
	  ImageUrl: "http://172.18.0.4:9000/test/grokaem.jpg",
	  Description: `«Грокаем алгоритмы» — это иллюстрированное руководство Адитьи Бхаграва,
	  				которое наглядно и доступно объясняет базовые алгоритмы и структуры данных
					через пошаговые примеры на Python и простые графические схемы.`,
	  AvgWordLen: 6.01,
	  LexicalDiversity: 0.19,
	  ConjunctionFreq: 0.064,
	  AvgSentenceLen: 8.41,
    },
	{
      ID:    4,
      Title: "Компьютерные сети",
	  ImageUrl: "http://172.18.0.5:9000/test/tanenb.jpg",
	  Description: `«Компьютерные сети» Эндрю Таненбаума — это книга, в которой последовательно
	                изложены основные концепции, определяющие современное состояние компьютерных
					сетей и тенденции их развития`,
	  AvgWordLen: 6.49,
	  LexicalDiversity: 0.21,
	  ConjunctionFreq: 0.058,
	  AvgSentenceLen: 11.25,
    },
}

var BooksInOrder = map[int]Application{
	1: {
		Books: []BooksToApplication{
			{Book: books[0], Description: "Description 1"},
			{Book: books[1], Description: "Description 2"},
		},
		AvgWordLen: 7.15,
		LexicalDiversity: 0.22,
		ConjunctionFreq: 0.055,
		AvgSentenceLen: 11.89,
		Result: "Результат",
	},
}

func (r *Repository) GetBooks() ([]Books, error) {
  if len(books) == 0 {
    return nil, fmt.Errorf("массив пустой")
  }

  return books, nil
}

func (r *Repository) GetBook(id int) (Books, error) {
	books, err := r.GetBooks()
	if err != nil {
		return Books{}, err 
	}

	for _, book := range books {
		if book.ID == id {
			return book, nil 
		}
	}
	return Books{}, fmt.Errorf("заказ не найден") 
}

func (r *Repository) GetBooksByTitle(title string) ([]Books, error) {
	books, err := r.GetBooks()
	if err != nil {
		return []Books{}, err
	}
	var result []Books
	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(title)) {
			result = append(result, book)
		}
	}
	return result, nil
}

func (r *Repository) GetBooksInOrder(id int) (Application, error) {
	return BooksInOrder[id], nil
}