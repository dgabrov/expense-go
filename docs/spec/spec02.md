@internal/proj/controller/post_login.go process function

- get request, payload parse to LoginData
- this payload post it against url that is at auth in the configData
- get back a response. If http code is 400+, return error. 
- if code is 200, parse to AuthLoginResponse
- create entry in table session (new method in server.Server etc)
  - session_id new v7 guid
  - token a new one created with createRandomToken in @common.go
  - user_id, login, name taken from AuthLoginResponse
  - expiry_dt in the future with tokenTTLSeconds from configData

- assemble LoggedUser and return it
