package errors

const (
	WrongPassword UserValidationError = "wrong password"
)

type UserValidationError string

func (e UserValidationError) Error() string {
	return string(e)
}
