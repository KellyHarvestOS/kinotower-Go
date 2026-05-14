package validator

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

func User(fio, email, birthday string, genderID int, password *string) error {
	if len(strings.TrimSpace(fio)) < 2 || len(fio) > 150 {
		return errors.New("fio must be 2..150 characters")
	}
	if !Email(email) {
		return errors.New("email is invalid")
	}
	if password != nil && (len(*password) < 6 || len(*password) > 65536) {
		return errors.New("password must be 6..65536 characters")
	}
	if _, err := time.Parse("2006-01-02", birthday); err != nil {
		return errors.New("birthday must be YYYY-MM-DD")
	}
	if genderID < 1 {
		return errors.New("gender_id is required")
	}
	return nil
}

func Email(value string) bool {
	_, err := mail.ParseAddress(value)
	return err == nil && len(value) >= 4 && len(value) <= 50
}
