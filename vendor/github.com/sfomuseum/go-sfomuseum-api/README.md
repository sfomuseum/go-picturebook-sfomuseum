# go-sfomuseum-api

Go package providing methods for interacting with the SFO Museum API endpoints.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-sfomuseum-api.svg)](https://pkg.go.dev/github.com/sfomuseum/go-sfomuseum-api)

## Tools

### api

api is a command-line tool for invoking a SFO Museum API emitting the results to STDOUT.

```
$> ./bin/api -h
api is a command-line tool for invoking a SFO Museum API emitting the results to STDOUT.

Usage:
	 ./bin/api [options]

Valid options are:
  -access-token-uri string
    	A valid gocloud.dev/runtime variable URI containing a value to replace '{ACCESS_TOKEN}' in the -api-client-uri flag.
  -api-client-uri string
    	 (default "oauth2://?access_token={ACCESS_TOKEN}")
  -param value
    	One or more KEY=VALUE SFO Museum API parameters
```

## See also

* https://api.sfomuseum.org/