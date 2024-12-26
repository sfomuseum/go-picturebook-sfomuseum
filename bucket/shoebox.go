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
	"strings"
	"time"

	pb_bucket "github.com/aaronland/go-picturebook/bucket"
	"github.com/jtacoma/uritemplates"
	"github.com/sfomuseum/go-sfomuseum-api/client"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-ioutil"
)

// ShoeboxBucket implements the `aaronland/go-picturebook/bucket.Bucket` interface for use with object images in a SFO Museum "shoebox".
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

// NewShoeboxBucket returns a new `ShoeboxBucket` instance implementing the `aaronland/go-picturebook/bucket.Bucket` interface for use with object images in a SFO Museum "shoebox".
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

// GatherPictures returns a new `iter.Seq2[string, error]` instance containing the URIs for object images in a SFO Museum "shoebox".
func (b *ShoeboxBucket) GatherPictures(ctx context.Context, uris ...string) iter.Seq2[string, error] {

	// https://api.sfomuseum.org/methods/sfomuseum.you.shoebox.listItems
	// https://api.sfomuseum.org/methods/sfomuseum.collection.objects.getImages

	return func(yield func(string, error) bool) {

		args := &url.Values{}
		args.Set("method", "sfomuseum.you.shoebox.listItems")

		list_cb := func(ctx context.Context, r io.ReadSeekCloser, err error) error {

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

				// fetch type map rather than hardcoding things...

				switch i.TypeId {
				case 1: // Objects

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
						}

						return nil
					}

					im_err := client.ExecuteMethodPaginatedWithClient(ctx, b.api_client, http.MethodGet, im_args, im_cb)

					if im_err != nil {
						return fmt.Errorf("Failed to retrieve images for object %d, %w", i.ItemId, im_err)
					}

				case 2: // Instagram

					str_id := strconv.FormatInt(i.ItemId, 10)

					ig_args := &url.Values{}
					ig_args.Set("method", "sfomuseum.millsfield.instagram.getInfo")
					ig_args.Set("post_id", str_id)

					ig_rsp, err := b.api_client.ExecuteMethod(ctx, http.MethodGet, ig_args)

					if err != nil {
						return fmt.Errorf("Failed to execute sfomuseum.millsfield.instagram.getInfo method, %w", err)
					}

					defer ig_rsp.Close()
					var ig_post_rsp *InstagramPostResponse

					dec := json.NewDecoder(ig_rsp)
					err = dec.Decode(&ig_post_rsp)

					if err != nil {
						return fmt.Errorf("Failed to unmarshal IG post response, %w", err)
					}

					ig_post := ig_post_rsp.Post

					// This will not work. Need to derive SFOM URI at the API layer...
					yield(ig_post.Path, nil)

				default:
					slog.Debug("Item type not supported", "item id", i.ItemId, "type", i.TypeId)
					continue
				}

			}

			return nil
		}

		err := client.ExecuteMethodPaginatedWithClient(ctx, b.api_client, http.MethodGet, args, list_cb)

		if err != nil {
			yield("", err)
			return
		}
	}
}

// NewReader returns a new `io.ReadSeekCloser` instance for an object image identified by 'key' in a SFO Museum "shoebox".
func (b *ShoeboxBucket) NewReader(ctx context.Context, key string, opts any) (io.ReadSeekCloser, error) {

	if !b.isValidKey(key) {
		return nil, fmt.Errorf("Invalid key")
	}

	rsp, err := http.Get(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve %s, %w", key, err)
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to retrieve %s, %d %s", key, rsp.StatusCode, rsp.Status)
	}

	return ioutil.NewReadSeekCloser(rsp.Body)
}

// NewWriter returns an error because this package only implements non-destructive methods of the `aaronland/go-picturebook/bucket.Bucket` interface.
func (b *ShoeboxBucket) NewWriter(ctx context.Context, key string, opts any) (io.WriteCloser, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Delete returns an error because this package only implements non-destructive methods of the `aaronland/go-picturebook/bucket.Bucket` interface.
func (b *ShoeboxBucket) Delete(ctx context.Context, key string) error {
	return fmt.Errorf("Not implemented")
}

// Attribute returns a new `aaronland/go-picturebook/bucket.Attributes` instance for an object image identified by 'key' in a SFO Museum "shoebox".
func (b *ShoeboxBucket) Attributes(ctx context.Context, key string) (*pb_bucket.Attributes, error) {

	if !b.isValidKey(key) {
		return nil, fmt.Errorf("Invalid key")
	}

	rsp, err := http.Head(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute request, %w", err)
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed: %d %s", rsp.StatusCode, rsp.Status)
	}

	/*

		> curl -I https://static.sfomuseum.org/media/191/366/340/9/1913663409_MSM9QjCaQmXnyemSonODPufdrayFWc4a_k.jpg
		HTTP/2 200
		content-type: image/jpeg
		content-length: 1737976
		date: Mon, 23 Dec 2024 19:53:51 GMT
		last-modified: Tue, 16 Apr 2024 07:51:12 GMT
		etag: "daf3bd2eb40e880602dd0f8333c9a09c"
		x-amz-server-side-encryption: AES256
		accept-ranges: bytes
		server: AmazonS3
		x-cache: Hit from cloudfront
		via: 1.1 7dbcbf3457f77b741952e31c6826a8dc.cloudfront.net (CloudFront)
		x-amz-cf-pop: SFO53-P7
		x-amz-cf-id: yGF5f7oWVxwXLDPCpWo8JuVdEfkdHKY6kre5sIOtL1QBhL4KteL3dA==
		age: 24

	*/

	str_len := rsp.Header.Get("Content-Length")
	str_lastmod := rsp.Header.Get("Last-Modified")

	content_len, err := strconv.ParseInt(str_len, 10, 64)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse content length, %w", err)
	}

	lastmod, err := time.Parse(time.RFC1123, str_lastmod)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse lastmod time, %w", err)
	}

	attrs := &pb_bucket.Attributes{
		ModTime: lastmod,
		Size:    content_len,
	}

	return attrs, nil
}

// Close completes and terminates any underlying code used by 'b'.
func (b *ShoeboxBucket) Close() error {
	return nil
}

func (b *ShoeboxBucket) isValidKey(key string) bool {
	return strings.HasPrefix(key, "https://static.sfomuseum.org/media")
}
