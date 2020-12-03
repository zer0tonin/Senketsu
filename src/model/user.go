package model

type UserRepository interface {
}

type User struct {
	Name     string `json:"name"`
	Images   []string `json:"images"`
}

func NewUser(name string) *User {
	return &User{
		Name: name,
		Images: make([]string, 0),
	}
}
