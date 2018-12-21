package main

import (
	"net/http"

	login "./Handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var person login.User
	httprouter.CleanPath("/")
	router := httprouter.New()
	//router.POST("/signup", person.SignUp)
	router.HandleMethodNotAllowed = false
	//default page
	//router.GET("/", index)
	router.PATCH("/confirm/:token", person.Confirm)
	router.POST("/signup", person.SignUp)
	http.HandleFunc("/signin", person.Login)
	http.ListenAndServe(":8080", router)
}
