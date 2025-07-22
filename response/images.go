package response

import (
	"fmt"
	"strings"
)

// ImageCaptionResponse defines the response object returned by the `sfomuseum.collection.images.getCaption` API method.
type ImageCaptionResponse struct {
	// Caption is an `ImageCaption` instance.
	Caption *ImageCaption
}

// ImageCaption defines the component parts of a caption for an object image in the SFO Museum Aviation Collection.
type ImageCaption struct {
	// Title is the title of the object.
	Title string `json:"title"`
	// Date is the date attributed to the object.
	Date string `json:"date"`
	// Creditline is the credit line for the object.
	CreditLine string `json:"creditline"`
	// AccessionNumber if the object's SFO Museum accession number.
	AccessionNumber string `json:"accession_number"`
	// URL is the collection.sfomuseum.org URL for the object.
	URL string `json:"url"`
}

// String() returns the image caption for an object image in the SFO Museum Aviation Collection as a line-separated string.
func (r *ImageCaption) String() string {

	lines := []string{
		fmt.Sprintf("%s, %s", r.Title, r.Date),
	}

	if r.CreditLine != "Collection of SFO Museum" {
		lines = append(lines, fmt.Sprintf("Collection of SFO Museum, %s", r.AccessionNumber))
		lines = append(lines, r.CreditLine)
	} else {
		lines = append(lines, fmt.Sprintf("%s, %s", r.CreditLine, r.AccessionNumber))
	}

	lines = append(lines, "")
	lines = append(lines, r.URL)

	return strings.Join(lines, "\n")
}
