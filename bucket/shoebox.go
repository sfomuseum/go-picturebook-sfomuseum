package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	pb_bucket "github.com/aaronland/go-picturebook/bucket"
	"github.com/jtacoma/uritemplates"
	"github.com/sfomuseum/go-sfomuseum-api/client"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-ioutil"
)

type ShoeboxBucket struct {
	pb_bucket.Bucket
	api_client client.Client
}

func init() {

	ctx := context.Background()
	err := pb_bucket.RegisterBucket(ctx, "shoebox", NewShoeboxBucket)

	if err != nil {
		panic(err)
	}
}

func NewShoeboxBucket(ctx context.Context, uri string) (pb_bucket.Bucket, error) {

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

			var items_rsp *ShoeboxListItemsResponse

			dec := json.NewDecoder(r)
			err = dec.Decode(&items_rsp)

			if err != nil {
				return err
			}

			for _, i := range items_rsp.Items {

				// Object (fetch type map rather than hardcoding things...)

				if i.TypeId != 1 {
					slog.Debug("Item type not supported", "item id", i.ItemId, "type", i.TypeId)
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

					for im_idx, r := range im_rsp.Array() {

						logger := slog.Default()
						logger = logger.With("object", i.ItemId)
						logger = logger.With("image", im_idx)

						template_r := r.Get("media:uri_template")

						if !template_r.Exists() {
							logger.Warn("Image record for object is missing media:uri_template, skipping")
							continue
						}

						sizes_r := r.Get("media:properties.sizes")

						if !sizes_r.Exists() {
							logger.Warn("Image record for object is missing media:properties.sizes, skipping")
							continue
						}

						uri_t, err := uritemplates.Parse(template_r.String())

						if err != nil {
							logger.Warn("Failed to parse URI template for object image", "error", err)
							continue
						}

						vars := make(map[string]interface{})

						labels := []string{
							"o",
							"k",
							"b",
							"c",
						}

						for _, label := range labels {

							l_rsp := sizes_r.Get(label)

							if !l_rsp.Exists() {
								continue
							}

							vars["label"] = label
							vars["secret"] = l_rsp.Get("secret").String()
							vars["extension"] = l_rsp.Get("extension").String()
							break
						}

						_, has_label := vars["label"]

						if !has_label {
							logger.Warn("Image missing label after reading sizes, skipping")
							continue
						}

						im_uri, err := uri_t.Expand(vars)

						if err != nil {
							logger.Warn("Failed to expand URI template vars, skipping", "error", err)
							continue
						}

						yield(im_uri, nil)

						/*
							id_r := r.Get("wof:id")

							if (! id_r.Exists()){
								slog.Warn("Image record for object is missing wof:id", "object", i.ItemId)
								continue
							}

							compound_id := fmt.Sprintf("%d:%d", i.TypeId, id_r.Int())
							yield(compound_id, nil)
						*/
					}

					return nil
				}

				im_err := client.ExecuteMethodPaginatedWithClient(ctx, b.api_client, http.MethodGet, im_args, im_cb)

				if im_err != nil {
					return fmt.Errorf("Failed to retrieve images for object %d, %w", i.ItemId, im_err)
				}
			}

			return nil
		}

		err := client.ExecuteMethodPaginatedWithClient(ctx, b.api_client, http.MethodGet, args, cb)

		if err != nil {
			yield("", err)
			return
		}
	}
}

func (b *ShoeboxBucket) NewReader(ctx context.Context, key string, opts any) (io.ReadSeekCloser, error) {

	// Valid key here...

	rsp, err := http.Get(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve %s, %w", key, err)
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to retrieve %s, %d %s", key, rsp.StatusCode, rsp.Status)
	}

	return ioutil.NewReadSeekCloser(rsp.Body)
}

func (b *ShoeboxBucket) NewWriter(ctx context.Context, key string, opts any) (io.WriteCloser, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *ShoeboxBucket) Delete(ctx context.Context, key string) error {
	return fmt.Errorf("Not implemented")
}

func (b *ShoeboxBucket) Attributes(ctx context.Context, key string) (*pb_bucket.Attributes, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *ShoeboxBucket) Close() error {
	return nil
}
