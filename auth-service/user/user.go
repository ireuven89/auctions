package user

import (
	"encoding/json"
)

type User struct {
	ID    string
	Name  string
	Email string
	// <-- do NOT include in public struct or JSON output
	Password string
	// <-- do NOT include in public struct or JSON output
}

func (user *User) toString() string {
	//Do not add the password to the ToString method
	return "id:" + user.ID + "name:" + user.Name + "email:" + user.Email
}

// ToJson - Do not add the password to the ToJson method
func ToJson(user User) string {
	safeUser := UserResponse{
		ID:       user.ID,
		Username: user.Name,
		Email:    user.Email,
	}
	data, err := json.Marshal(safeUser)
	if err != nil {
		return ""
	}
	return string(data)
}

func FromJson(jsonUser []byte) (*UserResponse, error) {
	var user UserResponse
	if err := json.Unmarshal(jsonUser, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
