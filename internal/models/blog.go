package models

import "time"

type Post struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Slug          string     `json:"slug"`
	Content       string     `json:"content"`
	Excerpt       string     `json:"excerpt"`
	FeaturedImage *Media     `json:"featured_image"`
	Published     bool       `json:"published"`
	PublishedAt   *time.Time `json:"published_at"`
	MetaTitle     string     `json:"meta_title"`
	MetaDesc      string     `json:"meta_description"`
	ReadingTime   int        `json:"reading_time"`
	Author        Author     `json:"author"`
	Tags          []Tag      `json:"tags"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Tag struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Color     string    `json:"color,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Author struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Bio         string            `json:"bio"`
	Avatar      *Media            `json:"avatar"`
	SocialLinks map[string]string `json:"social_links"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type Media struct {
	ID              uint         `json:"id"`
	Name            string       `json:"name"`
	AlternativeText string       `json:"alternativeText"`
	Caption         string       `json:"caption"`
	Width           int          `json:"width"`
	Height          int          `json:"height"`
	Formats         MediaFormats `json:"formats"`
	Hash            string       `json:"hash"`
	Ext             string       `json:"ext"`
	Mime            string       `json:"mime"`
	Size            float64      `json:"size"`
	URL             string       `json:"url"`
	PreviewURL      string       `json:"previewUrl"`
	Provider        string       `json:"provider"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
}

type MediaFormats struct {
	Large     *MediaFormat `json:"large,omitempty"`
	Medium    *MediaFormat `json:"medium,omitempty"`
	Small     *MediaFormat `json:"small,omitempty"`
	Thumbnail *MediaFormat `json:"thumbnail,omitempty"`
}

type MediaFormat struct {
	Name   string  `json:"name"`
	Hash   string  `json:"hash"`
	Ext    string  `json:"ext"`
	Mime   string  `json:"mime"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Size   float64 `json:"size"`
	Path   string  `json:"path"`
	URL    string  `json:"url"`
}
