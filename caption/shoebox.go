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
	"github.com/sfomuseum/go-picturebook-sfomuseum/response"
	"github.com/sfomuseum/go-sfomuseum-api/client"
)

// ShoeboxCaption implements the `aaronland/go-picturebook/caption.Caption` interface for use with object images in a SFO Museum "shoebox".
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

// NewShoeboxCaption returns a new `ShoeboxCaption` instance implementing the `aaronland/go-picturebook/caption.Caption` interface for use with object images in a SFO Museum "shoebox".
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

// Text returns the caption text for object image identified by 'key' in 'b'.
func (c *ShoeboxCaption) Text(ctx context.Context, b pb_bucket.Bucket, key string) (string, error) {

	// https://api.sfomuseum.org/methods/sfomuseum.collection.images.getCaption

	logger := slog.Default()
	logger = logger.With("key", key)

	str_caption, found := c.cache.Get(key)

	if found {
		return str_caption, nil
	}

	base := filepath.Base(key)
	parts := strings.Split(base, "#")

	switch len(parts) {
	case 2:

		slog.Info("YO")

		fragment := strings.Split(parts[1], ":")

		switch fragment[0] {
		case "ig":

			post_id := fragment[1]
			slog.Info("IG", "post id", post_id)

			ig_args := &url.Values{}
			ig_args.Set("method", "sfomuseum.millsfield.instagram.getInfo")
			ig_args.Set("post_id", post_id)

			ig_rsp, err := c.api_client.ExecuteMethod(ctx, http.MethodGet, ig_args)

			if err != nil {
				slog.Info("WOMP", "error", err)
				return "", fmt.Errorf("Failed to execute sfomuseum.millsfield.instagram.getInfo method, %w", err)
			}

			defer ig_rsp.Close()
			var ig_post_rsp *response.InstagramPostResponse

			dec := json.NewDecoder(ig_rsp)
			err = dec.Decode(&ig_post_rsp)

			if err != nil {
				slog.Info("WOMP 2", "error", err)
				return "", fmt.Errorf("Failed to unmarshal IG post response, %w", err)
			}

			text := []string{
				ig_post_rsp.Post.Caption.Excerpt,
				"This was posted to the SFO Museum Instagram account on {DATE}",
			}

			str_text := strings.Join(text, "\n")

			slog.Info("IG", "caption", str_text)

			return str_text, nil

		default:
			slog.Info("SAD")
			return "", fmt.Errorf("Unhandled or unsupported fragment, %s", fragment[0])
		}

	default:

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

		var caption_rsp *response.ImageCaptionResponse

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
}
