package helper

type BadRequestError struct {
	Message string
}

type AuthorizationError struct {
	Message string
}

type NotFoundError struct {
	Message string
}

type ForbiddenError struct {
	Message string
}

type InternalServerError struct {
	Message string
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{
		Message: message,
	}
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func NewAuthorizationError(message string) *AuthorizationError {
	return &AuthorizationError{
		Message: message,
	}
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{
		Message: message,
	}
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{
		Message: message,
	}
}

func (e *InternalServerError) Error() string {
	return e.Message
}
