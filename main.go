package main

import (
	"log"
	"net/http"
	"os"

	handler "github.com/UMichael/goweb/handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var person handler.User
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("HTTP_PLATFORM_PORT")
	}
	if port == "" {
		port = "3000"
	}
	//fileServer := http.FileServer(http.Dir("./template/images"))
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.RedirectTrailingSlash = true
	router.RedirectFixedPath = true
	router.GET("/confirm/:token", person.ConfirmToken)
	router.GET("/signup", person.SignUp)
	router.POST("/signup/", person.SignUpPost)
	router.GET("/login", person.Login)
	router.POST("/login/", person.LoginPost)
	router.GET("/success", person.Success)
	router.GET("/", person.Index)
	router.GET("/logout/", person.Logout)
	router.ServeFiles("/assets/*filepath", http.Dir("./template/assets"))
	log.Fatal(http.ListenAndServe(":"+port, router))
}
