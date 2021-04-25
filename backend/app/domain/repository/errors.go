package repository

import "errors"

var ErrNotExistsData = errors.New("not exists data error")
var ErrDuplicateData = errors.New("duplicate data error")
var ErrUserActivationNotFound = errors.New("wrong user_name or activation_key, or might be already activated")
