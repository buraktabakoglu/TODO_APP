package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "nickname") {
		return errors.New("nickname Already Taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("email Already Taken")
	}

	if strings.Contains(err, "title") {
		return errors.New("title Already Taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("ıncorrect Password")
	}
	return errors.New("ıncorrect Details")
}
