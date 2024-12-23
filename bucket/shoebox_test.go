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

func TestAttributes(t *testing.T) {

	ctx := context.Background()

	bucket_uri := fmt.Sprintf("shoebox://?token=%s", "TOKEN")

	b, err := NewShoeboxBucket(ctx, bucket_uri)

	if err != nil {
		t.Fatalf("Failed to create shoebox bucket, %v", err)
	}

	k := "https://static.sfomuseum.org/media/191/366/340/9/1913663409_MSM9QjCaQmXnyemSonODPufdrayFWc4a_k.jpg"

	_, err = b.Attributes(ctx, k)

	if err != nil {
		t.Fatalf("Failed to derive attributes for %s, %v", k, err)
	}
}
