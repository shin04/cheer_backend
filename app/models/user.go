package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type jwtCustomClaims struct {
	Email string
	jwt.StandardClaims
}

var signingKey = []byte("secret")
var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

type User struct {
	Id           string
	Username     string
	Email        string
	Password     string
	Created_date string
}

func (user *User) CreateUser(db *sql.DB) (err error) {
	query := "insert into users (username, email, password) values ($1, $2, $3) returning id, created_date"
	statement, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer statement.Close()

	err = statement.QueryRow(user.Username, user.Email, user.Password).Scan(&user.Id, &user.Created_date)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func FindUser(email string, db *sql.DB) (user User, err error) {
	query := "select * from users where email = $1"
	err = db.QueryRow(query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.Created_date)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func Login(email string, password string, db *sql.DB) (token string, err error) {
	user, err := FindUser(email, db)
	if err != nil {
		fmt.Println(err)
		return
	}
	if user.Password != password {
		fmt.Println("invalid password")
		return
	}

	// set claims
	claims := &jwtCustomClaims{
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// create token
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString(signingKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
