package usecase

import "github.com/pkg/errors"

var ErrNotExistsData = errors.New("not exists data error")
var ErrDuplicateData = errors.New("duplicate data error")
var ErrNotExitsUser = errors.New("the user doesn't exist")
var ErrNotCorrectPassword = errors.New("not correct password")
var ErrDecodeImage = errors.New("image decode failure")
