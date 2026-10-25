package data

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthLoginResponse struct {
	Id     string
	Name   string
	Login  string
	Rights []string
}

type LoggedUser struct {
	Id    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}
