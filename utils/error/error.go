package error

type typedError struct {
	Type    string
	Message string
}

func NewTypedError(t, msg string) error {
	return &typedError{
		Type:    t,
		Message: msg,
	}
}

func (t typedError) Error() string {
	return t.Message
}

func (t typedError) ErrorType() string {
	return t.Type
}
