package controller

var errorNotAllowedMethod = `
{
  "status": 405,
  "message": "not allowed method"
}
`

var errUserNameResetKeyNotExistsResetKeyExpired = "the user name doesn't exist or the password reset key doesn't exist or the password reset key is expired"
