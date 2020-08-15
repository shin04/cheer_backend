package main

import (
	"app/models"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

var Db *sql.DB // データベースに接続するためのハンドラ

type jwtCustomClaims struct {
	Email string
	jwt.StandardClaims
}

var signingKey = []byte("secret")
var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

func main() {
	dbInit()

	// routing
	e := echo.New()

	e.Use(middleware.CORS())

	e.POST("/signup", Signup)
	e.POST("/login", Login)

	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(Config)) // apiグループは認証が必要

	e.GET("/posts", getPostList)
	e.GET("/posts/:id", getPost)
	e.POST("/posts", createPost)
	e.DELETE("/posts/:id", deletePost)
	e.PUT("/posts", updatePost)
	e.POST("/posts/cheerup/:id", cheerup)

	e.Logger.Fatal(e.Start(":8080"))
}

func dbInit() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	// var err error
	connect_setting := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " sslmode=disable"
	Db, err = sql.Open("postgres", connect_setting)
	if err != nil {
		panic(err)
	}
}

func Signup(c echo.Context) error {
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		return err
	}

	if user.Username == "" || user.Password == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid name or password",
		}
	}

	if err := user.CreateUser(Db); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}

func Login(c echo.Context) error {
	var json map[string]interface{} = map[string]interface{}{}
	if err := c.Bind(&json); err != nil {
		return err
	}

	token, err := models.Login(json["email"].(string), json["password"].(string), Db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func getPostList(c echo.Context) error {
	posts, err := models.GetPostList(Db)
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, posts)
}

func getPost(c echo.Context) error {
	id := c.Param("id")
	post, err := models.GetPost(id, Db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, post)
}

func createPost(c echo.Context) error {
	post := &models.Post{}
	err := c.Bind(post)
	if err != nil {
		return err
	}

	err = post.CreatePost(Db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)
}

func deletePost(c echo.Context) error {
	id := c.Param("id")
	err := models.DeletePost(id, Db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func updatePost(c echo.Context) error {
	post := &models.Post{}
	err := c.Bind(post)
	if err != nil {
		return err
	}
	err = post.UpdatePost(Db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, post)
}

func cheerup(c echo.Context) error {
	id := c.Param("id")
	post, err := models.FindPost(id, Db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	post.Cheer = post.Cheer + 1
	err = post.UpdatePost(Db)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, post)
}
