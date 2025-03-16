package validate

import (
	"net/url"
	"regexp"
)

var Name = regexp.MustCompile(`^[a-zA-Z'\-]{1,49}$`)
var FirstNameValidator = Validator{
	Regex:   Name,
	Message: "Invalid first name",
}
var LastNameValidator = Validator{
	Regex:   Name,
	Message: "Invalid last name",
}
var NameValidator = Validator{
	Regex:   Name,
	Message: "Invalid name",
}
var EmailValidator = Validator{
	Regex:   regexp.MustCompile(`^[a-zA-Z0-9\._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`),
	Message: "Invalid email",
}
var PasswordValidator = Validator{
	Regex:   regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`),
	Message: "Passwords must be at least 8 characters long and can be a letter, digit, or special character",
}

var loginValidatorMap = map[string]*Validator{
	"email":    &EmailValidator,
	"password": &PasswordValidator,
}
var signupValidatorMap = map[string]*Validator{
	"email":      &EmailValidator,
	"password":   &PasswordValidator,
	"first_name": &FirstNameValidator,
	"last_name":  &LastNameValidator,
}

func ValidateSignup(form url.Values) []error {
	return validateForm(signupValidatorMap, form, "email", "password", "first_name", "last_name")
}

func ValidateLogin(form url.Values) []error {
	return validateForm(signupValidatorMap, form, "email", "password")
}
