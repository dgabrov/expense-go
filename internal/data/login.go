package data

type LoginData struct {
	Login    string
	Password string
}

type AuthLoginResponse struct {
	Id     string
	Name   string
	Login  string
	Rights []string
}

type LoggedUser struct {
	Id    string
	Login string
	Name  string
}
