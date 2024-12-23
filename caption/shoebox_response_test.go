package caption

import (
	"encoding/json"
	_ "fmt"
	"testing"
)

func TestImageCaptionResponse(t *testing.T) {

	body := `{
	"caption": {
		"title": "postcard: American Airlines, Canada",
		"date": "c. 1950",
		"creditline": "Gift of Thomas G. Dragges",
		"accession_number": "2015.166.0309",
		"url": "https://api.sfomuseum.org/objects/1762694275/"
	},
	"stat": "ok"
}`

	expected := `postcard: American Airlines, Canada
c. 1950
Gift of Thomas G. Dragges
Collection of SFO Museum
2015.166.0309`

	var caption_r *ImageCaptionResponse

	err := json.Unmarshal([]byte(body), &caption_r)

	if err != nil {
		t.Fatalf("Failed to unmarshal caption, %v", err)
	}

	str_caption := caption_r.Caption.String()

	if str_caption != expected {
		t.Fatalf("Unexpected string caption '%s'", str_caption)
	}
}
