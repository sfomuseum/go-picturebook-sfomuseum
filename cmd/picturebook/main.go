// picturebook is a command-line application for creating a PDF file from a folder containing images.
package main

import (
	"context"
	"log"

	_ "github.com/sfomuseum/go-picturebook-shoebox"
	_ "gocloud.dev/blob/fileblob"

	"github.com/aaronland/go-picturebook/app/picturebook"
)

func main() {

	ctx := context.Background()
	err := picturebook.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run picturebook application, %w", err)
	}

}
