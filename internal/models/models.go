package models

import "time"

type SimpleRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Film struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Duration      int         `json:"duration"`
	YearOfIssue   int         `json:"year_of_issue"`
	Age           int         `json:"age"`
	LinkImg       any         `json:"link_img"`
	LinkKinopoisk any         `json:"link_kinopoisk"`
	LinkVideo     string      `json:"link_video"`
	CreatedAt     string      `json:"created_at"`
	Country       SimpleRef   `json:"country"`
	Categories    []SimpleRef `json:"categories"`
	RatingAvg     float64     `json:"ratingAvg"`
	ReviewCount   int         `json:"reviewCount"`
}

type FilmFilters struct {
	Search   string
	SortBy   string
	SortDir  string
	Category int
	Country  int
	Limit    int
	Offset   int
}

type Category struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	ParentCategory *SimpleRef `json:"parentCategory"`
	FilmCount      int        `json:"filmCount"`
}

type Country struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	FilmCount int    `json:"filmCount"`
}

type Gender struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID          int       `json:"id"`
	FIO         string    `json:"fio"`
	Email       string    `json:"email"`
	Birthday    time.Time `json:"-"`
	GenderID    int       `json:"-"`
	GenderName  string    `json:"-"`
	Password    string    `json:"-"`
	ReviewCount int       `json:"reviewCount"`
	RatingCount int       `json:"ratingCount"`
}

type Review struct {
	ID         int       `json:"id"`
	FilmID     int       `json:"-"`
	FilmName   string    `json:"-"`
	UserID     int       `json:"-"`
	UserFIO    string    `json:"-"`
	Message    string    `json:"message"`
	IsApproved bool      `json:"-"`
	CreatedAt  time.Time `json:"-"`
}

type Rating struct {
	ID        int       `json:"id"`
	FilmID    int       `json:"-"`
	FilmName  string    `json:"-"`
	Score     int       `json:"score"`
	CreatedAt time.Time `json:"-"`
}
