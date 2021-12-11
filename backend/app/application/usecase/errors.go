package usecase

import "github.com/pkg/errors"

var ErrNotExitsUser = errors.New("the user doesn't exist")
var ErrNotCorrectPassword = errors.New("not correct password")
var ErrDecodeImage = errors.New("image decode failure")
