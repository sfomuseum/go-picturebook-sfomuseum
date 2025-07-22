// package client provides methods for SFO Museum API clients.
package client

import (
	"context"
	"fmt"
	"io"
	"iter"
	"net/url"
	"strconv"

	"github.com/sfomuseum/go-sfomuseum-api/v2/response"
)

// API_ENDPOINT is the default endpoint for the api.sfomuseum.org API.
const API_ENDPOINT string = "https://api.sfomuseum.org/rest"

// type Client is an interface for SFO Museum API client implementations.
type Client interface {
	// ExecuteMethod performs an API method request.
	ExecuteMethod(context.Context, string, *url.Values) (io.ReadSeekCloser, error)
}

// ExecuteMethodPaginatedWithClient performs as many paginated API requests for a given method. It returns
// a `iter.Seq2[io.ReadSeeker, error]` instance which can be iterated through to yield each API result.
func ExecuteMethodPaginatedWithClient(ctx context.Context, cl Client, verb string, args *url.Values) iter.Seq2[io.ReadSeeker, error] {

	return func(yield func(io.ReadSeeker, error) bool) {

		page := 1
		pages := -1

		if args.Get("page") == "" {
			args.Set("page", strconv.Itoa(page))
		} else {

			p, err := strconv.Atoi(args.Get("page"))

			if err != nil {
				yield(nil, fmt.Errorf("Invalid page number '%s', %v", args.Get("page"), err))
				return
			}

			page = p
		}

		for {

			select {
			case <-ctx.Done():
				break
			default:
				// pass
			}

			r, err := cl.ExecuteMethod(ctx, verb, args)

			if err != nil {
				yield(nil, err)
				break
			}

			defer r.Close()

			if !yield(r, nil) {
				break
			}

			_, err = r.Seek(0, 0)

			if err != nil {
				yield(nil, fmt.Errorf("Failed to rewind response, %v", err))
				break
			}

			if pages == -1 {

				pagination, err := response.DerivePagination(ctx, r)

				if err != nil {
					yield(nil, err)
					break
				}

				pages = pagination.Pages
			}

			page += 1

			if page <= pages {
				args.Set("page", strconv.Itoa(page))
			} else {
				break
			}
		}
	}

}
