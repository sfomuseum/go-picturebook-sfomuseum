package caption

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	pb_bucket "github.com/aaronland/go-picturebook/bucket"
	pb_caption "github.com/aaronland/go-picturebook/caption"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/sfomuseum/go-sfomuseum-api/client"
)

type ShoeboxCaption struct {
	pb_caption.Caption
	api_client client.Client
	cache      *ristretto.Cache[string, string]
}

func init() {

	ctx := context.Background()
	err := pb_caption.RegisterCaption(ctx, "shoebox", NewShoeboxCaption)

	if err != nil {
		panic(err)
	}
}

func NewShoeboxCaption(ctx context.Context, uri string) (pb_caption.Caption, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()
	token := q.Get("token")

	client_uri := fmt.Sprintf("oauth2://?access_token=%s", token)

	api_client, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new client, %w", err)
	}

	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	c := &ShoeboxCaption{
		cache:      cache,
		api_client: api_client,
	}

	return c, nil
}

func (c *ShoeboxCaption) Text(ctx context.Context, b pb_bucket.Bucket, key string) (string, error) {

	// https://api.sfomuseum.org/methods/sfomuseum.collection.images.getCaption

	logger := slog.Default()
	logger = logger.With("key", key)

	str_caption, found := c.cache.Get(key)

	if found {
		return str_caption, nil
	}

	base := filepath.Base(key)

	// Please use a regexp...
	parts := strings.Split(base, "_")
	image_id := parts[0]

	logger = logger.With("image", image_id)

	args := &url.Values{}
	args.Set("method", "sfomuseum.collection.images.getCaption")
	args.Set("image_id", image_id)

	r, err := c.api_client.ExecuteMethod(ctx, http.MethodGet, args)

	if err != nil {
		logger.Error("Failed to get caption", "error", err)
		return "", err
	}

	var caption_rsp *ImageCaptionResponse

	dec := json.NewDecoder(r)
	err = dec.Decode(&caption_rsp)

	if err != nil {
		logger.Error("Failed to decode caption", "error", err)
		return "", err
	}

	str_caption = caption_rsp.Caption.String()
	c.cache.Set(key, str_caption, 1)

	return str_caption, nil
}
