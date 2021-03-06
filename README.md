# Go Domain Search REST API
A RESTful API for request information about domains

It is a just simple application for making simple RESTful API with Go using **buaazp/fasthttprouter**, **valyala/fasthttp**, **lib/pq** (postgres driver) , **CockroachDB** (Database Engine). 

## Installation & Run Domain Search
```bash
# Download this projects
git clone https://github.com/jaoa12or/api-rest-golang.git
go get  github.com/buaazp/fasthttprouter
go get	github.com/lib/pq
go get	github.com/valyala/fasthttp
```

Before running API server, you should set the database config with yours.
```go
sql.Open("postgres", "postgresql://challenge@localhost:26257/challenge?sslmode=disable")
```
Modify the corsAllowOrigin var for the frontend endpoint choosen by you
```go
corsAllowOrigin      = "http://localhost:8084"
```
# Run cockroach database 3 nodes

```bash
cockroach start \
--insecure \
--store=node1 \
--listen-addr=localhost:26257 \
--http-addr=localhost:8080 \
--join=localhost:26257,localhost:26258,localhost:26259 \
--background

cockroach start \
--insecure \
--store=node2 \
--listen-addr=localhost:26258 \
--http-addr=localhost:8081 \
--join=localhost:26257,localhost:26258,localhost:26259 \
--background

cockroach start \
--insecure \
--store=node3 \
--listen-addr=localhost:26259 \
--http-addr=localhost:8082 \
--join=localhost:26257,localhost:26258,localhost:26259 \
--background

cockroach init --insecure --host=localhost:26257

```

# Stop cockroach nodes

```bash
cockroach quit --insecure --host=localhost:26257
cockroach quit --insecure --host=localhost:26258
cockroach quit --insecure --host=localhost:26259
```

# Build and Run

```bash
cd api-rest-golang
go run main.go

```
## API Endpoint 

```
http://127.0.0.1:9000
```
## Structure
```
├─ api-rest-golang
│  ├─ README.md
│  ├─ handlers
│  │  └─ handlers.go // Handlers for application
│  ├─ main.go
│  ├─ models
│  │  ├─ domain.go // Model for domain structure
│  │  └─ response.go // Models for differents responses
```

## API

#### /projects
* `GET` : Get all domains requested
* `POST` : Request for a new domain information

## Domain Search

- [x] Support basic REST APIs.
- [ ] Support Authentication with user for securing the APIs.
- [ ] Make convenient wrappers for creating API handlers.
- [ ] Write the tests for all APIs.
- [x] Organize the code with packages
- [ ] Make docs with GoDoc
- [ ] Building a deployment process 
