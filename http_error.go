package error_view

type HTTPError interface {
	error

	Code() int
	Message() string
}
