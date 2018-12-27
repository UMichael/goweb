package main

import (
	"net/http"

	login "./Handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var person login.User
	//httprouter.CleanPath("/")
	router := httprouter.New()
	router.ServeFiles("/*filepath", http.Dir("./template/"))
	//router.HandleMethodNotAllowed = false
	//default page
	//router.GET("/", index)
	router.GET("/confirm/:token", person.Confirm)
	router.POST("/signup", person.SignUp)
	router.POST("/signin", person.Login)
	http.ListenAndServe(":8080", router)
}
