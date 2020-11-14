package shlink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	apiKey      string      = os.Getenv("SHLINK_API_KEY")
	urlProvider string      = os.Getenv("SHORT_LINK_PROVIDER")
	storageHost string      = os.Getenv("STORAGE_HOST")
	client      http.Client = http.Client{}
)

type ShortLinkOption struct {
	LongURL         string   `json:"longUrl,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	ValidSince      string   `json:"validSince,omitempty"`
	ValidUntil      string   `json:"validUntil,omitempty"`
	CustomSlug      string   `json:"customSlug,omitempty"`
	MaxVisits       int      `json:"maxVisits,omitempty"`
	FindIfExist     bool     `json:"findIfExist,omitempty"`
	Domain          string   `json:"domain,omitempty"`
	ShortCodeLength int      `json:"shortCodeLength,omitempty"`
	ValidateURL     bool     `json:"validateUrl,omitempty"`
}

func CreateLink(filename string) ([]byte, error) {
	body, _ := json.Marshal(ShortLinkOption{
		LongURL: storageHost + "/" + filename,
		Tags:    []string{"upload"},
	})

	fmt.Println(string(body))

	r, _ := http.NewRequest(http.MethodPost, urlProvider+"/rest/v2/short-urls", bytes.NewBuffer(body))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Api-Key", apiKey)

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
