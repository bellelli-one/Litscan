package ds

import "time"

type BookDTO struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	Text             string  `json:"text"`
	Image            *string `json:"image"`
	AvgWordLen       float64 `json:"avg_word_len"`
	LexicalDiversity float64 `json:"lexical_diversity"`
	ConjunctionFreq  float64 `json:"conjunction_freq"`
	AvgSentenceLen   float64 `json:"avg_sentence_len"`
	Status           *bool   `json:"status"`
}

type BookCreateRequest struct {
	Title            string  `json:"title" binding:"required"`
	Text             string  `json:"text" binding:"required"`
	Image            *string `json:"image"`
	AvgWordLen       float64 `json:"avg_word_len"`
	LexicalDiversity float64 `json:"lexical_diversity"`
	ConjunctionFreq  float64 `json:"conjunction_freq"`
	AvgSentenceLen   float64 `json:"avg_sentence_len"`
}

type BookUpdateRequest struct {
	Title            *string  `json:"title"`
	Text             *string  `json:"text"`
	Image            *string  `json:"image"`
	AvgWordLen       *float64 `json:"avg_word_len"`
	LexicalDiversity *float64 `json:"lexical_diversity"`
	ConjunctionFreq  *float64 `json:"conjunction_freq"`
	AvgSentenceLen   *float64 `json:"avg_sentence_len"`
	Status           *bool    `json:"status"`
}

type AnalyseBooksDTO struct {
	ID               uint                   `json:"id"`
	Status           int                    `json:"status"`
	CreationDate     time.Time              `json:"creation_date"`
	CreatorID        uint                   `json:"creator_login"`
	ModeratorID      *uint                  `json:"moderator_login"`
	FormingDate      *time.Time             `json:"forming_date"`
	CompletionDate   *time.Time             `json:"completion_date"`
	AvgWordLen       float64                `json:"avg_word_len"`
	LexicalDiversity float64                `json:"lexical_diversity"`
	ConjunctionFreq  float64                `json:"conjunction_freq"`
	AvgSentenceLen   float64                `json:"avg_sentence_len"`
	CalculatedBooks  int                    `json:"calculated_books"` // Сколько посчитано (Similarity != nil)
	TotalBooks       int                    `json:"total_books"`
	Books            []BookInApplicationDTO `json:"books,omitempty"`
}

type BookInApplicationDTO struct {
	BookID           uint     `json:"book_id"`
	Title            string   `json:"title"`
	Text             string   `json:"text"`
	Image            *string  `json:"image"`
	AvgWordLen       float64  `json:"avg_word_len"`
	LexicalDiversity float64  `json:"lexical_diversity"`
	ConjunctionFreq  float64  `json:"conjunction_freq"`
	AvgSentenceLen   float64  `json:"avg_sentence_len"`
	Description      *string  `json:"description"`
	Similarity       *float64 `json:"similarity"`
}

type BookToApplicationUpdateRequest struct {
	Description *string `json:"description"`
}

// AnalyseBooksUpdateRequest - Исправлено: добавлены поля статуса и дат
type AnalyseBooksUpdateRequest struct {
	Status         *int       `json:"status"`
	FormingDate    *time.Time `json:"forming_date"`
	CompletionDate *time.Time `json:"completion_date"`
	Moderator      *bool      `json:"moderator"`

	// Новые поля метрик
	AvgWordLen       *float64        `json:"avg_word_len"`
	LexicalDiversity *float64        `json:"lexical_diversity"`
	ConjunctionFreq  *float64        `json:"conjunction_freq"`
	AvgSentenceLen   *float64        `json:"avg_sentence_len"`
	BookResults      []BookResultDTO `json:"book_results"`
}

type AnalyseBooksResolveRequest struct {
	Action string `json:"action" binding:"required"` // "complete" | "reject"
}

type BookMetric struct {
	AvgWordLen       float64 `json:"avg_word_len"`
	LexicalDiversity float64 `json:"lexical_diversity"`
	ConjunctionFreq  float64 `json:"conjunction_freq"`
	AvgSentenceLen   float64 `json:"avg_sentence_len"`
}

type AsyncRequest struct {
	ID     uint         `json:"id"`
	Target BookMetric   `json:"target"`
	Books  []BookMetric `json:"books"`
}

type AnalyseBooksBadgeDTO struct {
	ApplicationID *uint `json:"application_id"`
	Count         int   `json:"count"`
}

type UserRegisterRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDTO struct {
	ID        uint   `json:"id"`
	FullName  string `json:"full_name"`
	Username  string `json:"username"`
	Moderator bool   `json:"moderator"`
}

type UserUpdateRequest struct {
	Username *string `json:"username"`
	FullName *string `json:"full_name"`
	Password *string `json:"password"`
}

type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type PaginatedResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

type BookResultDTO struct {
	BookID     uint    `json:"book_id"`
	Similarity float64 `json:"similarity"`
}
