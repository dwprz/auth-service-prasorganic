package errors

type Response struct {
	Code    int
	Message string
}

func (err *Response) Error() string {
	return err.Message
}
