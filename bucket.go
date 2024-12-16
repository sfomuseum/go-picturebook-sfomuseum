package shoebox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/aaronland/go-picturebook/bucket"
	"github.com/sfomuseum/go-sfomuseum-api/client"
	"github.com/sfomuseum/go-sfomuseum-api/response"
	"github.com/tidwall/gjson"
)

type ShoeboxBucket struct {
	bucket.Bucket
	api_client client.Client
}

func init() {

	ctx := context.Background()
	err := bucket.RegisterBucket(ctx, "shoebox", NewShoeboxBucket)

	if err != nil {
		panic(err)
	}
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
	// https://api.sfomuseum.org/methods/sfomuseum.collection.objects.getImages

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

				// Object (fetch type map rather than hardcoding things...)

				if i.TypeId != 1 {
					slog.Warn("Item type not supported", "item id", i.ItemId, "type", i.TypeId)
					continue
				}

				str_id := strconv.FormatInt(i.ItemId, 10)

				im_args := &url.Values{}
				im_args.Set("method", "sfomuseum.collection.objects.getImages")
				im_args.Set("object_id", str_id)

				im_cb := func(ctx context.Context, r io.ReadSeekCloser, err error) error {

					if err != nil {
						return err
					}

					// Something something SPR...

					im_body, err := io.ReadAll(r)

					if err != nil {
						return err
					}

					im_rsp := gjson.GetBytes(im_body, "images")

					for _, r := range im_rsp.Array() {
						yield(r.Get("wof:id").String(), nil)
					}

					return nil
				}

				im_err := client.ExecuteMethodPaginatedWithClient(ctx, b.api_client, im_args, im_cb)

				if im_err != nil {
					return fmt.Errorf("Failed to retrieve images for object %d, %w", i.ItemId, im_err)
				}
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
