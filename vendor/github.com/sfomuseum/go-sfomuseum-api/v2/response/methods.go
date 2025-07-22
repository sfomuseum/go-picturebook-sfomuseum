package response

// Parameter defines a SFO Museum API method parameter.
type Parameter struct {
	// Name is the name of the API method parameter.
	Name string `json:"name"`
	// Type is a string label describing the type of the API method parameter.
	Type string `json:"type"`
	// Description descibes the API method	parameter.
	Description string `json:"description"`
	// Required indicates with the API method parameter must be included with an API method request.
	Required bool `json:"required"`
	// Example is a stringified example value for the API method	parameter.
	Example string `json:"example"`
}

// Error defines a SFO Museum API error response.
type Error struct {
	// Code is the numeric error code for the error response. See https://api.sfomuseum.org/errors/ for details.
	Code int `json:"code"`
	// Message is the human-readable message describing the error.
	Message string `json:"message"`
}

// Method defines a SFO Museum API method.
type Method struct {
	// Name is the name of the API method.
	Name string `json:"name"`
	// RequestMethod is the HTTP request method that the API method should be invoked with.
	RequestMethod string `json:"request_method"`
	// Description descibes	the API method.
	Description string `json:"description"`
	// RequiresPermissions is the minimum numeric permissions type necessary to invoke the method.
	RequiresPermissions int `json:"requires_perms"`
	// Parameters are zero or more `Parameter` instances associated with the API method.
	Parameters []Parameter `json:"parameters"`
	// Errors are zero or more `Error` instances associated with the API method (in addition to common errors).
	Errors []Error `json:"errors"`
	// Notes are zero or more text-based notes relating to the API method.
	Notes []string `json:"notes"`
	// Paginated is a boolean value indicating whether the API method returns paginated responses.
	Paginated bool `json:"paginated"`
}
