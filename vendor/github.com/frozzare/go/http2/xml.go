package http2

import (
	"encoding/xml"
	"net/http"
)

// GetXML will bind XML response from a url or return a error.
func (s *Client) GetXML(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/xml")

	res, err := s.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return xml.NewDecoder(res.Body).Decode(target)
}

// GetXML will bind XML response from a url or return a error.
func GetXML(url string, target interface{}) error {
	return NewClient(nil).GetXML(url, target)
}
