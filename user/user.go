package user

import "errors"

type User struct {
	Username string
	Room     string
}

func UserFromMap(m map[string]interface{}) (*User, error) {
	u := &User{}

	if name, ok := m["Username"].(string); ok {
		u.Username = name
	} else {
		return nil, errors.New("InvalidName")
	}

	if room, ok := m["Room"].(string); ok {
		u.Room = room
	} else {
		return nil, errors.New("InvalidRoom")
	}

	return u, nil
}
