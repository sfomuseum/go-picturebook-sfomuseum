package oembed

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OEmbed struct {
	Version                  string   `json:"version"`
	Type                     string   `json:"type"`
	Height                   int      `json:"height"`
	Width                    int      `json:"width"`
	Title                    string   `json:"title"`
	URL                      string   `json:"url"`
	AuthorName               string   `json:"author_name"`
	AuthorURL                string   `json:"author_url"`
	ProviderName             string   `json:"provider_name"`
	ProviderURL              string   `json:"provider_url"`
	SFOMuseumAccessionNumber string   `json:"sfomuseum:accession_number,omitempty"`
	SFOMuseumDate            string   `json:"sfomuseum:date,omitempty"`
	SFOMuseumCreditline      string   `json:"sfomuseum:creditline,omitempty"`
	InstagramHashTags        []string `json:"ig:hashtags,omitempty"`
	InstagramDate            string   `json:"ig:date,omitempty"`
}

func Unmarshal(r io.Reader) (*OEmbed, error) {

	var o *OEmbed
	dec := json.NewDecoder(r)
	err := dec.Decode(&o)

	return o, err
}

func Fetch(uri string) (*OEmbed, error) {

	rsp, err := http.Get(uri)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", rsp.StatusCode, rsp.Status)
	}

	return Unmarshal(rsp.Body)
}
