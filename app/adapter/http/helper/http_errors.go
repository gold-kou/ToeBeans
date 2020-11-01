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

type InternalServerError struct {
	Message string
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{
		Message: message,
	}
}

func (b *BadRequestError) Error() string {
	return b.Message
}

func NewAuthorizationError(message string) *AuthorizationError {
	return &AuthorizationError{
		Message: message,
	}
}

func (a *AuthorizationError) Error() string {
	return a.Message
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

func (i *NotFoundError) Error() string {
	return i.Message
}

func NewInternalServerError(message string) *InternalServerError {
	return &InternalServerError{
		Message: message,
	}
}

func (i *InternalServerError) Error() string {
	return i.Message
}
