package shoebox

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-picturebook/caption"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/sfomuseum/go-picturebook-shoebox/oembed"
)

type ShoeboxCaption struct {
	caption.Caption
	cache *ristretto.Cache[string, string]
}

func init() {

	ctx := context.Background()
	err := caption.RegisterCaption(ctx, "shoebox", NewShoeboxCaption)

	if err != nil {
		panic(err)
	}
}

func NewShoeboxCaption(ctx context.Context, uri string) (caption.Caption, error) {

	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	c := &ShoeboxCaption{
		cache: cache,
	}

	return c, nil
}

func (c *ShoeboxCaption) Text(ctx context.Context, b bucket.Bucket, key string) (string, error) {

	caption, found := c.cache.Get(key)

	if found {
		return caption, nil
	}

	oembed_uri := fmt.Sprintf("https://collection.sfomuseum.org/oembed?url=%s", key)
	slog.Info(oembed_uri)

	o, err := oembed.Fetch(oembed_uri)

	if err != nil {
		return "", err
	}

	lines := []string{
		o.Title,
		o.SFOMuseumDate,
		o.SFOMuseumCreditline,
		o.SFOMuseumAccessionNumber,
	}

	str_caption := strings.Join(lines, "\n")
	c.cache.Set(key, str_caption, 1)

	return str_caption, nil
}
