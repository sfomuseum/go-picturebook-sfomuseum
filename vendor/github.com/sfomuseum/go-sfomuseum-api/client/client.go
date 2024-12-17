// package client provides methods for SFO Museum API clients.
package client

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/sfomuseum/go-sfomuseum-api/response"
)

// API_ENDPOINT is the default endpoint for the api.sfomuseum.org API.
const API_ENDPOINT string = "https://api.sfomuseum.org/rest"

// type Client is an interface for SFO Museum API client implementations.
type Client interface {
	// ExecuteMethod performs an API method request.
	ExecuteMethod(context.Context, *url.Values) (io.ReadSeekCloser, error)
}

// ExecuteMethodPaginatedCallback is a custom callback function to be invoked by every response item
// seen by the `ExecuteMethodPaginatedWithClient` method.
type ExecuteMethodPaginatedCallback func(context.Context, io.ReadSeekCloser, error) error

// ExecuteMethodPaginatedWithClient performs as many paginated API requests for a given method to yield
// all the result. Each result is passed to the 'cb' callback method for final processing.
func ExecuteMethodPaginatedWithClient(ctx context.Context, cl Client, args *url.Values, cb ExecuteMethodPaginatedCallback) error {

	page := 1
	pages := -1

	if args.Get("page") == "" {
		args.Set("page", strconv.Itoa(page))
	} else {

		p, err := strconv.Atoi(args.Get("page"))

		if err != nil {
			return fmt.Errorf("Invalid page number '%s', %v", args.Get("page"), err)
		}

		page = p
	}

	for {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		fh, err := cl.ExecuteMethod(ctx, args)

		err = cb(ctx, fh, err)

		if err != nil {
			return err
		}

		_, err = fh.Seek(0, 0)

		if err != nil {
			return fmt.Errorf("Failed to rewind response, %v", err)
		}

		if pages == -1 {

			pagination, err := response.DerivePagination(ctx, fh)

			if err != nil {
				return err
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

	return nil
}
