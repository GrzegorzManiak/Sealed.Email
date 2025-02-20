package errors

type AppError interface {
	Message() string
	Title() string
	Fault() string
	Code() int

	Error() string
}
