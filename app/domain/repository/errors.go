package repository

import "errors"

var ErrNotExistsData = errors.New("not exists data error")
var ErrDuplicateData = errors.New("duplicate data error")
var ErrUserActivationNotFound = errors.New("no such user_name and activation_key")
