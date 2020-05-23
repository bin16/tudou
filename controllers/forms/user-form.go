package forms

import "github.com/bin16/tudou/models"

// SignupForm to create user
type SignupForm struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Valid method return if form is valid and error message
func (sf SignupForm) Valid() (bool, string) {
	// TODO: Email check`
	if sf.Email == "" {
		return false, "Email is required"
	}

	// TODO: Name check
	if sf.Name == "" {
		return false, "Name is required"
	}

	return true, ""
}

// ToUser get user struct
func (sf SignupForm) ToUser() models.User {
	return models.User{
		Name:  sf.Name,
		Email: sf.Email,
	}
}
