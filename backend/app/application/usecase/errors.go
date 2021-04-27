package usecase

import "github.com/pkg/errors"

var ErrNotVerifiedUser = errors.New("not email verified user")
var ErrNotCorrectPassword = errors.New("not correct password")
var ErrDecodeImage = errors.New("image decode failure")
var ErrOverPasswordResetCount = errors.New("you can't reset password as it exceeds limit counts")
var ErrLikeYourSelf = errors.New("you can't like your posting")
