package errors

func NewBaseError(message string, title string, fault string, code int) BaseError {
	return BaseError{
		message: message,
		title:   title,
		fault:   fault,
		code:    code,
	}
}

func User(message string, title string) BaseError {
	return NewBaseError(message, title, "user", 400)
}

func Access(message string) BaseError {
	return NewBaseError(message, "You are not allowed to access this resource", "access", 401)
}

func Server(message string, title string) BaseError {
	return NewBaseError(message, title, "server", 500)
}

func NotFound(message string, title string) BaseError {
	return NewBaseError(message, title, "user", 404)
}

func Validation(message string) BaseError {
	return NewBaseError(message, "Data validation error", "data", 400)
}
