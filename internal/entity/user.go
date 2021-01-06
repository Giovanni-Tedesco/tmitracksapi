package entity

type User struct {
	Id       string `json:id`
	Name     string `json:name`
	Location string `json:location`
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetLocation() string {
	return u.Location
}
