package utils

import (
	"errors"
	"fmt"
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func Validate(model interface{}) (string, bool) {
	err_required := ""
	err_email := ""
	err_number := ""
	err_alpha := ""
	err_alphanum := ""

	validate := validator.New()
	err := validate.Struct(model)

	if err != nil {
		errs := err.(validator.ValidationErrors)
		for i, e := range errs {
			field := ""
			if e.Tag() == "required" {
				switch e.Field() {
				case "Sku":
					field = "SKU"
				case "LocationCode":
					field = "LocationId Code"
				default:
					field = e.Field()
				}
				if i == len(errs)-1 {
					err_required = err_required + field + " "
				} else {
					err_required = err_required + field + ", "
				}
			}
			if e.Tag() == "email" {
				if i == len(errs)-1 {
					err_email = e.Field()
				} else {
					err_email = e.Field() + ","
				}
			}
			if e.Tag() == "number" {
				if i == len(errs)-1 {
					err_number = e.Field()
				} else {
					err_number = e.Field() + ","
				}
			}
			if e.Tag() == "alpha" {
				if i == len(errs)-1 {
					err_alpha = e.Field()
				} else {
					err_alpha = e.Field() + ","
				}
			}
			if e.Tag() == "alphanum" {
				fmt.Println("masuk alfa num")
				if i == len(errs)-1 {
					err_alphanum = e.Field()
				} else {
					err_alphanum = e.Field() + ","
				}
			}
		}
	}

	if len(err_required) > 0 {
		return " Please provide " + err_required + "fields", false
	} else if len(err_email) > 0 {
		return "Your email is wrong", false
	} else if len(err_number) > 0 {
		return " Format field " + string(err_number) + " must number", false
	} else if len(err_alpha) > 0 {
		return " Format field " + err_alphanum + " must alphabet numeric", false
	}
	return "", true
}

// CheckUsername is used to check if a username is valid
// It returns an error if the username is invalid
func CheckUsername(username string) error {
	if !regexp.MustCompile(`^[a-z0-9_]{4,16}$`).MatchString(username) {
		return errors.New("username must be have 4-16 characters and only contains lowercase letters, numbers and underscore")
	}

	return nil
}

// CheckEmail is used to check if an email is valid
// It returns an error if the email is invalid
func CheckEmail(email string) error {
	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email) {
		return errors.New("email must be a valid email address")
	}

	return nil
}

// CheckPassword is used to check if a password is valid
// It returns an error if the password is invalid
func CheckPassword(password string) error {
	if len(password) < 6 {
		return errors.New("password is too short")
	}

	var (
		lowercaseRegex = regexp.MustCompile(`[a-z]`)
		uppercaseRegex = regexp.MustCompile(`[A-Z]`)
		numberRegex    = regexp.MustCompile(`[0-9]`)
		specialRegex   = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	)

	if !lowercaseRegex.MatchString(password) || !uppercaseRegex.MatchString(password) || !numberRegex.MatchString(password) || !specialRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter, one uppercase letter, one number and one special character")
	}

	return nil
}

// CheckMobileNumber is used to check if a mobile phone number is valid
// It returns an error if the mobile phone number is invalid
func CheckMobileNumber(mobile string) error {
	if !regexp.MustCompile(`^\+?[0-9]{10,13}$`).MatchString(mobile) {
		return errors.New("mobile phone number must be 10-13 digits long and contain only numbers")
	}

	return nil
}
