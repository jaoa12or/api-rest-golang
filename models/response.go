package models

import (
	"encoding/json"
	"fmt"
)

// Endpoint model SslResponse
type Endpoint struct {
	Address string `json:"ipAddress"`
	Grade   string `json:"grade"`
	Country string `json:"country"`
	Owner   string `json:"owner"`
}

// Endpoints model SslResponse
type Endpoints struct {
	Endpoints []Endpoint `json:"endpoints"`
}

// WhoisResponse model WhoisResponse
type WhoisResponse struct {
	Country string `json:"Country"`
	Owner   string `json:"OrgName"`
}

// ScrapingResponse model ScrapingResponse
type ScrapingResponse struct {
	Icon  string
	Title string
}

// BadRequest :
type BadRequest struct{
	Response string `json:"response"`
}

// Response model Response
type Response struct {
	Servers          []Endpoint `json:"servers"`
	ServerChanged    bool       `json:"server_changed"`
	SslGrade         string     `json:"ssl_grade"`
	PreviousSslGrade string     `json:"previous_ssl_grade"`
	Logo             string     `json:"logo"`
	Title            string     `json:"title"`
	IsDown           bool       `json:"is_down"`
}

// Scan model Response
func (response *Response) Scan(src interface{}) error {
	strValue, ok := src.([]uint8)
	if !ok {
		return fmt.Errorf("metas field must be a []uint8, got %T instead", src)
	}
	return json.Unmarshal([]byte(strValue), response)
}
