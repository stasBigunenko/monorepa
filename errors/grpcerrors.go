package errors

type GrpcError string

func (e GrpcError) Error() string {
	return string(e)
}

const (
	DeadlineExceeded GrpcError = "deadline exceeded"
	NotFound         GrpcError = "not found"
	AlreadyExists    GrpcError = "already exists"
	ParseError       GrpcError = "failed to parse"
)
