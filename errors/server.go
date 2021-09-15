package errors

const (
	WrongPort ServerError = "wrong port"
)

type ServerError string

func (e ServerError) Error() string {
	return string(e)
}
