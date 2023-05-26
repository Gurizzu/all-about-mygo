package util

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

func GenerateID() (res string) {
	res = strings.ReplaceAll(uuid.New().String(), "-", "")
	return
}

func ConvertToReader(val interface{}) (*strings.Reader, error) {
	str, ok := val.(string)
	if !ok {
		return nil, errors.New("value is not a string")
	}
	return strings.NewReader(str), nil
}
