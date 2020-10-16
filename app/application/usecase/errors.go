package usecase

import "github.com/pkg/errors"

var ErrNotCorrectPassword = errors.New("not correct password")
var ErrDecodeImage = errors.New("can't decode this file")
