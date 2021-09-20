package customErrors

type HTTPError struct {
	Message string `json:"message,omitempty"`
}

func (e HTTPError) Error() string {
	return e.Message
}

var UUIDError = HTTPError{
	Message: "failed to parse uuid",
}

var RWError = HTTPError{
	Message: "failed to read / write to io",
}

var JSONError = HTTPError{
	Message: "failed to marshal / unmarshal json",
}
