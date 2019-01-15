package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var Db *sqlx.DB
var err error
var templates = template.Must(template.ParseGlob("template/*.html"))

//User ...
type User struct {
	Name      string `json:"names"`
	Email     string `json:"email"`
	Nickname  string `db:"user" json:"nickname"` //To implement this later
	Password  string
	Age       int    `json:"age"`
	Faculty   string `db:"dept" json:"department"`
	Super     bool   `json:"super"`
	Moderator bool   `db:"mod" json:"mod"`
	Confirmed bool
	Token     string
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func init() {

	Db, err = sqlx.Open("postgres", "user=DSC password=DSC sslmode=disable dbname=database port=5434")
	if err != nil {
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		panic(err)
	}
}

func executetemplate(file string, w http.ResponseWriter, r *http.Request, temp string) {
	templates = templates.Lookup(file + ".html")
	_, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		templates.Execute(w, temp)
	}

}

//LoginPost ...
func (user *User) LoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	fmt.Println(user)
	var hashpass string
	if err = Db.QueryRow("select password,names from users where email = $1", user.Email).Scan(&hashpass, &user.Name); err != nil {
		//Tell the user there was no email like that found
		//do something
		fmt.Fprintln(w, "error this has not been registered") //fix something explicit
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(user.Password)); err != nil {
		//Tell the user that the password is invalid
		fmt.Fprintln(w, "error this has a wrong pass") //fix something explicit
		return
	}

	//User login success
	//What to do after successfull login
	//Db.QueryRow("select nickname, age, dept, super, mod, token, created_at from users where email = $1", user.Email).Scan(&user.Nickname, &user.Age, &user.Faculty, &user.Super, &user.Moderator, &user.Token, &user.CreatedAt)
	//Find a way to use this
	Create(user, w, r)
	http.Redirect(w, r, "/", 302)
}

//Login ...
func (user *User) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("token")
	templates = templates.Lookup("login.html")
	if err == nil {
		http.Redirect(w, r, "/", 302)
	} else {
		templates.Execute(w, nil)
	}
}

//SignUpPost ...
func (user *User) SignUpPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	user.Name = r.FormValue("name")
	user.Nickname = r.FormValue("Nickname") //Ignored Nickname till fix
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("Password")
	user.Age, _ = strconv.Atoi(r.FormValue("Age"))
	user.Faculty = r.FormValue("Faculty")
	hashpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error with request ", 500)
	}
	//try to send email to user to confirm and tell user to login and give restrictions to users not confirmed
	statement := `
		insert into users (nickname, email, password, faculty, super, mod, age, confirmed, names)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
	if _, err = Db.Exec(statement, user.Nickname, user.Email, hashpass, user.Faculty, user.Super, user.Moderator, user.Age, user.Confirmed, user.Name); err != nil {
		//if it cant register
		fmt.Fprintln(w, "This email has been registered already please use another", err)
		return
	}
	//after success take it to successful page
	Create(user, w, r)
}

//SignUp ...
func (user *User) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	templates = templates.Lookup("signup.html")
	_, err := r.Cookie("token")
	if err == nil {
		http.Redirect(w, r, "/", 302)
	} else {
		templates.Execute(w, nil)
	}
}

//ResetPassword ...
//Work in progress
// func (user *User) ResetPassword(w http.ResponseWriter, r *http.Request) error {
// 	//get token and send to user
// 	user.Email = r.FormValue("email")
// 	Db.QueryRow("select token from users where nickname = $1", user.Nickname)

// 	return nil
// }

//To Update Users info
//Update ...
// func (user *User) Update() error {
// 	//stop
// 	Db.QueryRow("insert into user ()")
// 	return nil
// }

//ConfirmToken ...
func (user *User) ConfirmToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := ps.ByName("token")
	err := Db.QueryRow("select email, names from users where token = $1", token).Scan(&user.Email, &user.Name)
	if err != nil {
		fmt.Fprintln(w, "this is a wrong token", err)
		return
	}
	_, err = Db.Exec("update users set confirmed = $1 where token = $2", true, token)
	if err != nil {
		fmt.Println(w, "Error confirming user")
	}
}

//Index ...
func (user *User) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	executetemplate("success", w, r, "")
}

//Success ....
func (user *User) Success(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	cookie, _ := r.Cookie("token")
	token := cookie.Value
	err = Db.QueryRow("select email, names from users where token = $1", token).Scan(&user.Email, &user.Name)
	executetemplate("success", w, r, user.Name)
	fmt.Println("Username is ", user.Name)
}
func (user *User) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cook := http.Cookie{Name: "token", Value: "", HttpOnly: true, MaxAge: -1, Path: "/"}
	http.SetCookie(w, &cook)
	http.Redirect(w, r, "/login", 302)
}
