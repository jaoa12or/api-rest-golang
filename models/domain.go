package models

import (
	"database/sql"
	"encoding/json"

	// is nedded
	_ "github.com/lib/pq"
)

// Domain model Poll
type Domain struct {
	ID     int      `json:"id"`
	Domain string   `json:"domain"`
	Data   Response `json:"data"`
}

// DomainCollection return poll collection
type DomainCollection struct {
	Domains []Domain `json:"items"`
	Prev    int      `json:"prev"`
	Next    int      `json:"next"`
	Rows    int      `json:"rows"`
}

// Migrate return poll collection
func Migrate(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS domains(
			id SERIAL PRIMARY KEY,
			domain VARCHAR(256) NOT NULL,
			data JSON NOT NULL,
			UNIQUE (domain)
	);`
	_, err := db.Exec(sql)

	if err != nil {
		panic(err)
	}
}

// GetDomains return poll collection
func GetDomains(db *sql.DB, action string, id string) DomainCollection {
	var sql = ""
	switch action {
	case "prev":
		sql = "SELECT * FROM domains WHERE id < $1 LIMIT 10"
		break
	case "next":
		sql = "SELECT * FROM domains WHERE id > $1 LIMIT 10"
		break
	default:
		sql = "SELECT * FROM domains WHERE id > $1 LIMIT 10"
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		panic(err)
	}
	result := DomainCollection{}
	for rows.Next() {
		domain := Domain{}
		err := rows.Scan(&domain.ID, &domain.Domain, &domain.Data)
		if err != nil {
			panic(err)
		}
		result.Domains = append(result.Domains, domain)
	}
	rowsNumber := len(result.Domains)
	if rowsNumber > 0 {
		result.Prev = result.Domains[0].ID
		result.Next = result.Domains[rowsNumber-1].ID
		result.Rows = rowsNumber
	}
	return result
}

// CreateDomain create domain
func CreateDomain(db *sql.DB, domain string, data Response) (int64, error) {
	sql := "INSERT INTO domains (domain, data) VALUES($1, $2)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	dataJSONEencode, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	println(dataJSONEencode)
	result, err := stmt.Exec(domain, dataJSONEencode)
	if err != nil {
		panic(err)
	}
	return result.RowsAffected()
}

// UpdateDomain update domain
func UpdateDomain(db *sql.DB, id int, data Response) (int64, error) {
	sql := "UPDATE domains SET (data) = ($1) WHERE id = $2"
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	dataJSONEencode, _ := json.Marshal(data)
	result, err := stmt.Exec(dataJSONEencode, id)
	if err != nil {
		panic(err)
	}
	return result.RowsAffected()
}

// CheckIfDomainExists check if a domain already have been registered
func CheckIfDomainExists(db *sql.DB, domain string) (Domain, error) {
	sql := "SELECT id, domain, data FROM domains WHERE domain = $1"
	stmt, err := db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(domain)
	if err != nil {
		panic(err)
	}
	result := Domain{}
	for rows.Next() {
		err := rows.Scan(&result.ID, &result.Domain, &result.Data)
		if err != nil {
			return result, err
		}
	}
	return result, err
}
