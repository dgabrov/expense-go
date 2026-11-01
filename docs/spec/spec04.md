In @internal/proj/controller/common.go implement function checkLogged, returns error or of course nil 

- get the cookie that is set at /login
- if not present, error
- if present, look for it in the session table (field is token)
- if not found, error
- if found, see if expired_ind is 'Y', if so, error
- if not expired_ind, see if expired_dt is after the current timestamp. If not, error
- otherwise return nil 
