package response

import (
	_ "fmt"
	"strings"
)

type ImageCaptionResponse struct {
	Caption *ImageCaption
}

type ImageCaption struct {
	Title           string `json:"title"`
	Date            string `json:"date"`
	CreditLine      string `json:"creditline"`
	AccessionNumber string `json:"accession_number"`
	URL             string `json:"url"`
}

func (r *ImageCaption) String() string {

	lines := []string{
		r.Title,
		r.Date,
		r.CreditLine,
		r.AccessionNumber,
	}

	return strings.Join(lines, "\n")
}
