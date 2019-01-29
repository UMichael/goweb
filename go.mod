module github.com/UMichael/goweb

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/websocket v1.4.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/lib/pq v1.0.0
	golang.org/x/crypto v0.0.0-20190129200851-54b0dbbbf449
	github.com/UMichael/goweb/handlers v0.0.0
)

replace github.com/UMichael/goweb/handlers v0.0.0 => ./handlers
