package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Field string
}

func (e ValidationError) Error() string {
	return "invalid field: " + e.Field
}

func validate(age int) error {
	if age < 0 {
		return ValidationError{Field: "age"}
	}
	return nil
}

func main() {
	err := validate(20)
	fmt.Printf("Age 20 -> error: %v, is nil? %v, type: %T\n", err, err == nil, err)

	errNegative := validate(-1)
	if errNegative != nil {
		var valErr ValidationError
		if errors.As(errNegative, &valErr) {
			fmt.Printf("Validation failed on field: %s\n", valErr.Field)
		}
	}
}
