package models

type Comments struct {
	Id           string
	Post         Post
	Content      string
	Created_date string
}
