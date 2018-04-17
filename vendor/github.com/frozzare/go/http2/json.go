package http2

import (
	"encoding/json"
	"net/http"
)

// GetJSON will bind JSON response from a url or return a error.
func (s *Client) GetJSON(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := s.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(target)
}

// GetJSON will bind JSON response from a url or return a error.
func GetJSON(url string, target interface{}) error {
	return NewClient(nil).GetJSON(url, target)
}
