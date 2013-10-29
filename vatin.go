// Package vatin provides support for interacting with EU VIES VAT number validation
package vatin

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strings"
)

// Result of VATIN validation
type VATINResult struct {
	Valid   bool   `json:"valid"`
	Country string `json:"country,omitempty"`
	Number  string `json:"number,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

// Validates a VATIN number
func Validate(stateNumber string) (*VATINResult, error) {
	
	temp := strings.Replace(stateNumber, " ", "", -1)
	state := temp[0:2]
	number := temp[2:]
	
	resp, err := http.PostForm("http://ec.europa.eu/taxation_customs/vies/vatResponse.html", url.Values{"memberStateCode": {state}, "number": {number}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var result VATINResult
	if doc.Find(".validStyle").Text() == "Yes, valid VAT number" {
		result.Valid = true
	}
	if result.Valid {
		doc.Find(".labelStyle").Each(func(i int, s *goquery.Selection) {
			switch s.Text() {
			case "Member State":
				result.Country = s.Siblings().Text()
			case "VAT Number":
				result.Number = s.Siblings().Text()
			case "Name":
				result.Name = strings.Replace(s.Siblings().Text(), "\n", "", -1)
			case "Address":
				html, _ := s.Siblings().Html()
				result.Address = strings.Replace(strings.Replace(html, "<br/>", "", 1), "<br/>", "\n", -1)
			}
		})
	}
	return &result, nil
}

