package lib

import "errors"

var (
	ClassNotFoundError   = errors.New("Class not found!")
	UninitializedKBError = errors.New("Uninitialized KB!")
	InvalidClassError    = errors.New("Invalid class!")
)
