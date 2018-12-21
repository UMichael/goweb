package Handlers

import (
	"fmt"
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

type User struct {
	Email      string `json:"email"`
	Nickname   string `db:"user" json:"nickname"`
	Password   string
	Age        int    `json:"age"`
	Department string `db:"dept" json:"department"`
	Super      bool   `json:"super"`
	Moderator  bool   `db:"mod" json:"mod"`
	Confirmed  bool
	Token      string
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
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
func (user *User) Login(w http.ResponseWriter, r *http.Request) { //check how safe using pointer is to this
	//t, err := template.ParseFiles("./html/login.htm")
	r.ParseForm()
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("pass")
	var hashpass string
	if err = Db.QueryRow("select password from users where email = $1", user.Email).Scan(&hashpass); err != nil {
		//Tell the user there was no email like that found
		//do something
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(user.Password)); err != nil {
		//Tell the user that the password is invalid
		return
	}
	//User login success
	Db.QueryRow("select nickname, age, dept, super, mod, token, created_at from users where email = $1", user.Email).Scan(&user.Nickname, &user.Age, &user.Department, &user.Super, &user.Moderator, &user.Token, &user.CreatedAt)
	//Find a way to use this
	fmt.Println("success logging in")
}

func (user *User) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//implement email authentication
	//t,err:=template.ParseFiles("./html/signup.html")
	r.ParseForm()
	user.Nickname = r.FormValue("user")
	user.Email = r.FormValue("Email")
	user.Password = r.FormValue("pass")
	user.Age, _ = strconv.Atoi(r.FormValue("age"))
	user.Department = r.FormValue("dept")
	fmt.Println(user)
	hashpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error with request ", 500)
	}
	//try to send email to user to confirm and tell user to login and give restrictions to users not confirmed
	statement := `
		insert into users (email, password, department, super, mod, age, confirmed, nickname)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		`
	if _, err = Db.Exec(statement, user.Email, hashpass, user.Department, user.Super, user.Moderator, user.Age, user.Confirmed, user.Nickname); err != nil {
		//if it cant register
		fmt.Println(err)
	}
	fmt.Println(user)
}
func (user *User) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	//get token and send to user
	user.Email = r.FormValue("email")
	Db.QueryRow("select token from users where nickname = $1", user.Nickname)

	return nil
}
func (user *User) Update() error {
	//stop
	Db.QueryRow("insert into user ()")
	return nil
}
func (user *User) Confirm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := ps.ByName("token")
	err := Db.QueryRow("select email, nickname, from users where token = $1", token).Scan(&user.Email, user.Nickname)
	if err != nil {
		fmt.Fprintln(w, "this is a wrong token")
		return
	}
	_, err = Db.Exec("update users set confirmed = $1 where token = $2", true, token)
	if err != nil {
		fmt.Println(w, "Error confirming user")
	}
}

//Test
