package testing

var RespSimpleSuccess = `
{
  "status":200,
  "message":"success"
}
`

var ErrForbidden = `
{
  "status": 403,
  "message": "not allowed to guest user"
}
`

var ErrNotAllowedMethod = `
{
  "status": 405,
  "message": "not allowed method error"
}
`
