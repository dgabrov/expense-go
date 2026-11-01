package data

type Group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type CompleteGroup struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}
