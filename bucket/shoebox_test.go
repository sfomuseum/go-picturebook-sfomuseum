package bucket

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"testing"
)

var token = flag.String("token", "", "...")

func TestGatherPictures(t *testing.T) {

	if *token == "" {
		t.Skip()
	}

	ctx := context.Background()

	bucket_uri := fmt.Sprintf("shoebox://?token=%s", *token)

	b, err := NewShoeboxBucket(ctx, bucket_uri)

	if err != nil {
		t.Fatalf("Failed to create shoebox bucket, %v", err)
	}

	for id, err := range b.GatherPictures(ctx) {

		if err != nil {
			t.Fatalf("Failed to gather pictures, %v", err)
		}

		slog.Info(id)
	}
}
