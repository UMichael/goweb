package handlers

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	keys   = "Hello there DSC Unilag"
	expire = 259200
)

//Create ....
func Create(user *User, w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := make(jwt.MapClaims)
	claims["userid"] = user
	claims["exp"] = time.Now().Add(time.Hour * 72)
	claims["gen"] = "Test" //change this later
	token.Claims = claims

	key, err := token.SignedString([]byte(keys))
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	cookie := http.Cookie{
		Name:     "token",
		Value:    key,
		MaxAge:   expire,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	statement := `
		update users set token = $1 where email = $2`
	_, err = Db.Exec(statement, key, user.Email)
	fmt.Println("token err", err)
}

//Decode ....
func Decode(user *User, r *http.Request, w http.ResponseWriter) error {
	//Get cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("There was no token here")
		http.Redirect(w, r, "/login", 302)
		return err
	}
	tokens := cookie.Value
	var email string
	err = Db.QueryRow("select email from users where token = $1", tokens).Scan(email)
	if email != user.Email {
		w.Write([]byte("A cookie hijack scheme detected"))
		cook := http.Cookie{Name: "user", Value: "", HttpOnly: true, MaxAge: -1, Path: "/"}
		http.SetCookie(w, &cook)
		http.Redirect(w, r, "/login", 302)
		return err
	}
	token, err := jwt.Parse(tokens, func(token *jwt.Token) (interface{}, error) {
		return []byte(keys), nil
	})
	if err == nil && token.Valid {
		return nil
	}
	http.Redirect(w, r, "/login", 302)
	return err
}
