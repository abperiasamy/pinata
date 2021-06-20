/*
Copyright © 2021 Anand Babu Periasamy https://twitter.com/abperiasamy

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type LichessClient struct {
	BaseURL    string
	UserAgent  string
	AuthTok    string
	HttpClient *http.Client
}

// To generate a personal oAuth API access token from Lichess, follow
// https://lichess.org/account/oauth/token/create
func NewLichessClient(authToken, userAgent string) (lic *LichessClient) {
	lic = new(LichessClient)
	lic.BaseURL = "https://lichess.org/api/"
	lic.UserAgent = userAgent
	lic.AuthTok = authToken
	lic.HttpClient = new(http.Client)
	return lic
}

// Import a game PGN file into lichess
func (c *LichessClient) Import(pgnFile string) (id, gameURL string, err error) {
	pgn, err := ioutil.ReadFile(pgnFile)
	if err != nil {
		fmt.Println("Unable to read "+pgnFile, err)
		return
	}

	// PGN file has to be encoded into URL form.
	data := url.Values{}
	data.Set("pgn", string(pgn))

	req, err := http.NewRequest(http.MethodPost, c.BaseURL+"import", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Unable to initialize HTTP request,", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+c.AuthTok)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		fmt.Println("Unable to make a HTTP request,", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println("Failed to export \""+pgnFile+"\" to Lichess.org.", resp.Status)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read the response,", err)
		return
	}

	importResp := struct {
		Id  string
		URL string // URL to the uploaded game is of the form https://lichess.org/$Id
	}{}

	err = json.Unmarshal(body, &importResp)
	return importResp.Id, importResp.URL, err
}

/*
func main() {
	// lic := NewLichessClient("YSCaKB6WIGawIxI8")
	lic := NewLichessClient("PYN5ZDnwrorecbso", "Piñata")
	id, url, _ := lic.Import("/home/ab/pinata.pgn")
	fmt.Println(id, url)
}
*/
