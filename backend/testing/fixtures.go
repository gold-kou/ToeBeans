package testing

var RespSimpleSuccess = `
{
  "message":"success", 
  "status":200
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
