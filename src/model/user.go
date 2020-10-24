package model

type UserRepository interface {
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
