package helper

// BadRequestError represents error of bad request
type BadRequestError struct {
	Message string
}

// AuthorizationError represents error of authorization
type AuthorizationError struct {
	Message string
}

// ConflictError represents error of conflict
type ConflictError struct {
	Message string
}

// NotFoundError represents error of not found
type NotFoundError struct {
	Message string
}

// InternalServerError represents error of internal server
type InternalServerError struct {
	Message string
}

// NewBadRequestError returns a new pointer of BadRequestError that contains a message of an argument
func NewBadRequestError(message string, errors ...map[string]string) *BadRequestError {
	return &BadRequestError{
		Message: message,
	}
}

func (b *BadRequestError) Error() string {
	return b.Message
}

// NewAuthorizationError returns a new pointer of AuthorizationError that contains a message of an argument
func NewAuthorizationError(message string) *AuthorizationError {
	return &AuthorizationError{
		Message: message,
	}
}

func (a *AuthorizationError) Error() string {
	return a.Message
}

// NewNotFoundError returns a new pointer of NotFoundError that contains a message of an argument
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

func (i *NotFoundError) Error() string {
	return i.Message
}

// NewInternalServerError returns a new pointer of InternalServerError that contains a message of an argument
func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{
		Message: message,
	}
}

func (i *InternalServerError) Error() string {
	return i.Message
}
