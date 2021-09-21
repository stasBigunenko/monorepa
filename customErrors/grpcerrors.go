package customErrors

type GRPCError string

func (e GRPCError) Error() string {
	return string(e)
}

const (
	DeadlineExceeded GRPCError = "deadline exceeded"
	NotFound         GRPCError = "not found"
	AlreadyExists    GRPCError = "already exists"
	ParseError       GRPCError = "failed to parse"
)
