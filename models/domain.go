package models

import (
	"database/sql"
	"encoding/json"
	"strconv"

	// is nedded for sql querys
	_ "github.com/lib/pq"
)

// Domain : model for the domain struct
type Domain struct {
	ID     int      `json:"id"`
	Domain string   `json:"domain"`
	Data   Response `json:"data"`
}

// DomainCollection : return domain collection
type DomainCollection struct {
	Domains []Domain `json:"items"`
	Prev    string   `json:"prev"`
	Next    string   `json:"next"`
	Rows    int      `json:"rows"`
}

// Migrate : create the table domains in database
func Migrate(db *sql.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS domains(
			id SERIAL PRIMARY KEY,
			domain VARCHAR(256) NOT NULL,
			data JSON NOT NULL,
			UNIQUE (domain)
	);`
	_, err := db.Exec(sql)

	return err
}

// GetDomains : return domain collection
func GetDomains(db *sql.DB, action string, id string) (DomainCollection, error) {
	var sql = ""
	switch action {
	case "prev":
		sql = "SELECT * FROM domains WHERE id < $1 ORDER BY id ASC LIMIT 10"
		break
	case "next":
		sql = "SELECT * FROM domains WHERE id > $1 ORDER BY id ASC LIMIT 10"
		break
	default:
		sql = "SELECT * FROM domains WHERE id > $1 ORDER BY id ASC LIMIT 10"
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return DomainCollection{}, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		return DomainCollection{}, err
	}
	result := DomainCollection{}
	for rows.Next() {
		domain := Domain{}
		err := rows.Scan(&domain.ID, &domain.Domain, &domain.Data)
		if err != nil {
			return DomainCollection{}, err
		}
		result.Domains = append(result.Domains, domain)
	}
	rowsNumber := len(result.Domains)
	if rowsNumber > 0 {
		result.Prev = strconv.Itoa(result.Domains[0].ID)
		result.Next = strconv.Itoa(result.Domains[rowsNumber-1].ID)
		result.Rows = rowsNumber
	}
	return result, nil
}

// CreateDomain : create domain into database
func CreateDomain(db *sql.DB, domain string, data Response) (int64, error) {
	sql := "INSERT INTO domains (domain, data) VALUES($1, $2)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	dataJSONEencode, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(domain, dataJSONEencode)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// UpdateDomain : update domain in database
func UpdateDomain(db *sql.DB, id int, data Response) (int64, error) {
	sql := "UPDATE domains SET (data) = ($1) WHERE id = $2"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	dataJSONEencode, _ := json.Marshal(data)
	result, err := stmt.Exec(dataJSONEencode, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// CheckIfDomainExists : check if a domain already have been registered
func CheckIfDomainExists(db *sql.DB, domain string) (Domain, error) {
	sql := "SELECT id, domain, data FROM domains WHERE domain = $1"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return Domain{}, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(domain)
	if err != nil {
		return Domain{}, err
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
