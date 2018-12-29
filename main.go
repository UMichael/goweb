package main

import (
	"net/http"

	login "./Handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var person login.User
	// httprouter.CleanPath("/")
	router := httprouter.New()
	// router.HandleMethodNotAllowed = true
	// router.RedirectTrailingSlash = true
	// router.RedirectFixedPath = true
	//default page motive is to direct user to forum only to unverified and chatroom on login
	//router.GET("/", index)

	router.GET("/confirm/:token", person.ConfirmToken)
	router.GET("/signup", person.SignUp)
	router.POST("/signup/", person.SignUpPost)
	router.GET("/login", person.Login)
	router.POST("/login/", person.LoginPost)
	router.GET("/", person.Index)
	http.ListenAndServe(":8080", router)
}
