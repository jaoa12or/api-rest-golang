package main

import (
	"challenge-backend/handlers"
	"database/sql"
	"log"
	"github.com/buaazp/fasthttprouter"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
)

var (
	corsAllowCredentials = "true"
	corsAllowHeaders     = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	corsAllowMethods     = "HEAD, POST, GET, OPTIONS, PUT, DELETE"
	corsAllowOrigin      = "http://localhost:8084"
)

// CORS :
func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		next(ctx)
	}
}

func main() {
	// Connect to the "challenge" database.
	db, err := sql.Open("postgres", "postgresql://challenge@localhost:26257/challenge?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	handlers.MigrateDB(db)
	router := fasthttprouter.New()
	router.GET("/domains", handlers.GetDomains(db))
	router.POST("/domains/search", handlers.ConsultDomain(db))
	if err := fasthttp.ListenAndServe(":9000", CORS(router.Handler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
