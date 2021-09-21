package customErrors

type ServerError string

func (e ServerError) Error() string {
	return string(e)
}

const (
	WrongPort ServerError = "wrong port"
)
