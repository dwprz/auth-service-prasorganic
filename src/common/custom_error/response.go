package custom_error

type ResponseError struct {
	Code    int
	Message string
}

func (err *ResponseError) Error() string {
	return err.Message
}
