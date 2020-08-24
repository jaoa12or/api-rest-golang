package handlers

import (
	"bytes"
	"challenge-backend/models"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/valyala/fasthttp"
)

// GetDomains : handler for get domain collection from database
func GetDomains(db *sql.DB) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		action := ctx.QueryArgs().Peek("action")
		cursor := ctx.QueryArgs().Peek("cursor")
		domains, err := models.GetDomains(db, string(action), string(cursor))
		if err != nil {
			log.Println(err)
			ctx.Error("The domains can't be consulted", 500)
		}
		if err := json.NewEncoder(ctx).Encode(domains); err != nil {
			log.Println(err)
			ctx.Error("The domains can't be consulted", 500)
		}
	}
}

// GetPageInfo : method to get the web page information, like icon and title
func GetPageInfo(domain string) (models.ScrapingResponse, error) {
	// Request the HTML page.
	var route = "http://" + domain
	res, err := http.Get(route)
	if err != nil {
		log.Println(err)
		return models.ScrapingResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("status code error: ", res.StatusCode, res.Status)
		return models.ScrapingResponse{}, err
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return models.ScrapingResponse{}, err
	}
	var scrapingResponse models.ScrapingResponse
	// Find the head
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the rel and href
		s.Find("link").Each(func(j int, l *goquery.Selection) {
			rel, _ := l.Attr("rel")
			if strings.Contains(rel, "icon") {
				href, _ := l.Attr("href")
				if strings.Contains(href, "http://") || strings.Contains(href, "https://") {
					scrapingResponse.Icon = href
				} else {
					strTrim := strings.TrimLeft(href, "/")
					scrapingResponse.Icon = route + "/" + strTrim
				}
			}
		})

		title := s.Find("title").Text()
		scrapingResponse.Title = strings.Trim(title, "\t \n")
	})
	return scrapingResponse, nil
}

// GetOwnerData : get the owner data of the web site
func GetOwnerData(address string, key string) string {
	// get the owner data trough piped command between whois command and grep
	// to get the exact key information that is been searched
	output := exec.Command("whois", address)
	output2 := exec.Command("grep", "-m 1", key)
	// create the pipe for read and write exec command response
	reader, writer := io.Pipe()
	// buffer that stores response
	var buf bytes.Buffer
	// read and write response for outputs
	output.Stdout = writer
	output2.Stdin = reader
	output2.Stdout = &buf

	// starts commands in cmd
	output.Start()
	output2.Start()
	output.Wait()
	writer.Close()
	output2.Wait()
	reader.Close()

	// read buffer for loking response
	var name = buf.String()
	// delete spaces an tabs from response
	strTrim := strings.Trim(name, key+": \t \n")
	strTrim2 := strings.TrimLeft(strTrim, "\t \n")

	return strTrim2
}

// CalculateGrade : calculate current grade
func CalculateGrade(currentGrade string) string {
	previousGrade := 0
	finalGrade := ""
	// check if grade is set
	if len(currentGrade) > 0 {
		// convert current grade to ascii equivalent
		sslGrade := int([]rune(currentGrade)[0])
		// calculate value for grade ej: if grade A+ = 65 - 43
		if len(currentGrade) > 1 {
			for i := 1; i < len(currentGrade); i++ {
				sslGrade -= int([]rune(currentGrade)[i])
			}
			if sslGrade > previousGrade {
				previousGrade = sslGrade
				finalGrade = currentGrade
			}
		} else {
			if sslGrade > previousGrade {
				previousGrade = sslGrade
				finalGrade = currentGrade
			}
		}
	}
	return finalGrade
}

// setDomainState : set the domain state
func setDomainState(db *sql.DB, domain string) error {
	// if the system can't get any information from ssl server for a domain
	// then the domain is updated to unreachable (is_down = true)
	previousData, err := models.CheckIfDomainExists(db, domain)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(models.Domain{}, previousData) {
		previousData.Data.IsDown = true
		previousData.Data.ServerChanged = true
		_, err := models.UpdateDomain(db, previousData.ID, previousData.Data)
		if err != nil {
			return err
		}
	}
	return nil
}

// ConsultDomain : handler for get all the information required for an user domain request
func ConsultDomain(db *sql.DB) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var response models.Response
		var badRequest models.BadRequest
		var sendResponse = true

		strTrim := string(ctx.FormValue("domain"))
		// get domain from request
		domain := url.QueryEscape(strTrim)
		// get the info of domain from ssllabs
		url := "https://api.ssllabs.com/api/v3/analyze?host=" + domain
		resp, err := http.Get(url)
		if err != nil {
			// if domain is unrecheable, set is_down = true
			err = setDomainState(db, domain)
			if err != nil {
				ctx.SetStatusCode(500)
				log.Println(err)
				badRequest.Response = "Theres not info available for this domain"
			}
			sendResponse = false
			log.Println(err)
			ctx.SetStatusCode(500)
			badRequest.Response = "The domain can't be consulted"
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var sslResponse models.Endpoints
			err := json.NewDecoder(resp.Body).Decode(&sslResponse)
			if err != nil {
				sendResponse = false
				log.Println(err)
				ctx.SetStatusCode(500)
				badRequest.Response = "The domain can't be consulted"
			}

			if len(sslResponse.Endpoints) > 0 {
				finalGrade := ""
				for i := 0; i < len(sslResponse.Endpoints); i++ {
					// get the owner data
					var owner = GetOwnerData(sslResponse.Endpoints[i].Address, "OrgName")
					var country = GetOwnerData(sslResponse.Endpoints[i].Address, "Country")

					sslResponse.Endpoints[i].Owner = owner
					sslResponse.Endpoints[i].Country = country
					currentGrade := sslResponse.Endpoints[i].Grade
					// calculate the biggest grade
					finalGrade = CalculateGrade(currentGrade)
				}

				var pageInfo, err = GetPageInfo(domain)
				if err != nil {
					// if domain is unrecheable, set is_down = true
					err = setDomainState(db, domain)
					if err != nil {
						ctx.SetStatusCode(500)
						log.Println(err)
						badRequest.Response = "Theres not info available for this domain"
					}
					sendResponse = false
					log.Println(err)
					ctx.SetStatusCode(500)
					badRequest.Response = "The domain information can't be consulted"
				} else {
					response.Servers = sslResponse.Endpoints
					response.SslGrade = finalGrade
					response.Logo = pageInfo.Icon
					response.Title = pageInfo.Title
					response.IsDown = false

					// check if domain exists. then deside if create or update the info in database
					previousData, err := models.CheckIfDomainExists(db, domain)
					switch {
					case reflect.DeepEqual(models.Domain{}, previousData):
						response.ServerChanged = false
						response.PreviousSslGrade = finalGrade
						_, err := models.CreateDomain(db, domain, response)
						if err != nil {
							sendResponse = false
							log.Println(err)
							ctx.SetStatusCode(500)
							badRequest.Response = "The domain can't be created"
						}
					case err != nil:
						sendResponse = false
						log.Println(err)
						ctx.SetStatusCode(500)
						badRequest.Response = "The database can't be reached"

					default:
						response.ServerChanged = previousData.Data.ServerChanged
						response.PreviousSslGrade = previousData.Data.SslGrade
						if reflect.DeepEqual(previousData.Data, response) {
							response.ServerChanged = false
							_, err := models.UpdateDomain(db, previousData.ID, response)
							if err != nil {
								sendResponse = false
								log.Println(err)
								ctx.SetStatusCode(500)
								badRequest.Response = "The domain can't be updated"
							}
						} else {
							response.ServerChanged = true
							response.PreviousSslGrade = previousData.Data.SslGrade
							_, err := models.UpdateDomain(db, previousData.ID, response)
							if err != nil {
								sendResponse = false
								log.Println(err)
								ctx.SetStatusCode(500)
								badRequest.Response = "The domain can't be updated"
							}
						}
					}
				}

			} else {
				// if domain is unrecheable, set is_down = true
				err = setDomainState(db, domain)
				if err != nil {
					ctx.SetStatusCode(500)
					log.Println(err)
					badRequest.Response = "Theres not info available for this domain"
				}

				sendResponse = false
				ctx.SetStatusCode(500)
				log.Println("There's not info available for this domain")
				badRequest.Response = "Theres not info available for this domain"
			}
		} else {
			sendResponse = false
			log.Println("The ssl database can't be reached")
			ctx.SetStatusCode(500)
			badRequest.Response = "The ssl database can't be reached"
		}
		if sendResponse {
			if err := json.NewEncoder(ctx).Encode(response); err != nil {
				log.Println(err)
			}
		} else {
			if err := json.NewEncoder(ctx).Encode(badRequest); err != nil {
				log.Println(err)
			}
		}
	}
}
