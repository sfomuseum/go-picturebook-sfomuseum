package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/whosonfirst/go-ioutil"
)

func init() {

	ctx := context.Background()
	err := RegisterClient(ctx, "oauth2", NewOAuth2Client)

	if err != nil {
		panic(err)
	}
}

// OAuth2Client implements the `Client` interface for API endpoints that require OAuth2 access token authentication.
type OAuth2Client struct {
	Client
	// http_client is the underlying net/http client used to perform API requests.
	http_client *http.Client
	// api_endpoint is the URL of the API endpoint.
	api_endpoint string
	// access_token is the OAuth2 access token to append to API requests.
	access_token string
}

// NewOAuth2Client creates a new `OAuth2Client` instance configured by 'uri' which
// is expected to take the form of:
//
//	oauth2://{HOST}/{PATH}?{PARAMETERS}
//
// Where {PARAMETERS} is:
// - `?access_token={TOKEN}` A valid OAuth2 access token.
// - `?insecure={BOOLEAN}` A boolean flag signaling that TLS verification will be skipped.
func NewOAuth2Client(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	var api_endpoint string

	switch u.Host {
	case "", "api":
		api_endpoint = API_ENDPOINT
	default:
		api_endpoint = fmt.Sprintf("https://%s/%s", u.Host, u.Path)
	}

	q := u.Query()

	http_client := &http.Client{}

	if q.Has("insecure") {

		v, err := strconv.ParseBool(q.Get("insecure"))

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?insecure=, %w", err)
		}

		if v {

			config := &tls.Config{
				InsecureSkipVerify: true,
			}

			tr := &http.Transport{
				TLSClientConfig: config,
			}

			http_client = &http.Client{
				Transport: tr,
			}
		}
	}

	access_token := q.Get("access_token")

	cl := &OAuth2Client{
		http_client:  http_client,
		access_token: access_token,
		api_endpoint: api_endpoint,
	}

	return cl, nil
}

// ExecuteMethod will perform an API request derived from 'args'.
func (cl *OAuth2Client) ExecuteMethod(ctx context.Context, args *url.Values) (io.ReadSeekCloser, error) {

	endpoint, err := url.Parse(cl.api_endpoint)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse endpoint URI, %w", err)
	}

	http_method := "GET"

	if cl.access_token != "" {
		args.Set("access_token", cl.access_token)
	}

	endpoint.RawQuery = args.Encode()

	req, err := http.NewRequest(http_method, endpoint.String(), nil)

	if err != nil {
		return nil, err
	}

	return cl.executeRequest(ctx, req)
}

// executeRequest will perform an API request derived from 'req'.
func (cl *OAuth2Client) executeRequest(ctx context.Context, req *http.Request) (io.ReadSeekCloser, error) {

	req = req.WithContext(ctx)

	rsp, err := cl.http_client.Do(req)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		rsp.Body.Close()
		return nil, fmt.Errorf("API call failed with status '%s'", rsp.Status)
	}

	return ioutil.NewReadSeekCloser(rsp.Body)
}
