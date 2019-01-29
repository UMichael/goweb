module github.com/UMichael/goweb

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gorilla/websocket v1.4.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/lib/pq v1.0.0
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20190129210102-ccddf3741a0c
//github.com/UMichael/goweb/handlers v0.0.0
)

replace github.com/UMichael/goweb/handlers v0.0.0 => ./handlers
