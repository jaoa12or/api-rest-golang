package main

import (
	"challenge-backend/handlers"
	"database/sql"
	"log"
	"github.com/buaazp/fasthttprouter"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
	"challenge-backend/models"
)
// cors policies values
const (
	corsAllowCredentials = "true"
	corsAllowHeaders     = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	corsAllowMethods     = "HEAD, POST, GET, OPTIONS, PUT, DELETE"
	corsAllowOrigin      = "http://localhost:8083"
)

// CORS : the cors policy
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
	models.Migrate(db)
	router := fasthttprouter.New()
	// get collection of domains
	router.GET("/domains", handlers.GetDomains(db))
	// request by domain name the information
	router.POST("/domains/search", handlers.ConsultDomain(db))
	// starting backend information
	log.Println("Running server: Status ok")
	if err := fasthttp.ListenAndServe(":9000", CORS(router.Handler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
	
}

