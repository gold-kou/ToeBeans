package usecase

import "github.com/pkg/errors"

var ErrNotCorrectPassword = errors.New("not correct password")
var ErrDecodeImage = errors.New("can't decode this file")
var ErrOverPasswordResetCount = errors.New("can't reset password today as it exceeds limit counts")
var ErrLikeYourSelf = errors.New("you can't like your posting")
