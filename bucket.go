package shoebox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	_ "log/slog"
	"net/url"

	"github.com/aaronland/go-picturebook/bucket"
	"github.com/sfomuseum/go-sfomuseum-api/client"
	"github.com/sfomuseum/go-sfomuseum-api/response"
)

type ShoeboxBucket struct {
	bucket.Bucket
	api_client client.Client
}

func NewShoeboxBucket(ctx context.Context, uri string) (bucket.Bucket, error) {

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

	b := &ShoeboxBucket{
		api_client: api_client,
	}

	return b, nil
}

func (b *ShoeboxBucket) GatherPictures(ctx context.Context, uris ...string) iter.Seq2[string, error] {

	// https://api.sfomuseum.org/methods/sfomuseum.you.shoebox.listItems

	return func(yield func(string, error) bool) {

		args := &url.Values{}
		args.Set("method", "sfomuseum.you.shoebox.listItems")

		cb := func(ctx context.Context, r io.ReadSeekCloser, err error) error {

			if err != nil {
				return err
			}

			var items_rsp *response.ShoeboxListItemsResponse

			dec := json.NewDecoder(r)
			err = dec.Decode(&items_rsp)

			if err != nil {
				return err
			}

			for _, i := range items_rsp.Items {
				str_id := fmt.Sprintf("%d:%d", i.TypeId, i.ItemId)
				yield(str_id, nil)
			}

			return nil
		}

		err := client.ExecuteMethodPaginatedWithClient(ctx, b.api_client, args, cb)

		if err != nil {
			yield("", err)
			return
		}
	}
}

func (b *ShoeboxBucket) NewReader(ctx context.Context, key string, opts any) (io.ReadSeekCloser, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *ShoeboxBucket) NewWriter(ctx context.Context, key string, opts any) (io.WriteCloser, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *ShoeboxBucket) Delete(ctx context.Context, key string) error {
	return nil
}

func (b *ShoeboxBucket) Attributes(ctx context.Context, key string) (*bucket.Attributes, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *ShoeboxBucket) Close() error {
	return nil
}
