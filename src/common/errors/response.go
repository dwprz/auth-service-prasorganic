package errors

import "google.golang.org/grpc/codes"

type Response struct {
	HttpCode int
	GrpcCode codes.Code
	Message  string
}

func (err *Response) Error() string {
	return err.Message
}