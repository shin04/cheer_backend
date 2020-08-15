package models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Post struct {
	// 小文字から始まると外部パッケージからアクセスできへん
	Id           string
	Title        string
	Content      string
	Cheer        int
	Created_date string
	Update_date  string
	User         User
}

func (post *Post) CreatePost(db *sql.DB) (err error) {
	query := "insert into posts (title, content) values ($1, $2) returning id, created_date, update_date"
	statement, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer statement.Close()

	err = statement.QueryRow(post.Title, post.Content).Scan(&post.Id, &post.Created_date, &post.Update_date)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (post *Post) UpdatePost(db *sql.DB) (err error) {
	update_date := time.Now()
	post.Update_date = update_date.String()
	query := "update posts set title = $2, content = $3, update_date = $4 where id = $1"
	_, err = db.Exec(query, post.Id, post.Title, post.Content, update_date)
	if err != nil {
		return
	}

	return
}

func FindPost(id string, db *sql.DB) (post Post, err error) {
	query := "select * from posts where id = $1"
	err = db.QueryRow(query, id).Scan(&post.Id, &post.Title, &post.Content, &post.Created_date, &post.Update_date, &post.User)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func GetPostList(db *sql.DB) (posts []Post, err error) {
	query := "select * from posts"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.Created_date, &post.Update_date)
		if err != nil {
			fmt.Println(err)
			return
		}
		posts = append(posts, post)
	}
	rows.Close()

	return
}

func GetPost(id string, db *sql.DB) (post Post, err error) {
	query := "select * from posts where id = $1"
	err = db.QueryRow(query, id).Scan(&post.Id, &post.Title, &post.Content, &post.Created_date, &post.Update_date)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func DeletePost(id string, db *sql.DB) (err error) {
	query := "delete from posts where id=$1"
	_, err = db.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

// func CheerUp(id string, db *sql.DB) (post Post, err error) {
// 	post, err = FindPost(id, db)
// 	query := "update posts set cheer $2 where id = $1"
// 	a, err := db.Exec(query, post.Id, post.Cheer+1)
// 	fmt.Println(a)
// 	if err != nil {
// 		return
// 	}

// 	return
// }
