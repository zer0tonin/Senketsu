package model

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

func (u *User) GetID() string {
	return u.Name
}
